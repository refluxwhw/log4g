package log4g

import (
	"time"
)

type logSource struct {
	Tid  uint64
	File string
	Func string
	Line int
}

// A logRecord contains all of the pertinent information for each message
type logRecord struct {
	Category string     // The log group
	Level    Level      // The log level
	Created  time.Time  // The time at which the log message was created (nanoseconds)
	Message  string     // The log message
	Source   *logSource // The message source
}

type formattedRecord struct {
	Created   time.Time // 时间
	Formatted string    // 格式化后的字符串，没有附件换行符
}
