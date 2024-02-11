package golib

import (
	"fmt"
	"log"
	"os"
)

const (
	logFlag = log.Ldate | log.Ltime | log.Lshortfile | log.Lmsgprefix
)

var (
	stdLogger *log.Logger
)

// GetStdLogger ...
func GetStdLogger() *log.Logger {
	if stdLogger != nil {
		return stdLogger
	}
	return NewStdLogger("")
}

// NewStdLogger ...
func NewStdLogger(prefix string) *log.Logger {
	return log.New(os.Stdout, prefix, logFlag)
}

// NewThirdPartyLogger third party logger
func NewThirdPartyLogger(key string) *log.Logger {
	prefix := fmt.Sprintf("[third party] [%s] ", key)
	l := log.New(os.Stdout, prefix, logFlag)
	return l
}
