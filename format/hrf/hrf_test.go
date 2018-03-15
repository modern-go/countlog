package hrf_test

import (
	"testing"
	"github.com/modern-go/test"
	"context"
	"github.com/modern-go/countlog/format/hrf"
	"github.com/modern-go/countlog/logger"
	"github.com/modern-go/test/must"
)

func TestHrf(t *testing.T) {
	t.Run("event", test.Case(func(ctx context.Context) {
		fmt := (&hrf.Format{}).FormatterOf(&logger.LogSite{
			Event: "event!hello",
		})
		must.Equal("hello\n", string(fmt.Format(nil, &logger.Event{})))
	}))
	t.Run("callee", test.Case(func(ctx context.Context) {
		fmt := (&hrf.Format{}).FormatterOf(&logger.LogSite{
			Event: "callee!hello",
		})
		must.Equal("call hello\n", string(fmt.Format(nil, &logger.Event{})))
	}))
	t.Run("msg", test.Case(func(ctx context.Context) {
		fmt := (&hrf.Format{}).FormatterOf(&logger.LogSite{
			Event: "hello {var}",
			Sample: []interface{}{"var", "world"},
		})
		must.Equal("hello world\n", string(fmt.Format(nil, &logger.Event{
			Properties: []interface{}{"var", "world"},
		})))
	}))
}