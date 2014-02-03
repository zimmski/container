package linkedlist

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/zimmski/container/list"
)

func TestRunAllTests(t *testing.T) {
	lt := &list.ListTest{
		New: func(t *testing.T) list.List {
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
