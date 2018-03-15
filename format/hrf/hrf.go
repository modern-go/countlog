package hrf

import (
	"github.com/modern-go/countlog/logger"
	"github.com/modern-go/countlog/format"
	"strings"
)

type Format struct {
}

func (f *Format) FormatterOf(site *logger.LogSite) format.Formatter {
	var formatters format.Formatters
	if strings.HasPrefix(site.Event, "event!") {
		formatters = append(formatters, formatLiteral(site.Event[len("event!"):]))
	} else if strings.HasPrefix(site.Event, "callee!") {
		formatters = append(formatters, formatLiteral("call "+site.Event[len("callee!"):]))
	} else {
		formatters = append(formatters, formatProperties(site.Event, site.Sample))
	}
	formatters = append(formatters, formatLiteral("\n"))
	return formatters
}
