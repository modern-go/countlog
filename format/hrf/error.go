package hrf

import (
	"github.com/modern-go/countlog/format"
	"github.com/modern-go/countlog/logger"
)

func formatError() format.Formatter {
	return format.FuncFormatter(func(space []byte, event *logger.Event) []byte {
		if event.Error == nil {
			return space
		}
		msg := event.Error.Error()
		if msg == "" {
			msg = "error"
		}
		space = append(space, "\n\x1b[90;1merror: "...)
		space = append(space, msg...)
		space = append(space, "\x1b[0m"...)
		return space
	})
}
