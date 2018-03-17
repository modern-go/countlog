package compact

import (
	"github.com/modern-go/countlog/format"
	"github.com/modern-go/countlog/logger"
)

func formatLevel() format.Formatter {
	return format.FuncFormatter(func(space []byte, event *logger.Event) []byte {
		space = append(space, '[')
		space = append(space, logger.LevelName(event.Level)...)
		space = append(space, ']', ' ')
		return space
	})
}
