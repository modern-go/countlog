package stats

import (
	"testing"
	"time"
	"github.com/modern-go/test"
	"context"
	"github.com/modern-go/countlog/logger"
	"github.com/modern-go/test/must"
)

func TestStats(t *testing.T) {
	t.Run("counter", test.Case(func(ctx context.Context) {
		dumpPoints := &dumpPoint{}
		aggregator := NewEventAggregator(func(aggregator *EventAggregator) {
			aggregator.Collector = dumpPoints
		})
		counter := aggregator.HandlerOf(&logger.LogSite{
			Event: "event!abc",
			Sample: []interface{}{
				"agg", "counter",
				"dim", "city,ver",
				"city", "beijing",
				"ver", "1.0",
			},
		}).(State)
		counter.Handle(&logger.Event{
			Properties: []interface{}{
				"agg", "counter",
				"dim", "city,ver",
				"city", "beijing",
				"ver", "1.0",
			},
		})
		counter.Handle(&logger.Event{
			Properties: []interface{}{
				"agg", "counter",
				"dim", "city,ver",
				"city", "beijing",
				"ver", "1.0",
			},
		})
		window := counter.GetWindow()
		window.Export(time.Now())
		points := *dumpPoints
		must.Equal(1, len(points))
		must.Equal(float64(2), points[0].Value)
		must.Equal([]string{"city", "beijing", "ver", "1.0"}, points[0].Dimension)
	}))
}

type dumpPoint []*Point

func (points *dumpPoint) Collect(point *Point) {
	*points = append(*points, point)
}
