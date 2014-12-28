package linkedlist

import (
	"testing"

	. "github.com/zimmski/container/test/assert"

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

func TestFindParentNode(t *testing.T) {
	l := New()

	for i := 0; i < 5; i++ {
		l.Push(i)
	}

	Nil(t, l.findParentNode(l.first))
	Equal(t, l.findParentNode(l.first.next), l.first)
	Equal(t, l.findParentNode(l.last).value, 3)

	// find parent to nil
	Nil(t, l.findParentNode(nil))

	// not existing node
	n := &node{}

	Nil(t, l.findParentNode(n))
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
