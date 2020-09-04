package log4g

/**
 * logger interface
 */
type Logger interface {
	Critical(args ...interface{})
	CriticalF(format string, args ...interface{})

	Error(args ...interface{})
	ErrorF(format string, args ...interface{})

	Warn(args ...interface{})
	WarnF(format string, args ...interface{})

	Info(args ...interface{})
	InfoF(format string, args ...interface{})

	Debug(args ...interface{})
	DebugF(format string, args ...interface{})

	Trace(args ...interface{})
	TraceF(format string, args ...interface{})

	Log(level Level, args ...interface{})
	LogF(level Level, format string, args ...interface{})
}
