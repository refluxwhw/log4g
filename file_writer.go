package log4g

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

func newFileLogWriter(cfg fileSecCfg) *fileLogWriter {
	w := &fileLogWriter{
		filename: cfg.Filename,
		rotate:   cfg.Rotate,
		maxsize:  strToNumSuffix(cfg.Maxsize, 1024),
		maxline:  strToNumSuffix(cfg.Maxline, 1000),
		daily:    cfg.Daily,

		ch:     make(chan *formattedRecord, 32),
		wg:     sync.WaitGroup{},
		opened: false,

		file:     nil,
		currDay:  0,
		currSize: 0,
	}

	if !strings.HasSuffix(w.filename, ".log") {
		w.filename += ".log"
	}

	go w.Run()

	return w
}

type fileLogWriter struct {
	filename string
	rotate   bool
	maxsize  int64
	maxline  int64 // unused
	daily    bool

	ch     chan *formattedRecord
	wg     sync.WaitGroup
	opened bool

	file     *os.File // 文件句柄
	currDay  int      // 文件创建时间，在一年中的某天(认为日志不会存在一年都没有写一条记录)
	currSize int64    // 当前文件大小
}

func (w *fileLogWriter) Write(msg *formattedRecord) {
	w.ch <- msg
}

func (w *fileLogWriter) Close() {
	if !w.opened {
		return
	}
	close(w.ch)
	w.wg.Wait()
}

func (w *fileLogWriter) renameFile(oldpath, newpath string, idx int) bool {
	filename := ""
	for ; ; idx++ {
		if idx == 0 {
			filename = newpath
		} else {
			filename = newpath[:len(newpath)-4] + fmt.Sprintf("_%d.log", idx)
		}

		if isFileExist(filename) {
			continue
		}

		if err := os.Rename(oldpath, filename); err != nil {
			fmt.Println(err)
			return false
		}
		break
	}
	return true
}

func (w *fileLogWriter) rotateFile(now time.Time) bool {
	if w.file != nil {
		_ = w.file.Close()
		w.file = nil
	}

	// 计算当前时间是否到达更换文件时间
	if w.rotate {
		info, err := os.Stat(w.filename)
		if err == nil { // 文件存在
			if w.daily {
				// 按时间重命名文件， filename_20060102.log
				// 如果存在重名，增加后缀 filename_20060102_n.log
				mtime := info.ModTime()
				// 此文件不是为当天的文件，或文件过大
				if (mtime.Year() != now.Year() || mtime.YearDay() != now.YearDay()) ||
					(w.maxsize > 0 && info.Size() > w.maxsize) {
					basename := w.filename[:len(w.filename)-4]
					basename += fmt.Sprintf("_%s.log", mtime.Format("20060102"))
					if !w.renameFile(w.filename, basename, 0) {
						return false
					}
				}
			} else {
				// 仅需判断文件大小，直接增加后缀
				if w.maxsize > 0 && info.Size() > w.maxsize {
					if !w.renameFile(w.filename, w.filename, 1) {
						return false
					}
				}
			}
		}
	}

	f, err := os.OpenFile(w.filename, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return false
	}
	w.file = f
	w.currSize = 0
	w.currDay = now.YearDay()

	return true
}

func (w *fileLogWriter) ProcessMsg(msg *formattedRecord) bool {
	if w.file == nil {
		if !w.rotateFile(msg.Created) {
			return false
		}
	}

	if w.rotate {
		if msg.Created.YearDay() != w.currDay ||
			w.currSize+int64(len(msg.Formatted))+1 > w.maxsize {
			if !w.rotateFile(msg.Created) {
				return false
			}
		}
	}

	if _, err := fmt.Fprintln(w.file, msg.Formatted); err != nil {
		fmt.Println(err)
		w.file.Close()
		w.file = nil
		return false
	}

	w.currSize += int64(len(msg.Formatted))
	return true
}

func (w *fileLogWriter) Run() {
	defer doRecover()

	w.wg.Add(1)
	w.opened = true

	for msg := range w.ch {
		w.ProcessMsg(msg)
	}

	if w.file != nil {
		w.file.Close()
		w.file = nil
	}

	w.opened = false
	w.wg.Done()
}

//// ---------------------------------------------------
