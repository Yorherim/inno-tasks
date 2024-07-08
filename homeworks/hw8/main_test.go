package main

import "testing"

func BenchmarkOptimizeLogger(b *testing.B) {
	l := NewOptimizeLogger()
	for i := 0; i < b.N; i++ {
		l.Info("hello NewOptimizeLogger")
	}

}

func BenchmarkNoOptimizeLogger(b *testing.B) {
	l := NewNoOptimizeLogger()
	for i := 0; i < b.N; i++ {
		l.Info("hello NewNoOptimizeLogger")
	}
}
