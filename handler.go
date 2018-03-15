package countlog

import (
	"errors"
	"context"
	"github.com/modern-go/countlog/logger"
	"time"
	"github.com/modern-go/concurrent"
	"github.com/modern-go/msgfmt/formatter"
	"runtime"
)

var handlerCache = concurrent.NewMap()

func log(level int, eventName string, agg string, ctx context.Context, err error, properties []interface{}) error {
	handler := getHandler(eventName, agg, ctx, properties)
	event := &logger.Event{
		Level:      level,
		Context:    ctx,
		Error:      err,
		Timestamp:  time.Now(),
		Properties: properties,
	}
	handler.Handle(event)
	if event.Error != nil {
		fmt := formatter.Of(eventName, properties)
		errMsg := fmt.Format(nil, properties)
		errMsg = append(errMsg, ": "...)
		errMsg = append(errMsg, event.Error.Error()...)
		event.Error = errors.New(string(errMsg))
	}
	return event.Error
}

func getHandler(event string, agg string, ctx context.Context, properties []interface{}) logger.EventHandler {
	handler, found := handlerCache.Load(event)
	if found {
		return handler.(logger.EventHandler)
	}
	return newHandler(event, agg, ctx, properties)
}

func newHandler(eventName string, agg string, ctx context.Context, properties []interface{}) logger.EventHandler {
	pc, callerFile, callerLine, _ := runtime.Caller(4)
	site := &logger.LogSite{
		Context: ctx,
		Func:    runtime.FuncForPC(pc).Name(),
		Event:   eventName,
		Agg:     agg,
		File:    callerFile,
		Line:    callerLine,
		Sample:  properties,
	}
	handler := newRootHandler(site, nil)
	handlerCache.Store(eventName, handler)
	return handler
}
