package countlog_test

import (
	"testing"
	"github.com/modern-go/test"
	"context"
	"github.com/modern-go/countlog"
	"github.com/modern-go/countlog/output"
	"bytes"
	"github.com/modern-go/countlog/format/jsonlog"
	"github.com/json-iterator/go"
	"github.com/modern-go/test/must"
)

func TestCountlog(t *testing.T) {
	t.Run("happy path", test.Case(func(ctx context.Context) {
		buf := bytes.NewBuffer(nil)
		countlog.SetMinLevel(countlog.LevelTrace)
		countlog.EventWriter = output.NewEventWriter(func(eventWriter *output.EventWriter) {
			eventWriter.Writer = buf
			eventWriter.Format = &jsonlog.Format{}
		})
		countlog.Trace("hello", "key", "value")
		must.Equal("hello", jsoniter.Get(buf.Bytes(), "event").ToString())
	}))
}