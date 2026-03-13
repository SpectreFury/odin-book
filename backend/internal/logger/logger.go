package logger

import (
	"log"
	"os"
)

type logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

var l *logger

func init() {
	l = &logger {
		infoLogger : log.New(os.Stdout, "[INFO]: ", log.LstdFlags | log.Lmsgprefix),
		errorLogger : log.New(os.Stdout, "[ERROR]: ", log.LstdFlags | log.Lmsgprefix),
	}
}


func Info(v ...any) {
	l.infoLogger.Println(v...)
}

func Error(v ...any) {
	l.errorLogger.Println(v...)
}

func Fatal(v ...any) {
	l.errorLogger.Fatal(v...)
}
