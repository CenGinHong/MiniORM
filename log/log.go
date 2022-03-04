package log

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

// 日志分级
// 不同层级日志使用不同颜色
// 显示打印日志代码对应的文件名和行号

var (
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m", log.LstdFlags|log.Lshortfile) // 使用红色作为logger
	infoLog  = log.New(os.Stdout, "\033[34m[info ]\033[0m", log.LstdFlags|log.Lshortfile) // 使用lanselogger
	loggers  = []*log.Logger{errorLog, infoLog}
	mu       sync.Mutex
)

var (
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
)

const (
	InfoLevel = iota
	ErrorLevel
	Disable
)

func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()
	// 将所有logger的输出设为标准输出流
	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}
	// 设置为ErrorLevel和infoLog的输出会被定向到ioutil.Discard,即不打印
	if ErrorLevel < level {
		errorLog.SetOutput(ioutil.Discard)
	}
	if InfoLevel < level {
		infoLog.SetOutput(ioutil.Discard)
	}
}
