package compact

import (
	"github.com/modern-go/countlog/format"
	"github.com/modern-go/countlog/logger"
)

func formatError() format.Formatter {
	return format.FuncFormatter(func(space []byte, event *logger.Event) []byte {
		if event.Error == nil {
			return space
		}
		space = append(space, ':', ' ')
		space = append(space, event.Error.Error()...)
		return space
	})
}
