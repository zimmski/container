package unrolledlinkedlist

import (
	"testing"

	"github.com/zimmski/container/list"
)

func TestRunAllTests(t *testing.T) {
	lt := &list.ListTest{
		New: func(t *testing.T) list.List {
			return New(7)
		},
	}

	lt.Run(t)
}
