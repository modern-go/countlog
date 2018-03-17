package compact

import (
	"github.com/modern-go/countlog/format"
	"github.com/modern-go/countlog/logger"
	"github.com/modern-go/msgfmt/formatter"
)

func formatMessage(pattern string, sample []interface{}) format.Formatter {
	formatter := formatter.Of(pattern, sample)
	return format.FuncFormatter(func(space []byte, event *logger.Event) []byte {
		return formatter.Format(space, event.Properties)
	})
}
