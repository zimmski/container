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
	Nil(t, l.Shift())

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

	for i, va := range v {
		l.Unshift(va)

		Equal(t, l.Len(), i+1)
		Equal(t, l.First().Value, va)
		Equal(t, l.Last().Value, v[0])
	}

	Equal(t, l.Len(), vLen)

	i = vLen - 1
	n = l.Shift()

	for i > -1 && n != nil {
		Equal(t, v[i], n.Value)
		Nil(t, n.Next())
		Equal(t, l.Len(), i)

		i--
		n = l.Shift()
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

func TestFind(t *testing.T) {
	l := New()

	for _, vi := range v {
		f, ok := l.IndexOf(vi)

		Equal(t, f, -1)
		Equal(t, ok, false)

		ok = l.Contains(vi)

		Equal(t, ok, false)
	}

	fillList(t, l)

	for i, vi := range v {
		f, ok := l.IndexOf(vi)

		Equal(t, f, i)
		Equal(t, ok, true)

		ok = l.Contains(vi)

		Equal(t, ok, true)
	}
}

func TestGetSet(t *testing.T) {
	l := New()

	for i, _ := range v {
		n, err := l.Get(i)

		Nil(t, n)
		NotNil(t, err)

		err = l.Set(i, i+10)

		NotNil(t, err)

		n, err = l.Get(i)

		Nil(t, n)
		NotNil(t, err)
	}

	fillList(t, l)

	for i, vi := range v {
		n, err := l.Get(i)

		Equal(t, n.Value, vi)
		Nil(t, err)

		err = l.Set(i, i+10)

		Nil(t, err)

		n, err = l.Get(i)

		Equal(t, n.Value, i+10)
		Nil(t, err)
	}
}

func TestRemove(t *testing.T) {
	l1 := newFilledList(t)
	l2 := newFilledList(t)

	// do not allow removing elements from another list
	n1, _ := l1.Get(1)
	n2, _ := l2.Get(2)

	Nil(t, l1.Remove(n2))
	Equal(t, l1.Len(), vLen)
	Nil(t, l2.Remove(n1))
	Equal(t, l2.Len(), vLen)

	Equal(t, n1.Value, l1.Remove(n1).Value)
	Equal(t, l1.Len(), vLen-1)
	Equal(t, n2.Value, l2.Remove(n2).Value)
	Equal(t, l2.Len(), vLen-1)
}
