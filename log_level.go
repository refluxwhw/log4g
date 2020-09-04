package log4g

import "strings"

type Level int

const (
	DEBUG Level = iota
	TRACE
	INFO
	WARNING
	ERROR
	CRITICAL
)

// Logging level strings
var (
	longLevelStrings = [...]string{"DEBUG", "TRACE", "INFO ", "WARN ", "ERROR", "CRITI"}
	//longLevelStrings = [...]string{"DEBG", "TRAC", "INFO", "WARN", "EROR", "CRIT"}
	shortLevelStrings = [...]string{"D", "T", "I", "W", "E", "C"}
)

func (l Level) LongString() string {
	if l < 0 || int(l) > len(longLevelStrings) {
		return "UNKNOWN"
	}
	return longLevelStrings[l]
}

func (l Level) ShortString() string {
	if l < 0 || int(l) > len(shortLevelStrings) {
		return "UNKNOWN"
	}
	return shortLevelStrings[l]
}

func strToLevel(s string) Level {
	s = strings.ToUpper(s)
	l := DEBUG
	switch s {
	case "DEBUG":
		l = DEBUG
	case "TRACE":
		l = TRACE
	case "INFO":
		l = INFO
	case "WARNING":
		l = WARNING
	case "ERROR":
		l = ERROR
	case "CRITICAL":
		l = CRITICAL
	}
	return l
}
