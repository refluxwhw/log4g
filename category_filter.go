package log4g

import "fmt"

func newFilter(category string, level Level, layout *layoutInfo) *categoryFilter {
	filter := &categoryFilter{
		category: category,
		level:    level,
		writers:  make([]logWriter, 0),
		layout:   layout,
	}
	return filter
}

type categoryFilter struct {
	category string
	level    Level
	writers  []logWriter
	layout   *layoutInfo
}

func (f *categoryFilter) logMessage(rec *logRecord) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	if rec.Level < f.level {
		return
	}

	fr := &formattedRecord{
		Created: rec.Created,
		Formatted: recordFormatToString(rec, f.layout),
	}

	for _, writer := range f.writers {
		writer.Write(fr)
	}
}
