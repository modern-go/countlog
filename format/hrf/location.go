package hrf

import (
	"github.com/modern-go/countlog/format"
	"github.com/modern-go/countlog/logger"
)

func formatLocation(site *logger.LogSite) format.Formatter {
	loc := site.Location()
	return format.FuncFormatter(func(space []byte, event *logger.Event) []byte {
		if event.Level < logger.LevelWarn {
			return space
		}
		space = append(space, "\n\x1b[90;1mlocation: "...)
		space = append(space, loc...)
		space = append(space, "\x1b[0m"...)
		return space
	})
}