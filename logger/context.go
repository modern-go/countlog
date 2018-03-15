package logger

import "context"

var LogContextKey = 1010010001

type LogContext struct {
	Memos      [][]byte
	Properties []interface{}
}

func (ctx *LogContext) Add(key string, value interface{}) {
	if ctx == nil {
		return
	}
	ctx.Properties = append(ctx.Properties, key)
	ctx.Properties = append(ctx.Properties, value)
}

func GetLogContext(ctx context.Context) *LogContext {
	if ctx == nil {
		return nil
	}
	logContext, _ := ctx.Value(LogContextKey).(*LogContext)
	return logContext
}

func WithLogContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, LogContextKey, &LogContext{})
}
