package log4g

import (
	"fmt"
	"os"
	"sync"
)

// 不会被关闭的
var gSingleConsoleWriter = newConsoleLogWriter()

func newConsoleLogWriter() *consoleLogWriter {
	writer := &consoleLogWriter{
		ch: make(chan *formattedRecord, 16),
		wg: sync.WaitGroup{},
	}
	go writer.Run()
	return writer
}

type consoleLogWriter struct {
	ch chan *formattedRecord
	wg sync.WaitGroup
	open bool
}

func (w *consoleLogWriter) Write(msg *formattedRecord) {
	w.ch <- msg
}

func (w *consoleLogWriter) Close() {
	if !w.open {
		return
	}
	close(w.ch)
	w.wg.Wait()
}

func (w *consoleLogWriter) Run() {
	defer doRecover()

	w.wg.Add(1)
	w.open = true

	for rec := range w.ch {
		fmt.Fprintln(os.Stdout, rec.Formatted)
	}

	w.open = false
	w.wg.Done()
}
