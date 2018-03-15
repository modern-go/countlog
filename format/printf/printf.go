package printf

import (
	"time"
	"github.com/modern-go/msgfmt/formatter"
	"github.com/modern-go/countlog/logger"
	"github.com/modern-go/countlog/format"
)

type Format struct {
	Layout string
}

func (f *Format) FormatterOf(site *logger.LogSite) format.Formatter {
	logFmt := formatter.Of(f.Layout+"\n",
		[]interface{}{
			"message", []byte{},
			"timestamp", time.Time{},
			"level", "",
			"event", "",
			"func", "",
			"file", "",
			"line", 0,
		})
	messageFmt := formatter.Of(site.Event, site.Sample)
	return format.FuncFormatter(func(space []byte, event *logger.Event) []byte {
		return logFmt.Format(space,
			[]interface{}{
				"message", messageFmt.Format(nil, event.Properties),
				"timestamp", event.Timestamp,
				"level", logger.LevelName(event.Level),
				"event", site.Event,
				"func", site.Func,
				"file", site.File,
				"line", site.Line,
			})
	})
}
