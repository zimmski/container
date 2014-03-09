package list

import (
	"runtime/debug"
	"testing"
)

// ListBenchmark is the base for all benchmarks of lists
type ListBenchmark struct {
	New func(b *testing.B) List
}

func (lb *ListBenchmark) BenchmarkPushSequentiel(b *testing.B) {
	l := lb.New(b)

	debug.FreeOSMemory()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Push(i)
	}
}

func (lb *ListBenchmark) BenchmarkUnshiftSequentiel(b *testing.B) {
	l := lb.New(b)

	debug.FreeOSMemory()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Unshift(i)
	}
}
