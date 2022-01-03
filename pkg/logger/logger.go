package logger

import (
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger
)

func init() {
	flag := log.Ldate | log.Ltime | log.Lshortfile

	InfoLogger = log.New(os.Stdout, "INFO: ", flag)
	WarnLogger = log.New(os.Stderr, "WARN: ", flag)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", flag)
}

func Info(str string) {
	InfoLogger.Println(str)
}

func Warn(str string) {
	WarnLogger.Println(str)
}

func Error(str string) {
	ErrorLogger.Println(str)
}
