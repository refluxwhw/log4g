package log4g

import "fmt"

/**
 * 输出接口，可以输出到文件、控制台、网络 ...
 */
type logWriter interface {
	Write(msg *formattedRecord)
	Close()
}


func doRecover() {
	if err := recover(); err != nil {
		fmt.Println("recover:", err)
	}
}