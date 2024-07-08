package main

import (
	"fmt"
	"os"
)

type OptimizeLogger struct {
}

func NewOptimizeLogger() OptimizeLogger {
	return OptimizeLogger{}
}

// Info на данный метод работает только со string
func (l OptimizeLogger) Info(v string) {
	var data []byte
	data = append([]byte(v), '\n')

	if _, err := os.Stdout.Write(data); err != nil {
		return
	}
}

type NoOptimizeLogger struct {
}

func NewNoOptimizeLogger() *NoOptimizeLogger {
	return &NoOptimizeLogger{}
}

func (l *NoOptimizeLogger) Info(args ...any) {
	fmt.Println(args...)
}

// go build -gcflags "-m=1" main.go

func main() {
	optimizeLogger := NewOptimizeLogger()
	noOptimizeLogger := NewNoOptimizeLogger()

	optimizeLogger.Info("hello optimizeLogger")
	noOptimizeLogger.Info("hello noOptimizeLogger")
}
