package log4g

import (
	"bytes"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
)

func isFileExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// Parse a number with K/M/G suffixes based on thousands (1000) or 2^10 (1024)
func strToNumSuffix(str string, mult int) int64 {
	num := 1
	if len(str) > 1 {
		switch str[len(str)-1] {
		case 'G', 'g':
			num *= mult
			fallthrough
		case 'M', 'm':
			num *= mult
			fallthrough
		case 'K', 'k':
			num *= mult
			str = str[0 : len(str)-1]
		}
	}
	parsed, _ := strconv.Atoi(str)
	return int64(parsed) * int64(num)
}

func getGoroutineID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func getSource(skip int) *logSource {
	src := &logSource{}
	src.Tid = getGoroutineID()

	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return src
	}

	funcName := runtime.FuncForPC(pc).Name()
	if idx := strings.LastIndex(funcName, "."); idx >= 0 {
		funcName = funcName[idx+1:]
	}

	src.File = path.Base(file)
	src.Func = funcName
	src.Line = line

	return src
}
