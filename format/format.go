package format

import "github.com/modern-go/countlog/logger"

type Format interface {
	FormatterOf(site *logger.LogSite) Formatter
}

type Formatter interface {
	Format(space []byte, event *logger.Event) []byte
}

type Formatters []Formatter

func (formatters Formatters) Format(space []byte, event *logger.Event) []byte {
	for _, formatter := range formatters {
		space = formatter.Format(space, event)
	}
	return space
}

type FuncFormatter func(space []byte, event *logger.Event) []byte

func (f FuncFormatter) Format(space []byte, event *logger.Event) []byte {
	return f(space, event)
}
