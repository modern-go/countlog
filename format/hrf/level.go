package hrf

import (
	"github.com/modern-go/countlog/format"
	"github.com/modern-go/countlog/logger"
)

func formatLevel() format.Formatter {
	return format.FuncFormatter(func(space []byte, event *logger.Event) []byte {
		return append(space, logger.ColoredLevelName(event.Level)...)
	})
}
