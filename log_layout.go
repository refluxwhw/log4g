package log4g

import (
	"bytes"
	"fmt"
	"strings"
)

type sectionType int

const (
	kMsg sectionType = iota
	kLongLevel
	kShortLevel
	kSource
	kCategory
	kTime
	kString
	kGoroutineID
)

const (
	kDefaultLayout         = "[%T] %L %C (%S) %M"
 	kDefaultDatetimeFormat = "2006-01-02 15:04:05.000"
)

var gDefaultLayout = newLayoutConf(kDefaultLayout)

type section struct {
	T sectionType
	V string
}

type layoutInfo struct {
	Sections []section
}

/**
 * TODO:
 * %{section:format} 	section is required, must be one of (time level category file line function thread message),
 * 						format is optional
 *
 * - time		: datetime, format must be valid datetime format,
 *					default format is '2006-01-02 15:04:05.000'
 * - level      : log level, format is '<long|short>', short is (D, T, W, E, C), long is (DEBG, TRAC, WARN, EROR, CRIT),
 *              	default format is 'short'
 * - category   : log category, format is SPrintf string format, eg. '%10s', '%-10s',
 *					default format is '%s'
 * - file       : source code file name, format is '<long|short>[;string-format]'
 * 					default format is 'short;%s'
 * - line       : source code line, format is SPrintf int format
 *					default format is '%d'
 * - function   : source code function name, format is '<long|short>[;string-format]'
 *					default format is 'short;%s'
 * - thread     : goroutines id, format is ignored
 * - message    : output message, format is ignored
 *
 * Ignores unknown formats
 * Recommended: "[%{time}] %{level} %{category} (%{file}:%{line}) %{message}"
 */

// Known format codes:
// %T - DateTime with format string, default format is '2006-01-03 15:04:05.000'
//  	eg. %T{2006-01-02 15:04:05}
// %L - Level (DEBG, TRAC, WARN, EROR, CRIT)
// %l - Level (D, T, W, E, C)
// %C - category
// %S - Source filename:line
// %G - GoroutineID
// %M - Message
// Ignores unknown formats
// Recommended: "[%T] %L %C (%S) %M"
func newLayoutConf(layout string) *layoutInfo {
	lc := &layoutInfo{}

	// 1. 找到 '%' ，将字符串分割为前后两段， prefix / suffix
	// 2. 如果 suffix 不存在，即为没有找到 '%'，或者为空，即为 '%' 在字符串结尾，将整体作为一个字符串追加到之前没有处理的字符串上，
	//	并添加一个section，然后结束
	// 3. 否则，prefix 追加到之前没有处理的字符串上，suffix 将第一个字符 ch 分割出来
	// 4. 判断 ch
	// 		T: 时间，后半段解析时间格式，
	//			判断第一个字符是否为 '{' 并且能找到 '}'，
	//			如果是，那么 '{}' 括号中的内容认为是时间格式化字符串，后半段'}' 之后部分为
	//			如果不是，设置默认格式化字符串
	//		L: 日志等级
	//		C: 日志分类
	//		S: 日志源
	//		G: go协程ID
	//		M: 日志信息
	//		%: 将字符 '%' 追加到 prefix 上
	//		其他: 将 '%' 和 ch 追加到 prefix 上

	prefix := ""
	suffix := layout
	add := func(t sectionType, v string) {
		if len(prefix) > 0 {
			lc.Sections = append(lc.Sections, section{T: kString, V: prefix})
			prefix = ""
		}
		lc.Sections = append(lc.Sections, section{T: t, V: v})
	}

	for {
		if len(suffix) == 0 {
			if len(prefix) > 0 {
				lc.Sections = append(lc.Sections, section{T: kString, V: prefix})
			}
			break
		}

		sl := strings.SplitN(suffix, "%", 2)
		if len(sl) == 1 || sl[1] == "" {
			prefix += suffix
			lc.Sections = append(lc.Sections, section{T: kString, V: prefix})
			break
		}

		prefix += sl[0]
		suffix = sl[1]

		ch := suffix[0]
		suffix = suffix[1:]
		switch ch {
		case 'T':
			format := ""
			format, suffix = getTimeFormat(suffix)
			add(kTime, format)
		case 'L':
			add(kLongLevel, "")
		case 'l':
			add(kShortLevel, "")
		case 'C':
			add(kCategory, "")
		case 'S':
			add(kSource, "")
		case 'G':
			add(kGoroutineID, "")
		case 'M':
			add(kMsg, "")
		case '%':
			prefix += "%"
		default:
			prefix += "%" + string(ch)
		}
	}

	return lc
}

func getTimeFormat(layout string) (format string, newLayout string) {
	format = kDefaultDatetimeFormat
	newLayout = layout
	if len(newLayout) < 2 || newLayout[0] != '{' {
		return
	}

	idx := strings.Index(newLayout, "}")
	if idx < 0 {
		return
	}

	format = newLayout[1:idx]
	newLayout = newLayout[idx+1:]

	return
}

func recordFormatToString(rec *logRecord, layout *layoutInfo) string {
	if rec == nil || layout == nil {
		return "<nil>"
	}

	out := bytes.NewBuffer(make([]byte, 0, 64))

	for i := 0; i < len(layout.Sections); i++ {
		switch layout.Sections[i].T {
		case kTime:
			out.WriteString(rec.Created.Format(layout.Sections[i].V))
		case kCategory:
			out.WriteString(rec.Category)
		case kLongLevel:
			out.WriteString(rec.Level.LongString())
		case kShortLevel:
			out.WriteString(rec.Level.ShortString())
		case kSource:
			out.WriteString(fmt.Sprintf("%s:%d", rec.Source.File, rec.Source.Line))
		case kGoroutineID:
			out.WriteString(fmt.Sprintf("%016x", rec.Source.Tid))
		case kMsg:
			out.WriteString(rec.Message)
		case kString:
			out.WriteString(layout.Sections[i].V)
		}
	}

	return out.String()
}
