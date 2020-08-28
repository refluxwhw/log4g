package logger_mgr

import (
	"fmt"
	"runtime"
	"strings"

	"third_part/log4go"
)

type Level int

const (
	FATAL Level = iota
	ERROR
	WARN
	INFO
	DEBUG
	TRACE
)

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
	Logf(level Level, format string, args ...interface{})
}

type logger struct {
	log    *log4go.Filter
	source SourceGetter
}

func newLoggerWithDefaultSourceGetter(category string, skip int) *logger {
	return &logger{
		log: log4go.LOGGER(category),
		source: &defaultSourceGetter{
			skip: skip,
		},
	}
}

func newLoggerWithSourceGetter(category string, source SourceGetter) *logger {
	return &logger{
		log:    log4go.LOGGER(category),
		source: source,
	}
}

func loadConfig(path string)  {
	log4go.LoadConfiguration(path)
}

func closeLog4go()  {
	log4go.Close()
}

/// ----------------------------------------------------------------------

func formatter(args ...interface{}) string {
	msg := fmt.Sprintf(strings.Repeat("%v", len(args)), args...)
	return msg
}

func getCallStack() string {
	buf := make([]byte, 1024*4)
	n := runtime.Stack(buf, false)
	return string(buf[0:n])
}

func (l *logger) Critical(args ...interface{}) {
	l.log.Log(log4go.CRITICAL, l.source.getSource(), formatter(args...))
}
func (l *logger) CriticalF(format string, args ...interface{}) {
	l.log.Log(log4go.CRITICAL, l.source.getSource(), fmt.Sprintf(format, args...))
}

func (l *logger) Error(args ...interface{}) {
	l.log.Log(log4go.ERROR, l.source.getSource(), formatter(args...))
}
func (l *logger) ErrorF(format string, args ...interface{}) {
	l.log.Log(log4go.ERROR, l.source.getSource(), fmt.Sprintf(format, args...))
}

func (l *logger) Warn(args ...interface{}) {
	l.log.Log(log4go.WARNING, l.source.getSource(), formatter(args...))
}
func (l *logger) WarnF(format string, args ...interface{}) {
	l.log.Log(log4go.WARNING, l.source.getSource(), fmt.Sprintf(format, args...))
}

func (l *logger) Info(args ...interface{}) {
	l.log.Log(log4go.INFO, l.source.getSource(), formatter(args...))
}
func (l *logger) InfoF(format string, args ...interface{}) {
	l.log.Log(log4go.INFO, l.source.getSource(), fmt.Sprintf(format, args...))
}

func (l *logger) Debug(args ...interface{}) {
	l.log.Log(log4go.DEBUG, l.source.getSource(), formatter(args...))
}
func (l *logger) DebugF(format string, args ...interface{}) {
	l.log.Log(log4go.DEBUG, l.source.getSource(), fmt.Sprintf(format, args...))
}

func (l *logger) Trace(args ...interface{}) {
	l.log.Log(log4go.TRACE, l.source.getSource(), formatter(args...)+"\n"+getCallStack())
}
func (l *logger) TraceF(format string, args ...interface{}) {
	l.log.Log(log4go.TRACE, l.source.getSource(), fmt.Sprintf(format, args...)+"\n"+getCallStack())
}

func (l *logger) Log(level Level, args ...interface{}) {
	source := l.source.getSource()
	msg := fmt.Sprint(args...)
	l.log.Log(log4go.Level(level), source, msg)
}
func (l *logger) Logf(level Level, format string, args ...interface{}) {
	source := l.source.getSource()
	msg := fmt.Sprintf(format, args...)
	l.log.Log(log4go.Level(level), source, msg)
}
