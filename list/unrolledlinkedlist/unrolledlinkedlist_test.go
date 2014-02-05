package unrolledlinkedlist

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/container/list"
	"github.com/zimmski/container/util"
)

func TestRunAllTests(t *testing.T) {
	lt := &list.ListTest{
		New: func(t *testing.T) list.List {
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
	l.RemoveAt(4)
	l.RemoveAt(4)

	// remove a node from the front
	l.RemoveAt(0)
	l.RemoveAt(0)

	// remove a node from the back
	l.RemoveAt(l.Len() - 1)
	l.RemoveAt(l.Len() - 1)
}
