package doublylinkedlist

import (
	"testing"

	List "github.com/zimmski/container/list"
)

func TestRunAllTests(t *testing.T) {
	lt := &List.ListTest{
		New: func(t *testing.T) List.List {
			return New()
		},
	}

	lt.Run(t)
}

func BenchmarkPushSequentiel(b *testing.B) {
	lb := &List.ListBenchmark{
		New: func(b *testing.B) List.List {
			return New()
		},
	}

	lb.BenchmarkPushSequentiel(b)
}

func BenchmarkUnshiftSequentiel(b *testing.B) {
	lb := &List.ListBenchmark{
		New: func(b *testing.B) List.List {
			return New()
		},
	}

	lb.BenchmarkUnshiftSequentiel(b)
}
