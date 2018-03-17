package countlog

import (
	"context"
	"github.com/modern-go/countlog/logger"
	"github.com/modern-go/reflect2"
	"unsafe"
)

func Ctx(ctx context.Context) *Context {
	wrapped, isWrapped := ctx.(*Context)
	if isWrapped {
		return wrapped
	}
	return &Context{Context: ctx, logContext: &logger.LogContext{}}
}

type Context struct {
	context.Context
	logContext *logger.LogContext
}

func (ctx *Context) Value(key interface{}) interface{} {
	if ctx == nil {
		return nil
	}
	if key == logger.LogContextKey {
		return ctx.logContext
	}
	return ctx.Context.Value(key)
}

func (ctx *Context) Trace(event string, properties ...interface{}) {
	if LevelTrace < logger.MinLevel {
		return
	}
	ptr := reflect2.NoEscape(unsafe.Pointer(&properties))
	log(LevelTrace, event, "", ctx, nil, *(*[]interface{})(ptr))
}

func (ctx *Context) TraceCall(event string, err error, properties ...interface{}) error {
	if err != nil {
		ptr := reflect2.NoEscape(unsafe.Pointer(&properties))
		return log(LevelWarn, event, "call", ctx, err, *(*[]interface{})(ptr))
	}
	if LevelTrace < logger.MinLevel {
		return nil
	}
	ptr := reflect2.NoEscape(unsafe.Pointer(&properties))
	log(LevelTrace, event, "call", ctx, err, *(*[]interface{})(ptr))
	return nil
}

func (ctx *Context) Debug(event string, properties ...interface{}) {
	if LevelDebug < logger.MinLevel {
		return
	}
	ptr := reflect2.NoEscape(unsafe.Pointer(&properties))
	log(LevelDebug, event, "", ctx, nil, *(*[]interface{})(ptr))
}

func (ctx *Context) DebugCall(event string, err error, properties ...interface{}) error {
	if err != nil {
		ptr := reflect2.NoEscape(unsafe.Pointer(&properties))
		return log(LevelWarn, event, "call", ctx, err, *(*[]interface{})(ptr))
	}
	if LevelDebug < logger.MinLevel {
		return nil
	}
	ptr := reflect2.NoEscape(unsafe.Pointer(&properties))
	log(LevelDebug, event, "call", ctx, err, *(*[]interface{})(ptr))
	return nil
}

func (ctx *Context) Info(event string, properties ...interface{}) {
	if LevelInfo < logger.MinLevel {
		return
	}
	ptr := reflect2.NoEscape(unsafe.Pointer(&properties))
	log(LevelInfo, event, "", ctx, nil, *(*[]interface{})(ptr))
}

func (ctx *Context) InfoCall(event string, err error, properties ...interface{}) error {
	if err != nil {
		ptr := reflect2.NoEscape(unsafe.Pointer(&properties))
		return log(LevelWarn, event, "call", ctx, err, *(*[]interface{})(ptr))
	}
	if LevelInfo < logger.MinLevel {
		return nil
	}
	ptr := reflect2.NoEscape(unsafe.Pointer(&properties))
	log(LevelInfo, event, "call", ctx, err, *(*[]interface{})(ptr))
	return nil
}

func (ctx *Context) LogAccess(event string, err error, properties ...interface{}) error {
	if err != nil {
		ptr := reflect2.NoEscape(unsafe.Pointer(&properties))
		return log(LevelError, event, "call", ctx, err, *(*[]interface{})(ptr))
	}
	if LevelInfo < logger.MinLevel {
		return nil
	}
	ptr := reflect2.NoEscape(unsafe.Pointer(&properties))
	log(LevelInfo, event, "call", ctx, err, *(*[]interface{})(ptr))
	return nil
}

func (ctx *Context) Warn(event string, properties ...interface{}) {
	ptr := reflect2.NoEscape(unsafe.Pointer(&properties))
	log(LevelWarn, event, "", ctx, nil, *(*[]interface{})(ptr))
}

func (ctx *Context) Error(event string, properties ...interface{}) {
	ptr := reflect2.NoEscape(unsafe.Pointer(&properties))
	log(LevelError, event, "", ctx, nil, *(*[]interface{})(ptr))
}

func (ctx *Context) Fatal(event string, properties ...interface{}) {
	ptr := reflect2.NoEscape(unsafe.Pointer(&properties))
	log(LevelFatal, event, "", ctx, nil, *(*[]interface{})(ptr))
}

func (ctx *Context) Add(key string, value interface{}) {
	ctx.logContext.Properties = append(ctx.logContext.Properties, key)
	ctx.logContext.Properties = append(ctx.logContext.Properties, value)
}
