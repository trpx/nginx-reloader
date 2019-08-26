package utils

import (
	"fmt"
	"log"
	"os"
)

const TOKEN = "[nginx-reloader] "

var stdoutLogger = log.New(os.Stdout, "", log.LstdFlags)

func Fatalf(format string, v ...interface{}) {
	log.Fatalf(TOKEN+"[fatal] "+format, v...)
}

func Stdoutf(format string, v ...interface{}) {
	stdoutLogger.Printf(TOKEN+format, v...)
}

// func Stderrf(format string, v ...interface{}) {
//     log.Printf(TOKEN+"[error] "+format, v...)
// }

func Panicf(format string, v ...interface{}) {
	panic(fmt.Sprintf(format, v...))
}
