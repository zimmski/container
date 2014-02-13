package binarysearchtree

import (
	"testing"

	Tree "github.com/zimmski/container/tree"
)

func TestRunAllTests(t *testing.T) {
	tt := &Tree.TreeTest{
		New: func(t *testing.T) Tree.Tree {
			return New(func(a, b interface{}) int {
				switch {
				case a.(int) == b.(int):
					return 0
				case a.(int) < b.(int):
					return -1
				default:
					return 1
				}
			})
		},
	}

	tt.Run(t)
}
