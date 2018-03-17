package hrf

import (
	"github.com/modern-go/countlog/format"
	"github.com/modern-go/countlog/logger"
	"time"
)

func formatTime(timeFormat string) format.Formatter {
	if timeFormat == "" {
		timeFormat = time.RFC3339
	}
	return format.FuncFormatter(func(space []byte, event *logger.Event) []byte {
		space = append(space, "\n\x1b[90;1mtimestamp: "...)
		space = event.Timestamp.AppendFormat(space, timeFormat)
		space = append(space, "\x1b[0m"...)
		return space
	})
}
