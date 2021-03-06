package compact

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
		space = append(space, '[')
		space = event.Timestamp.AppendFormat(space, timeFormat)
		space = append(space, ']', ' ')
		return space
	})
}
