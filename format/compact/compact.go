package compact

import (
	"github.com/modern-go/countlog/logger"
	"github.com/modern-go/countlog/format"
	"fmt"
	"strings"
)

type Format struct {
	timeFormat   string
	hideLevel    bool
	hideTime     bool
	hideLocation bool
}

func (f Format) FormatterOf(site *logger.LogSite) format.Formatter {
	var formatters format.Formatters
	if !f.hideLevel {
		formatters = append(formatters, formatLevel())
	}
	if !f.hideTime {
		formatters = append(formatters, formatTime(f.timeFormat))
	}
	if !f.hideLocation {
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
	sample := site.Sample
	for i := 0; i < len(sample); i += 2 {
		key := sample[i].(string)
		pattern := "||" + key + "={" + key + "}"
		formatters = append(formatters, formatProperties(pattern, sample))
	}
	formatters = append(formatters, formatLiteral("\n"))
	return formatters
}
