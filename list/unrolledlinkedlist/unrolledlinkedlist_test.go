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
