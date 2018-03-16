package countlog

import (
	"github.com/modern-go/countlog/logger"
	"unsafe"
	"github.com/modern-go/reflect2"
	"runtime/debug"
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
	ptr := reflect2.NoEscape(unsafe.Pointer(&properties))
	log(LevelTrace, event, "", nil, nil, *(*[]interface{})(ptr))
}

// TraceCall will calculate stats in TRACE level
// TraceCall will output individual log entries in TRACE_CALL level
func TraceCall(event string, err error, properties ...interface{}) error {
	if err != nil {
		ptr := reflect2.NoEscape(unsafe.Pointer(&properties))
		return log(LevelWarn, event, "call", nil, err, *(*[]interface{})(ptr))
	}
	if LevelTrace < logger.MinLevel {
		return nil
	}
	ptr := reflect2.NoEscape(unsafe.Pointer(&properties))
	log(LevelTrace, event, "call", nil, err, *(*[]interface{})(ptr))
	return nil
}

func Debug(event string, properties ...interface{}) {
	if LevelDebug < logger.MinLevel {
		return
	}
	ptr := reflect2.NoEscape(unsafe.Pointer(&properties))
	log(LevelDebug, event, "", nil, nil, *(*[]interface{})(ptr))
}

// DebugCall will calculate stats in DEBUG level
// DebugCall will output individual log entries in DEBUG_CALL level (TRACE includes DEBUG_CALL)
func DebugCall(event string, err error, properties ...interface{}) error {
	if err != nil {
		ptr := reflect2.NoEscape(unsafe.Pointer(&properties))
		return log(LevelWarn, event, "call", nil, err, *(*[]interface{})(ptr))
	}
	if LevelDebug < logger.MinLevel {
		return nil
	}
	ptr := reflect2.NoEscape(unsafe.Pointer(&properties))
	log(LevelDebug, event, "call", nil, err, *(*[]interface{})(ptr))
	return nil
}

func Info(event string, properties ...interface{}) {
	if LevelInfo < logger.MinLevel {
		return
	}
	ptr := reflect2.NoEscape(unsafe.Pointer(&properties))
	log(LevelInfo, event, "", nil, nil, *(*[]interface{})(ptr))
}

// InfoCall will calculate stats in INFO level
// InfoCall will output individual log entries in INFO_CALL level (DEBUG includes INFO_CALL)
func InfoCall(event string, err error, properties ...interface{}) error {
	if err != nil {
		ptr := reflect2.NoEscape(unsafe.Pointer(&properties))
		return log(LevelWarn, event, "call", nil, err, *(*[]interface{})(ptr))
	}
	if LevelInfo < logger.MinLevel {
		return nil
	}
	ptr := reflect2.NoEscape(unsafe.Pointer(&properties))
	log(LevelInfo, event, "call", nil, err, *(*[]interface{})(ptr))
	return nil
}

func Warn(event string, properties ...interface{}) {
	ptr := reflect2.NoEscape(unsafe.Pointer(&properties))
	log(LevelWarn, event, "", nil, nil, *(*[]interface{})(ptr))
}

func Error(event string, properties ...interface{}) {
	ptr := reflect2.NoEscape(unsafe.Pointer(&properties))
	log(LevelError, event, "", nil, nil, *(*[]interface{})(ptr))
}

func Fatal(event string, properties ...interface{}) {
	ptr := reflect2.NoEscape(unsafe.Pointer(&properties))
	log(LevelFatal, event, "", nil, nil, *(*[]interface{})(ptr))
}

func Log(level int, event string, properties ...interface{}) {
	if level < logger.MinLevel {
		return
	}
	ptr := reflect2.NoEscape(unsafe.Pointer(&properties))
	log(level, event, "", nil, nil, *(*[]interface{})(ptr))
}

func LogPanic(recovered interface{}, properties ...interface{}) interface{} {
	if recovered == nil {
		return nil
	}
	buf := debug.Stack()
	if len(properties) > 0 {
		properties = append(properties, "err", recovered, "stacktrace", string(buf))
		Fatal("event!panic", properties...)
	} else {
		Fatal("event!panic", "err", recovered, "stacktrace", string(buf))
	}
	return recovered
}
