package output

import (
	"io"
	"sync"
	"os"
	"github.com/modern-go/countlog/format"
	"github.com/modern-go/countlog/logger"
)

type EventWriter struct {
	Format format.Format
	Writer io.Writer
}

func NewEventWriter(initWriter func(eventWriter *EventWriter)) *EventWriter {
	eventWriter := &EventWriter{
		Writer: os.Stdout,
	}
	initWriter(eventWriter)
	eventWriter.Writer = &recylceWriter{eventWriter.Writer}
	return eventWriter
}

func (sink *EventWriter) HandlerOf(site *logger.LogSite) logger.EventHandler {
	formatter := sink.Format.FormatterOf(site)
	return &writeEvent{
		site:      site,
		formatter: formatter,
		writer:    sink.Writer,
	}
}

type writeEvent struct {
	site      *logger.LogSite
	formatter format.Formatter
	writer    io.Writer
}

func (handler *writeEvent) Handle(event *logger.Event) {
	space := bufPool.Get().([]byte)[:0]
	formatted := handler.formatter.Format(space, event)
	_, err := handler.writer.Write(formatted)
	if err != nil {
		logger.ErrorLogger.Println("failed to write formatted event", err)
	}
}

func (handler *writeEvent) LogSite() *logger.LogSite {
	return handler.site
}

var bufPool = &sync.Pool{
	New: func() interface{} {
		return make([]byte, 0, 128)
	},
}

type recylceWriter struct {
	writer io.Writer
}

func (writer *recylceWriter) Write(buf []byte) (int, error) {
	n, err := writer.writer.Write(buf)
	bufPool.Put(buf)
	return n, err
}
