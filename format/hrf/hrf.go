package hrf

import (
	"github.com/modern-go/countlog/logger"
	"github.com/modern-go/countlog/format"
	"strings"
)

type Format struct {
	TimeFormat     string
	HideLevel      bool
	HideTime       bool
	HideContext    bool
	HideProperties bool
	HideLocation   bool
}

func (f *Format) FormatterOf(site *logger.LogSite) format.Formatter {
	var formatters format.Formatters
	if !f.HideLevel {
		formatters = append(formatters, formatLevel())
	}
	if strings.HasPrefix(site.Event, "event!") {
		formatters = append(formatters, formatLiteral(site.Event[len("event!"):]))
	} else if strings.HasPrefix(site.Event, "callee!") {
		formatters = append(formatters, formatLiteral("call "+site.Event[len("callee!"):]))
	} else {
		formatters = append(formatters, formatProperties(site.Event, site.Sample))
	}
	formatters = append(formatters, formatError())
	if !f.HideTime {
		formatters = append(formatters, formatTime(f.TimeFormat))
	}
	ctx := site.LogContext()
	if !f.HideContext && ctx != nil {
		for i := 0; i < len(ctx.Properties); i += 2 {
			key := ctx.Properties[i].(string)
			formatters = append(formatters, formatContext(key, ctx.Properties))
		}
	}
	if !f.HideLocation {
		formatters = append(formatters, formatLocation(site))
	}
	formatters = append(formatters, formatLiteral("\n"))
	return formatters
}
