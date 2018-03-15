package stats

import (
	"github.com/modern-go/countlog/logger"
)

type countEvent struct {
	*Window
	site    *logger.LogSite
	extract dimensionExtractor
}

func (state *countEvent) Handle(event *logger.Event) {
	lock, dimensions := state.Window.Mutate()
	lock.Lock()
	counter := state.extract(event, dimensions, NewCounterMonoid)
	*(counter.(*CounterMonoid)) += CounterMonoid(1)
	lock.Unlock()
}

func (state *countEvent) GetWindow() *Window {
	return state.Window
}

func (state *countEvent) LogSite() *logger.LogSite {
	return state.site
}
