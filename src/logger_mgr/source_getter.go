package logger_mgr

import (
	"fmt"
	"path"
	"runtime"
)

type SourceGetter interface {
	getSource() string
}

type defaultSourceGetter struct {
	skip int
}

func (f *defaultSourceGetter) getSource() string {
	_, file, line, _ := runtime.Caller(f.skip)
	file = path.Base(file)
	return fmt.Sprintf("%s:%d", file, line)
}