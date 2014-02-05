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
