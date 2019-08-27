package utils

import (
	"fmt"
	"log"
	"os"
)

const TOKEN = "[nginx-reloader] "

var stdoutLogger = log.New(os.Stdout, "", 0)
var stderrLogger = log.New(os.Stderr, "", 0)

func Fatalf(format string, v ...interface{}) {
	stderrLogger.Fatalf(TOKEN+"[fatal] "+format, v...)
}

func Stdoutf(format string, v ...interface{}) {
	stdoutLogger.Printf(TOKEN+format, v...)
}

// func Stderrf(format string, v ...interface{}) {
//     stderrLogger.Printf(TOKEN+"[error] "+format, v...)
// }

func Panicf(format string, v ...interface{}) {
	panic(fmt.Sprintf(format, v...))
}
