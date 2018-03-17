package hrf

import (
	"github.com/modern-go/countlog/format"
	"github.com/modern-go/countlog/logger"
	"github.com/modern-go/msgfmt/formatter"
)

func formatMessage(msgFormat string, sample []interface{}) format.Formatter {
	formatter := formatter.Of(msgFormat, sample)
	return format.FuncFormatter(func(space []byte, event *logger.Event) []byte {
		return formatter.Format(space, event.Properties)
	})
}
