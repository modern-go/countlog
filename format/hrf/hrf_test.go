package hrf_test

import (
	"testing"
	"github.com/modern-go/test"
	"context"
	"github.com/modern-go/countlog/format/hrf"
	"github.com/modern-go/countlog/logger"
	"github.com/modern-go/test/must"
	"errors"
)

func TestHrf(t *testing.T) {
	t.Run("event", test.Case(func(ctx context.Context) {
		fmt := (&hrf.Format{
			HideLevel: true,
			HideTime:  true,
		}).FormatterOf(&logger.LogSite{
			Event: "event!hello",
		})
		must.Equal("hello\n", string(fmt.Format(nil, &logger.Event{})))
	}))
	t.Run("callee", test.Case(func(ctx context.Context) {
		fmt := (&hrf.Format{
			HideLevel: true,
			HideTime:  true,
		}).FormatterOf(&logger.LogSite{
			Event: "callee!hello",
		})
		must.Equal("call hello\n", string(fmt.Format(nil, &logger.Event{})))
	}))
	t.Run("msg", test.Case(func(ctx context.Context) {
		fmt := (&hrf.Format{
			HideLevel: true,
			HideTime:  true,
		}).FormatterOf(&logger.LogSite{
			Event:  "hello {var}",
			Sample: []interface{}{"var", "world"},
		})
		must.Equal("hello world\n", string(fmt.Format(nil, &logger.Event{
			Properties: []interface{}{"var", "world"},
		})))
	}))
	t.Run("error", test.Case(func(ctx context.Context) {
		fmt := (&hrf.Format{
		}).FormatterOf(&logger.LogSite{
			Event: "hello",
		})
		must.Equal("\x1b[31;1m[ERROR]\x1b[0m hello\n"+
			"\x1b[90;1merror: err\x1b[0m\n"+
			"\x1b[90;1mtimestamp: 0001-01-01T00:00:00Z\x1b[0m\n"+
			"\x1b[90;1mlocation:  @ :0\x1b[0m\n", string(fmt.Format(nil, &logger.Event{
			Level: logger.LevelError,
			Error: errors.New("err"),
		})))
	}))
	t.Run("context", test.Case(func(ctx context.Context) {
		ctx = logger.WithContext(ctx)
		logger.AddLogContext(ctx, "thread", 100)
		fmt := (&hrf.Format{
			HideLevel: true,
			HideTime:  true,
		}).FormatterOf(&logger.LogSite{
			Event:   "hello",
			Context: ctx,
		})
		must.Equal("hello\n"+
			"\x1b[90;1mthread: 100\x1b[0m\n", string(fmt.Format(nil, &logger.Event{
			Context: ctx,
		})))
	}))
}
