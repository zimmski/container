package list

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

var v = []interface{}{1, "a", 2, "b"}
var vLen = len(v)

func fillList(t *testing.T, l *LinkedList) {
	for i, va := range v {
		l.Push(va)

		Equal(t, l.Len(), i+1)
		Equal(t, l.First().Value, v[0])
		Equal(t, l.Last().Value, va)
	}

	Equal(t, l.Len(), vLen)
}

func newFilledList(t *testing.T) *LinkedList {
	l := New()

	fillList(t, l)

	return l
}

func TestBasic(t *testing.T) {
	l := New()

	Equal(t, l.Len(), 0)
	Nil(t, l.First())
	Nil(t, l.Last())
	Nil(t, l.Pop())

	fillList(t, l)

	i := 0
	n := l.First()

	for i < vLen && n != nil {
		Equal(t, v[i], n.Value)

		i++
		n = n.Next()
	}

	Nil(t, n)

	i = vLen - 1
	n = l.Pop()

	for i > -1 && n != nil {
		Equal(t, v[i], n.Value)
		Nil(t, n.Next())
		Equal(t, l.Len(), i)

		i--
		n = l.Pop()
	}

	Equal(t, i, -1)
	Nil(t, n)
	Equal(t, l.Len(), 0)
}

func TestClear(t *testing.T) {
	l := newFilledList(t)

	l.Clear()

	Equal(t, l.Len(), 0)
	Nil(t, l.First())
	Nil(t, l.Last())
	Nil(t, l.Pop())
}

func TestCopy(t *testing.T) {
	l1 := newFilledList(t)

	l2 := l1.Copy()

	Equal(t, l1.Len(), l2.Len())

	n1 := l1.First()
	n2 := l2.First()

	for n1 != nil && n2 != nil {
		Equal(t, n1.Value, n2.Value)

		n1 = n1.Next()
		n2 = n2.Next()
	}

	Nil(t, n1)
	Nil(t, n2)
}
