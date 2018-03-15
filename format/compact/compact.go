package compact

import (
	"github.com/modern-go/countlog/logger"
	"github.com/modern-go/countlog/format"
	"fmt"
	"strings"
)

type Format struct {
	TimeFormat     string
	HideLevel      bool
	HideTime       bool
	HideLocation   bool
	HideContext    bool
	HideProperties bool
}

func (f *Format) FormatterOf(site *logger.LogSite) format.Formatter {
	var formatters format.Formatters
	if !f.HideLevel {
		formatters = append(formatters, formatLevel())
	}
	if !f.HideTime {
		formatters = append(formatters, formatTime(f.TimeFormat))
	}
	if !f.HideLocation {
		formatters = append(formatters, formatLiteral(fmt.Sprintf(
			"[%s] ", site.Location())))
	}
	eventName := site.Event
	if strings.HasPrefix(eventName, "event!") {
		formatters = append(formatters, formatLiteral(eventName[len("event!"):]))
	} else if strings.HasPrefix(eventName, "callee!") {
		msg := "call " + eventName[len("callee!"):]
		formatters = append(formatters, formatLiteral(msg))
	} else {
		formatters = append(formatters, formatProperties(eventName, site.Sample))
	}
	formatters = append(formatters, formatError())
	ctx := logger.GetLogContext(site.Context)
	if !f.HideContext && ctx != nil {
		sample := ctx.Properties
		for i := 0; i < len(sample); i += 2 {
			key := sample[i].(string)
			formatters = append(formatters, formatContext(key, sample))
		}
	}
	if !f.HideProperties {
		sample := site.Sample
		for i := 0; i < len(sample); i += 2 {
			key := sample[i].(string)
			formatters = append(formatters, formatProperties(key, sample))
		}
	}
	formatters = append(formatters, formatLiteral("\n"))
	return formatters
}
