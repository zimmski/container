package unrolledlinkedlist

import (
	"testing"

	"github.com/zimmski/container/list"
)

type UnrolledLinkedListTest struct {
	list.ListTest
}

func TestRunAllTests(t *testing.T) {
	lt := &UnrolledLinkedListTest{
		list.ListTest{
			New: func(t *testing.T) list.List {
				return New(7)
			},
		},
	}

	lt.Run(t)
}
