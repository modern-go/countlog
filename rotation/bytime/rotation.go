package bytime

import (
	"context"
	"github.com/modern-go/concurrent"
	"github.com/modern-go/countlog/logger"
	"github.com/modern-go/msgfmt"
	"math/rand"
	"os"
	"path"
	"sync/atomic"
	"time"
	"unsafe"
)

// normal => triggered => opened new => normal
const statusNormal = 0
const statusTriggered = 1
const statusArchived = 2

type Writer struct {
	WritePath           string
	FileMode            os.FileMode
	DirectoryMode       os.FileMode
	Interval            time.Duration
	ArchiveFilePattern  string
	ArchiveKeepDuration time.Duration
	// file is owned by the write goroutine
	file *os.File
	// newFile and status shared between write and rotate goroutine
	newFile  unsafe.Pointer
	status   int32
	executor *concurrent.UnboundedExecutor
}

func NewRotationWriter(init func(writer *Writer)) (*Writer, error) {
	writer := &Writer{
		FileMode:      0644,
		DirectoryMode: 0755,
		executor:      concurrent.NewUnboundedExecutor(),
		Interval:      time.Hour,
	}
	init(writer)
	err := writer.reopen()
	if err != nil {
		return nil, err
	}
	writer.executor.Go(writer.rotateInBackground)
	return writer, nil
}

func (writer *Writer) reopen() error {
	file, err := os.OpenFile(writer.WritePath, os.O_WRONLY|os.O_APPEND, writer.FileMode)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		os.MkdirAll(path.Dir(writer.WritePath), writer.DirectoryMode)
		file, err = os.OpenFile(writer.WritePath,
			os.O_CREATE|os.O_WRONLY|os.O_TRUNC, writer.FileMode)
		if err != nil {
			return err
		}
	}
	writer.file = file
	return nil
}

func (writer *Writer) Close() error {
	writer.executor.StopAndWaitForever()
	return writer.file.Close()
}

func (writer *Writer) Write(buf []byte) (int, error) {
	if atomic.LoadInt32(&writer.status) == statusArchived {
		err := writer.file.Close()
		if err != nil {
			logger.ErrorLogger.Println("close rotation log failed", err)
		}
		err = writer.reopen()
		if err != nil {
			logger.ErrorLogger.Println("open rotation log failed", err)
		}
	}
	return writer.file.Write(buf)
}

func (writer *Writer) rotateInBackground(ctx context.Context) {
	var timer <-chan time.Time
	for {
		duration := writer.Interval
		if duration > 0 {
			randomMax := duration
			if randomMax > time.Minute*5 {
				randomMax = time.Minute * 5
			}
			randomGap := time.Duration(rand.Int63n(int64(randomMax)))
			duration += randomGap
			timer = time.NewTimer(duration).C
		}
		select {
		case <-ctx.Done():
			return
		case <-timer:
		}
		archivePath := msgfmt.Sprintf(writer.ArchiveFilePattern, "time", time.Now())
		archivePath = path.Join(path.Dir(writer.WritePath), archivePath)
		err := os.Rename(writer.WritePath, archivePath)
		if err != nil {
			logger.ErrorLogger.Println("failed to move archive", err)
			continue
		}
		atomic.StoreInt32(&writer.status, statusArchived)
		writer.purgeExpired()
	}
}

func (writer *Writer) purgeExpired() {
	archivePath := path.Dir(writer.WritePath)
	dir, err := os.Open(archivePath)
	if err != nil {
		logger.ErrorLogger.Println("failed to open dir", err)
		return
	}
	files, err := dir.Readdir(0)
	if err != nil {
		logger.ErrorLogger.Println("failed to list dir", err)
		return
	}
	now := time.Now()
	for _, file := range files {
		if now.Sub(file.ModTime()) > writer.ArchiveKeepDuration {
			os.Remove(path.Join(archivePath, file.Name()))
		}
	}
}
