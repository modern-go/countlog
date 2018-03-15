package jsonlog

import (
	"github.com/modern-go/countlog/logger"
	"github.com/modern-go/countlog/format"
	"github.com/json-iterator/go"
	"strings"
)

type Format struct {
	Json           jsoniter.API
	HideLocation   bool
	HideTime       bool
	HideContext    bool
	HideProperties bool
}

func (f *Format) FormatterOf(site *logger.LogSite) format.Formatter {
	json := f.Json
	if json == nil {
		json = jsoniter.ConfigDefault
	}
	baseMap := map[interface{}]interface{}{}
	if strings.HasPrefix(site.Event, "event!") {
		baseMap["event"] = site.Event[len("event!"):]
	} else if strings.HasPrefix(site.Event, "callee!") {
		baseMap["event"] = "call " + site.Event[len("callee!"):]
	} else {
		baseMap["event"] = site.Event
	}
	if !f.HideLocation {
		baseMap["location"] = site.Location()
	}
	return format.FuncFormatter(func(space []byte, event *logger.Event) []byte {
		eventMap := map[interface{}]interface{}{}
		for k, v := range baseMap {
			eventMap[k] = v
		}
		f.fillEventMap(eventMap, event)
		output, err := json.Marshal(eventMap)
		if err != nil {
			logger.ErrorLogger.Println("failed to marshal json", err)
			return nil
		}
		if space == nil {
			return output
		}
		return append(space, output...)
	})
}

func (f *Format) fillEventMap(eventMap map[interface{}]interface{}, event *logger.Event) {
	if event.Error != nil {
		eventMap["error"] = event.Error.Error()
	}
	if !f.HideTime {
		eventMap["timestamp"] = event.Timestamp.UnixNano()
	}
	ctx := logger.GetLogContext(event.Context)
	if !f.HideContext && ctx != nil {
		for i := 0; i < len(ctx.Properties); i+=2 {
			eventMap[ctx.Properties[i]] = ctx.Properties[i+1]
		}
	}
	if !f.HideProperties {
		for i := 0; i < len(event.Properties); i+=2 {
			eventMap[event.Properties[i]] = event.Properties[i+1]
		}
	}
}
