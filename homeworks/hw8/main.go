package main

import (
	"fmt"
	"os"
	"time"
)

const (
	dateFormat = "2006-01-02 15:04:05"
)

type OptimizeLogger struct {
}

func NewOptimizeLogger() OptimizeLogger {
	return OptimizeLogger{}
}

// Info на данный метод работает только со string
func (l OptimizeLogger) Info(v string) {
	date := time.Now().Format(dateFormat) + " "

	if _, err := os.Stdout.Write([]byte(date + v + "\n")); err != nil {
		return
	}
}

type NoOptimizeLogger struct {
}

func NewNoOptimizeLogger() *NoOptimizeLogger {
	return &NoOptimizeLogger{}
}

func (l *NoOptimizeLogger) Info(args ...any) {
	date := time.Now().Format(dateFormat) + " "
	fmt.Printf("%s%s\n", date, fmt.Sprint(args...))
}

// go build -gcflags "-m=1" main.go

func main() {
	optimizeLogger := NewOptimizeLogger()
	noOptimizeLogger := NewNoOptimizeLogger()

	optimizeLogger.Info("hello optimizeLogger")
	noOptimizeLogger.Info("hello noOptimizeLogger")
}
