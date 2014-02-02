package linkedlist

import (
	"testing"

	"github.com/zimmski/container/list"
)

type LinkedListTest struct {
	list.ListTest
}

func TestRunAllTests(t *testing.T) {
	lt := &LinkedListTest{
		list.ListTest{
			New: func(t *testing.T) list.List {
				return New()
			},
		},
	}

	lt.Run(t)
}
