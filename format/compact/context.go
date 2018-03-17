package compact

import (
	"github.com/modern-go/countlog/format"
	"github.com/modern-go/countlog/logger"
	"github.com/modern-go/msgfmt/formatter"
)

func formatContext(key string, sample []interface{}) format.Formatter {
	pattern := "||" + key + "={" + key + "}"
	formatter := formatter.Of(pattern, sample)
	return format.FuncFormatter(func(space []byte, event *logger.Event) []byte {
		ctx := logger.GetLogContext(event.Context)
		return formatter.Format(space, ctx.Properties)
	})
}
