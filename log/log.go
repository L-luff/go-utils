package log

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile)
	infoLog  = log.New(os.Stdout, "\033[34m[info ]\033[0m ", log.LstdFlags|log.Lshortfile)
	debugLog = log.New(os.Stdout, "\033[34m[debug ]\033[0m ", log.LstdFlags|log.Lshortfile)
	loggers  = []*log.Logger{errorLog, debugLog, infoLog}
	mu       sync.Mutex
)

// log methods
var (
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
	Debug  = debugLog.Println
	Debugf = debugLog.Printf
)

const (
	debugLevel = iota
	InfoLevel
	ErrorLevel
	Disabled
)

func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()
	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}
	if debugLevel < level {
		debugLog.SetOutput(ioutil.Discard)
	}
	if InfoLevel < level {
		infoLog.SetOutput(ioutil.Discard)
	}
	if ErrorLevel < level {
		errorLog.SetOutput(ioutil.Discard)
	}

}
