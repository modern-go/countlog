package bytime

import (
	"os"
	"unsafe"
	"sync/atomic"
	"path"
	"github.com/modern-go/concurrent"
	"context"
	"time"
	"math/rand"
	"github.com/modern-go/countlog/logger"
)

// normal => triggered => opened new => normal
const statusNormal = 0
const statusTriggered = 1
const statusArchived = 2

type Writer struct {
	WritePath       string
	FileMode        os.FileMode
	DirectoryMode   os.FileMode
	Interval time.Duration
	// file is owned by the write goroutine
	file *os.File
	// newFile and status shared between write and rotate goroutine
	newFile unsafe.Pointer
	status int32
	executor *concurrent.UnboundedExecutor
}

func NewRotationWriter(init func(writer *Writer)) (*Writer, error) {
	writer := &Writer{
		FileMode: 0644,
		DirectoryMode: 0755,
		executor: concurrent.NewUnboundedExecutor(),
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
			if randomMax > time.Minute * 5 {
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
		//archives, err := archiveStrategy.Archive(writer.cfg.WritePath)
		//if err != nil {
		//	loglog.Error(err)
		//	// retry after one minute
		//	timer = time.NewTimer(time.Minute).C
		//	continue
		//}
		//atomic.StoreInt32(&writer.status, statusArchived)
		//purgeSet := retainStrategy.PurgeSet(archives)
		//err = purgeStrategy.Purge(purgeSet)
		//if err != nil {
		//	loglog.Error(err)
		//}
	}
}