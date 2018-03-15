package hrf

import (
	"github.com/modern-go/countlog/logger"
	"github.com/modern-go/countlog/format"
)

func formatLevel() format.Formatter {
	return format.FuncFormatter(func(space []byte, event *logger.Event) []byte {
		return append(space, logger.ColoredLevelName(event.Level)...)
	})
}