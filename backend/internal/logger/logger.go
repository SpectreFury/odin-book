package logger

import (
	"fmt"
	"time"
)

type Logger struct {}

func (l* Logger) Print(level string, text string) {
	now := time.Now()
	timestamp := now.Format("2006-01-02 15:04:05")
	fmt.Printf("[%s]: [%s] %s\n", level, timestamp, text)
}

func (l* Logger) Log(text string) {
	l.Print("INFO", text)
}

func (l* Logger) Debug(text string) {
	l.Print("DEBUG", text)
}

func (l* Logger) Error(text string) {
	l.Print("ERROR", text)
}

func (l* Logger) Warning(text string) {
	l.Print("WARNING", text)
}
