package stats

import (
	"reflect"
	"sync"
	"time"
	"context"
	"github.com/modern-go/countlog/logger"
	"github.com/modern-go/gls"
	"github.com/modern-go/reflect2"
	"strings"
)

type Window struct {
	collector          Collector
	event              string
	dimensionElemCount int
	shards             [16]windowShard
}

type windowShard struct {
	MapMonoid
	lock *sync.Mutex
}

func newWindow(executor Executor, collector Collector, dimensionElemCount int) *Window {
	window := &Window{
		collector:          collector,
		dimensionElemCount: dimensionElemCount,
	}
	for i := 0; i < 16; i++ {
		window.shards[i] = windowShard{
			MapMonoid: MapMonoid{},
			lock:      &sync.Mutex{},
		}
	}
	executor(window.exportEverySecond)
	return window
}

func (window *Window) exportEverySecond(ctx context.Context) {
	timer := time.NewTimer(time.Second)
	for {
		select {
		case <-timer.C:
			window.Export(time.Now())
		case <-ctx.Done():
			return
		}
	}
}

func (window *Window) Mutate() (*sync.Mutex, MapMonoid) {
	shardId := gls.GoID() % 16
	shard := window.shards[shardId]
	return shard.lock, shard.MapMonoid
}

func (window *Window) Export(now time.Time) {
	for i := 0; i < 16; i++ {
		window.exportShard(now, window.shards[i])
	}
}

func (window *Window) exportShard(now time.Time, shard windowShard) {
	shard.lock.Lock()
	defer shard.lock.Unlock()
	// batch allocate memory to hold dimensions
	space := make([]string, len(shard.MapMonoid)*window.dimensionElemCount)
	arrayType := reflect.ArrayOf(window.dimensionElemCount, reflect.TypeOf(""))
	arrayType2 := reflect2.Type2(arrayType).(*reflect2.UnsafeArrayType)
	for dimensionObj, monoid := range shard.MapMonoid {
		dimension := space[:window.dimensionElemCount]
		space = space[window.dimensionElemCount:]
		pObj := arrayType2.PackEFace(reflect2.PtrOf(dimensionObj))
		for i := 0; i < window.dimensionElemCount; i++ {
			dimension[i] = *(arrayType2.GetIndex(pObj, i).(*string))
		}
		window.collector.Collect(&Point{
			Event:     window.event,
			Timestamp: now,
			Dimension: dimension,
			Value:     monoid.Export(),
		})
	}
}

type propIdx int

type dimensionExtractor func(event *logger.Event, monoid MapMonoid, createElem func() Monoid) Monoid

func newDimensionExtractor(site *logger.LogSite) (dimensionExtractor, int) {
	var dimensionElems []string
	for i := 0; i < len(site.Sample); i += 2 {
		key := site.Sample[i].(string)
		if key == "dim" {
			dimensionElems = strings.Split(site.Sample[i+1].(string), ",")
		}
	}
	indices := make([]propIdx, 0, len(dimensionElems))
	for i := 0; i < len(site.Sample); i += 2 {
		key := site.Sample[i].(string)
		for _, dimension := range dimensionElems {
			if key == dimension {
				indices = append(indices, propIdx(i))
				indices = append(indices, propIdx(i+1))
			}
		}
	}
	return func(event *logger.Event, monoid MapMonoid, createElem func() Monoid) Monoid {
		arrayType := reflect.ArrayOf(len(indices), reflect.TypeOf(""))
		arrayType2 := reflect2.Type2(arrayType).(*reflect2.UnsafeArrayType)
		dimension := arrayType2.New()
		for i, idx := range indices {
			val := event.Properties[idx].(string)
			arrayType2.SetIndex(dimension, i, &val)
		}
		dimensionObj := arrayType2.Indirect(dimension)
		elem := monoid[dimensionObj]
		if elem == nil {
			elem = createElem()
			monoid[dimensionObj] = elem
		}
		return elem
	}, len(indices)
}
