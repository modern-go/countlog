package countlog

import (
	"os"
	"github.com/modern-go/countlog/logger"
	"github.com/modern-go/countlog/output"
	"github.com/modern-go/countlog/format/hrf"
	"github.com/modern-go/countlog/stats"
)

var EventWriter logger.EventSink = output.NewEventWriter(func(eventWriter *output.EventWriter) {
	eventWriter.Format = &hrf.Format{}
	eventWriter.Writer = os.Stdout
})

var EventAggregator logger.EventSink = stats.NewEventAggregator(func(aggregator *stats.EventAggregator) {
	aggregator.Collector = nil
})