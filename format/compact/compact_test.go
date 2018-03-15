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
		fmt := compact.Format{}.FormatterOf(&logger.LogSite{
			Event: "event!hello",
			Func: "func",
			File: "a.go",
			Line: 100,
		})
		must.Equal("[TRACE] [0001-01-01T00:00:00Z] [func @ a.go:100] hello", string(fmt.Format(nil, &logger.Event{
			Level: logger.LevelTrace,
		})))
	}))
	t.Run("callee", test.Case(func(ctx context.Context) {
		fmt := compact.Format{}.FormatterOf(&logger.LogSite{
			Event: "callee!hello",
			Func: "func",
			File: "a.go",
			Line: 100,
		})
		must.Equal("[TRACE] [0001-01-01T00:00:00Z] [func @ a.go:100] call hello", string(fmt.Format(nil, &logger.Event{
			Level: logger.LevelTrace,
		})))
	}))
	t.Run("msg", test.Case(func(ctx context.Context) {
		fmt := compact.Format{}.FormatterOf(&logger.LogSite{
			Event: "hello {var}",
			Func: "func",
			File: "a.go",
			Line: 100,
			Sample: []interface{}{"var", "world"},
		})
		must.Equal("[TRACE] [0001-01-01T00:00:00Z] [func @ a.go:100] hello world", string(fmt.Format(nil, &logger.Event{
			Level: logger.LevelTrace,
			Properties: []interface{}{"var", "world"},
		})))
	}))
	t.Run("error", test.Case(func(ctx context.Context) {
		fmt := compact.Format{}.FormatterOf(&logger.LogSite{
			Event: "hello",
			Func: "func",
			File: "a.go",
			Line: 100,
		})
		must.Equal("[TRACE] [0001-01-01T00:00:00Z] [func @ a.go:100] hello: err", string(fmt.Format(nil, &logger.Event{
			Level: logger.LevelTrace,
			Error: errors.New("err"),
		})))
	}))
}
