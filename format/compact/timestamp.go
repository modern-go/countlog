package compact

import (
	"github.com/modern-go/countlog/format"
	"github.com/modern-go/countlog/logger"
	"time"
)

func formatTime(dateFormat string) format.Formatter {
	if dateFormat == "" {
		dateFormat = time.RFC3339
	}
	return format.FuncFormatter(func(space []byte, event *logger.Event) []byte {
		space = append(space, '[')
		space = append(space, logger.LevelName(event.Level)...)
		space = append(space, ']', ' ', '[')
		space = append(event.Timestamp.AppendFormat(space, dateFormat))
		space = append(space, ']', ' ')
		return space
	})
}
