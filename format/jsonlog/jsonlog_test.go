package jsonlog_test

import (
	"testing"
	"github.com/modern-go/test"
	"context"
	"github.com/modern-go/countlog/format/jsonlog"
	"github.com/modern-go/countlog/logger"
	"github.com/modern-go/test/must"
	"github.com/json-iterator/go"
	"errors"
	"time"
)

func TestJsonLog(t *testing.T) {
	t.Run("event", test.Case(func(ctx context.Context) {
		fmt := (&jsonlog.Format{
			Json: jsoniter.ConfigCompatibleWithStandardLibrary,
			HideTime: true,
			HideLocation: true,
			HideContext: true,
			HideProperties: true,
		}).FormatterOf(&logger.LogSite{
			Event: "event!hello",
		})
		must.Equal(`{"event":"hello"}`, string(fmt.Format(nil, &logger.Event{})))
	}))
	t.Run("callee", test.Case(func(ctx context.Context) {
		fmt := (&jsonlog.Format{
			Json: jsoniter.ConfigCompatibleWithStandardLibrary,
			HideTime: true,
			HideLocation: true,
			HideContext: true,
			HideProperties: true,
		}).FormatterOf(&logger.LogSite{
			Event: "callee!hello",
		})
		must.Equal(`{"event":"call hello"}`, string(fmt.Format(nil, &logger.Event{})))
	}))
	t.Run("msg", test.Case(func(ctx context.Context) {
		fmt := (&jsonlog.Format{
			Json: jsoniter.ConfigCompatibleWithStandardLibrary,
			HideTime: true,
			HideLocation: true,
			HideContext: true,
			HideProperties: true,
		}).FormatterOf(&logger.LogSite{
			Event: "hello",
		})
		must.Equal(`{"event":"hello"}`, string(fmt.Format(nil, &logger.Event{})))
	}))
	t.Run("error", test.Case(func(ctx context.Context) {
		fmt := (&jsonlog.Format{
			Json: jsoniter.ConfigCompatibleWithStandardLibrary,
		}).FormatterOf(&logger.LogSite{
			Event: "hello",
		})
		must.Equal(`{"error":"err","event":"hello","location":" @ :0","timestamp":0}`, string(fmt.Format(nil, &logger.Event{
			Error: errors.New("err"),
			Timestamp: time.Unix(0, 0).In(time.UTC),
		})))
	}))
	t.Run("context", test.Case(func(ctx context.Context) {
		ctx = logger.WithLogContext(ctx)
		logger.GetLogContext(ctx).Add("thread", 100)
		fmt := (&jsonlog.Format{
			Json: jsoniter.ConfigCompatibleWithStandardLibrary,
			HideTime: true,
			HideLocation: true,
		}).FormatterOf(&logger.LogSite{
			Event: "hello",
			Context: ctx,
		})
		must.Equal(`{"event":"hello","thread":100}`, string(fmt.Format(nil, &logger.Event{
			Context: ctx,
		})))
	}))
	t.Run("context", test.Case(func(ctx context.Context) {
		ctx = logger.WithLogContext(ctx)
		logger.GetLogContext(ctx).Add("thread", 100)
		fmt := (&jsonlog.Format{
			Json: jsoniter.ConfigCompatibleWithStandardLibrary,
			HideTime: true,
			HideLocation: true,
		}).FormatterOf(&logger.LogSite{
			Event: "hello",
			Context: ctx,
		})
		must.Equal(`{"event":"hello","thread":100}`, string(fmt.Format(nil, &logger.Event{
			Context: ctx,
		})))
	}))
	t.Run("properties", test.Case(func(ctx context.Context) {
		ctx = logger.WithLogContext(ctx)
		logger.GetLogContext(ctx).Add("thread", 100)
		fmt := (&jsonlog.Format{
			Json: jsoniter.ConfigCompatibleWithStandardLibrary,
			HideTime: true,
			HideLocation: true,
		}).FormatterOf(&logger.LogSite{
			Event: "hello",
		})
		must.Equal(`{"a":"b","event":"hello"}`, string(fmt.Format(nil, &logger.Event{
			Properties: []interface{}{"a", "b"},
		})))
	}))
}
