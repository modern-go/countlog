package countlog

import (
	"github.com/modern-go/countlog/logger"
	"runtime"
	"runtime/debug"
)

type panicHandler func(recovered interface{}, event *logger.Event, site *logger.LogSite)

func newRootHandler(site *logger.LogSite, onPanic panicHandler) logger.EventHandler {
	statsHandler := EventAggregator.HandlerOf(site)
	if statsHandler == nil {
		return &oneHandler{
			site:    site,
			handler: EventWriter.HandlerOf(site),
			onPanic: onPanic,
		}
	}
	return &statsAndOutput{
		site:          site,
		statsHandler:  statsHandler,
		outputHandler: EventWriter.HandlerOf(site),
		onPanic:       onPanic,
	}
}

type oneHandler struct {
	site    *logger.LogSite
	handler logger.EventHandler
	onPanic panicHandler
}

func (handler *oneHandler) Handle(event *logger.Event) {
	defer func() {
		recovered := recover()
		if recovered != nil {
			handler.onPanic(recovered, event, handler.site)
		}
	}()
	handler.handler.Handle(event)
}

func (handler *oneHandler) LogSite() *logger.LogSite {
	return handler.site
}

type statsAndOutput struct {
	site          *logger.LogSite
	statsHandler  logger.EventHandler
	outputHandler logger.EventHandler
	onPanic       panicHandler
}

func (handler *statsAndOutput) Handle(event *logger.Event) {
	defer func() {
		recovered := recover()
		if recovered != nil {
			handler.onPanic(recovered, event, handler.site)
		}
	}()
	if event.Level >= logger.MinCallLevel {
		handler.outputHandler.Handle(event)
	}
	handler.statsHandler.Handle(event)
}

func (handler *statsAndOutput) LogSite() *logger.LogSite {
	return handler.site
}

func nomalModeOnPanic(recovered interface{}, event *logger.Event, site *logger.LogSite) {
	redirector := &redirector{
		site: *site,
	}
	handlerCache.Store(site.Event, redirector)
	newSite := *site
	newSite.File = "unknown"
	newSite.Line = 0
	newSite.Sample = event.Properties
	newRootHandler(&newSite, fallbackModeOnPanic).Handle(event)
}

func fallbackModeOnPanic(recovered interface{}, event *logger.Event, site *logger.LogSite) {
	logger.ErrorLogger.Println("panic:", recovered)
	if logger.MinLevel <= logger.LevelDebug {
		logger.ErrorLogger.Println(string(debug.Stack()))
	}
}

type redirector struct {
	site logger.LogSite
}

func (redirector *redirector) Handle(event *logger.Event) {
	_, callerFile, callerLine, _ := runtime.Caller(3)
	key := accurateHandlerKey{callerFile, callerLine}
	handlerObj, found := handlerCache.Load(key)
	if found {
		handlerObj.(logger.EventHandler).Handle(event)
		return
	}
	site := redirector.site
	site.File = callerFile
	site.Line = callerLine
	site.Sample = event.Properties
	handler := newRootHandler(&site, fallbackModeOnPanic)
	handlerCache.Store(key, handler)
	handler.Handle(event)
	return
}

type accurateHandlerKey struct {
	File string
	Line int
}
