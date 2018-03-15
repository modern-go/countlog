package countlog

import (
	"github.com/modern-go/countlog/logger"
)

type panicHandler func(recovered interface{}, event *logger.Event, site *logger.LogSite)

func newRootHandler(site *logger.LogSite, onPanic panicHandler) logger.EventHandler {
	return nil
}
