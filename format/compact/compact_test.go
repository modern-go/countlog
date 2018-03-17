package compact_test

import (
	"testing"
	"github.com/modern-go/test"
	"context"
	"github.com/modern-go/countlog/format/compact"
	"github.com/modern-go/countlog/logger"
	"github.com/modern-go/test/must"
	"errors"
)

func TestCompact(t *testing.T) {
	t.Run("event", test.Case(func(ctx context.Context) {
		fmt := (&compact.Format{
			HideLevel:    true,
			HideTime:     true,
			HideLocation: true,
		}).FormatterOf(&logger.LogSite{
			Event: "event!hello",
		})
		must.Equal("hello\n", string(fmt.Format(nil, &logger.Event{
		})))
	}))
	t.Run("callee", test.Case(func(ctx context.Context) {
		fmt := (&compact.Format{
			HideLevel:    true,
			HideTime:     true,
			HideLocation: true,
		}).FormatterOf(&logger.LogSite{
			Event: "callee!hello",
		})
		must.Equal("call hello\n", string(fmt.Format(nil, &logger.Event{
		})))
	}))
	t.Run("msg", test.Case(func(ctx context.Context) {
		fmt := (&compact.Format{
			HideLevel:    true,
			HideTime:     true,
			HideLocation: true,
		}).FormatterOf(&logger.LogSite{
			Event:  "hello {var}",
			Sample: []interface{}{"var", "world"},
		})
		must.Equal("hello world||var=world\n", string(fmt.Format(nil, &logger.Event{
			Properties: []interface{}{"var", "world"},
		})))
	}))
	t.Run("error", test.Case(func(ctx context.Context) {
		fmt := (&compact.Format{}).FormatterOf(&logger.LogSite{
			Event: "hello",
			Func:  "func",
			File:  "a.go",
			Line:  100,
		})
		must.Equal("[TRACE] [0001-01-01T00:00:00Z] [func @ a.go:100] hello: err\n", string(fmt.Format(nil, &logger.Event{
			Level: logger.LevelTrace,
			Error: errors.New("err"),
		})))
	}))
	t.Run("context", test.Case(func(ctx context.Context) {
		ctx = logger.WithLogContext(ctx)
		logger.GetLogContext(ctx).Add("thread", 123)
		fmt := (&compact.Format{
			HideLevel:    true,
			HideTime:     true,
			HideLocation: true,
		}).FormatterOf(&logger.LogSite{
			Event:   "hello",
			Context: ctx,
		})
		must.Equal("hello||thread=123\n", string(fmt.Format(nil, &logger.Event{
			Context: ctx,
		})))
	}))
}
