package stats

import (
	"github.com/modern-go/countlog/logger"
)

type EventAggregator struct {
	Executor  Executor
	Collector Collector
}

func NewEventAggregator(init func(aggregator *EventAggregator)) *EventAggregator {
	aggregator := &EventAggregator{
		Executor: DefaultExecutor,
	}
	init(aggregator)
	return aggregator
}

func (aggregator *EventAggregator) HandlerOf(site *logger.LogSite) logger.EventHandler {
	if site.Agg != "" {
		return aggregator.createHandler(site.Agg, site)
	}
	sample := site.Sample
	for i := 0; i < len(sample); i += 2 {
		if sample[i].(string) == "agg" {
			return aggregator.createHandler(sample[i+1].(string), site)
		}
	}
	return nil
}

func (aggregator *EventAggregator) createHandler(agg string, site *logger.LogSite) logger.EventHandler {
	if aggregator.Collector == nil {
		// disable aggregation if collector not set
		return &logger.DummyEventHandler{Site: site}
	}
	extractor, dimensionElemCount := newDimensionExtractor(site)
	window := newWindow(aggregator.Executor, aggregator.Collector, dimensionElemCount)
	switch agg {
	case "counter":
		return &countEvent{
			Window:  window,
			extract: extractor,
			site:    site,
		}
	default:
		// TODO: log unknown agg
	}
	return nil
}
