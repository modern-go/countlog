package compact

import (
	"github.com/modern-go/countlog/format"
	"github.com/modern-go/countlog/logger"
)

func formatLiteral(literal string) format.Formatter {
	return format.FuncFormatter(func(space []byte, event *logger.Event) []byte {
		return append(space, literal...)
	})
}
