package unrolledlinkedlist

import (
	"testing"

	. "github.com/zimmski/container/test/assert"

	List "github.com/zimmski/container/list"
	"github.com/zimmski/container/util"
)

func TestRunAllTests(t *testing.T) {
	lt := &List.ListTest{
		New: func(t *testing.T) List.List {
			return New(7)
		},
	}

	lt.Run(t)
}

func TestNewWrongParameters(t *testing.T) {
	True(t, util.Panics(New, -1))
}

func TestGetNode(t *testing.T) {
	l := New(2)

	for i := 0; i < 5; i++ {
		l.Push(i)

		_, ic := l.getNode(i)
		Equal(t, ic, i%2)
	}

	n, ic := l.getNode(100)

	Nil(t, n)
	Equal(t, ic, -1)
}

func TestSmallMaxElementList(t *testing.T) {
	l := New(2)

	for i := 0; i < 10; i++ {
		l.Push(i)
	}

	// remove a node from the middle
	l.Remove(4)
	l.Remove(4)

	// remove a node from the front
	l.Remove(0)
	l.Remove(0)

	// remove a node from the back
	l.Remove(l.Len() - 1)
	l.Remove(l.Len() - 1)
}

func BenchmarkPushSequentiel(b *testing.B) {
	lb := &List.ListBenchmark{
		New: func(b *testing.B) List.List {
			return New(7)
		},
	}

	lb.BenchmarkPushSequentiel(b)
}

func BenchmarkUnshiftSequentiel(b *testing.B) {
	lb := &List.ListBenchmark{
		New: func(b *testing.B) List.List {
			return New(7)
		},
	}

	lb.BenchmarkUnshiftSequentiel(b)
}
