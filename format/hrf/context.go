package hrf

import (
	"github.com/modern-go/msgfmt/formatter"
	"github.com/modern-go/countlog/logger"
	"github.com/modern-go/countlog/format"
)

func formatContext(key string, sample []interface{}) format.Formatter {
	pattern := "\n\x1b[90;1m" + key + ": {" + key + "}\x1b[0m"
	formatter := formatter.Of(pattern, sample)
	return format.FuncFormatter(func(space []byte, event *logger.Event) []byte {
		ctx := logger.GetLogContext(event.Context)
		return formatter.Format(space, ctx.Properties)
	})
}