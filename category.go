package log4g

import (
	"fmt"
	"time"
)

func newDefaultCategory(name string) Logger {
	return &category{
		category: name,
		extSkip:  0,
		filters: []*categoryFilter{
			{
				category: name,
				level:    DEBUG,
				writers: []logWriter{
					gSingleConsoleWriter,
				},
				layout: gDefaultLayout,
			},
		},
	}
}

/**
 * 日志分类，继承自 Logger ，并包含多个过滤器
 */
type category struct {
	category string            // 分类名称
	extSkip  int               // 如果外部增加调用层级，就需要使用这个值来修正源代码位置
	filters  []*categoryFilter // 过滤器，日志输出最终会进入到过滤器中，然后过滤器再送到各个输出对象中
}

func (c *category) Critical(args ...interface{}) {
	c.internalLog(CRITICAL, fmt.Sprint(args...))
}
func (c *category) CriticalF(format string, args ...interface{}) {
	c.internalLog(CRITICAL, fmt.Sprintf(format, args...))
}

func (c *category) Error(args ...interface{}) {
	c.internalLog(ERROR, fmt.Sprint(args...))
}
func (c *category) ErrorF(format string, args ...interface{}) {
	c.internalLog(ERROR, fmt.Sprintf(format, args...))
}

func (c *category) Warn(args ...interface{}) {
	c.internalLog(WARNING, fmt.Sprint(args...))
}
func (c *category) WarnF(format string, args ...interface{}) {
	c.internalLog(WARNING, fmt.Sprintf(format, args...))
}

func (c *category) Info(args ...interface{}) {
	c.internalLog(INFO, fmt.Sprint(args...))
}
func (c *category) InfoF(format string, args ...interface{}) {
	c.internalLog(INFO, fmt.Sprintf(format, args...))
}

func (c *category) Debug(args ...interface{}) {
	c.internalLog(DEBUG, fmt.Sprint(args...))
}
func (c *category) DebugF(format string, args ...interface{}) {
	c.internalLog(DEBUG, fmt.Sprintf(format, args...))
}

func (c *category) Trace(args ...interface{}) {
	c.internalLog(TRACE, fmt.Sprint(args...))
}
func (c *category) TraceF(format string, args ...interface{}) {
	c.internalLog(TRACE, fmt.Sprintf(format, args...))
}

func (c *category) Log(level Level, args ...interface{}) {
	c.internalLog(level, fmt.Sprint(args...))
}
func (c *category) LogF(level Level, format string, args ...interface{}) {
	c.internalLog(level, fmt.Sprintf(format, args...))
}

func (c *category) internalLog(level Level, msg string) {
	rec := &logRecord{
		Category: c.category,
		Level:    level,
		Created:  time.Now(),
		Message:  msg,
		Source:   getSource(3 + c.extSkip),
	}

	for _, filter := range c.filters {
		filter.logMessage(rec)
	}
}

func (c *category) internalAddFilter(cfg categoryFilterCfg, writers map[string]logWriter,
	layouts map[string]*layoutInfo) error {
	layout, ok := layouts[cfg.Layout]
	if !ok {
		return fmt.Errorf("layout not found in layouts config")
	}

	filter := newFilter(c.category, strToLevel(cfg.Level), layout)
	for _, output := range cfg.Output {
		var writer logWriter
		if output == "console" {
			writer = gSingleConsoleWriter
		} else {
			var ok bool
			writer, ok = writers[output]
			if !ok {
				return fmt.Errorf("output not found in files config")
			}
		}
		filter.writers = append(filter.writers, writer)
	}

	c.filters = append(c.filters, filter)
	return nil
}
