package printf_test

import (
	"testing"
	"github.com/modern-go/test"
	"context"
	"github.com/modern-go/test/must"
	"github.com/modern-go/countlog/logger"
	"github.com/modern-go/countlog/format/printf"
)

func TestPrintf(t *testing.T) {
	t.Run("printf", test.Case(func(ctx context.Context) {
		format := &printf.Format{"{message}"}
		formatter := format.FormatterOf(&logger.LogSite{
			Event:  "hello {key}",
			Sample: []interface{}{"key", "world"},
		})
		output := formatter.Format(nil, &logger.Event{
			Properties: []interface{}{"key", "world"},
		})
		must.Equal("hello world\n", string(output))
	}))
}
