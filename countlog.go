package countlog

import (
	"github.com/modern-go/countlog/logger"
	"unsafe"
	"github.com/modern-go/reflect2"
)

const LevelTraceCall = logger.LevelTraceCall
const LevelTrace = logger.LevelTrace
const LevelDebugCall = logger.LevelDebugCall
const LevelDebug = logger.LevelDebug
const LevelInfoCall = logger.LevelInfoCall
const LevelInfo = logger.LevelInfo
const LevelWarn = logger.LevelWarn
const LevelError = logger.LevelError
const LevelFatal = logger.LevelFatal

func SetMinLevel(level int) {
	logger.MinLevel = level
	logger.MinCallLevel = level + 5
}

func ShouldLog(level int) bool {
	return level >= logger.MinLevel
}

func Trace(event string, properties ...interface{}) {
	if LevelTrace < logger.MinLevel {
		return
	}
	ptr := unsafe.Pointer(&properties)
	ptr = reflect2.NoEscape(ptr)
	log(LevelTrace, event, "", nil, nil, *(*[]interface{})(ptr))
}
