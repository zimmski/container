package unrolledlinkedlist

import (
	"fmt"
	"testing"

	. "github.com/stretchr/testify/assert"
)

var v = []interface{}{1, "a", 2, "b", 3, "c", 4, "d"}
var vLen = len(v)

func fillList(t *testing.T, l *UnrolledLinkedList) {
	for i, va := range v {
		l.Push(va)

		Equal(t, l.Len(), i+1)
		Equal(t, l.First().Value(), v[0])
		Equal(t, l.Last().Value(), va)
	}

	Equal(t, l.Len(), vLen)
}

func newFilledList(t *testing.T) *UnrolledLinkedList {
	l := New(4)

	fillList(t, l)

	return l
}

func printList(l *UnrolledLinkedList) {
	fmt.Printf("Len: %d\n", l.len)
	for n := l.first; n != nil; n = n.Next() {
		fmt.Printf("\t%+v\n", n.values)
	}
}

func TestBasic(t *testing.T) {
	l := New(4)

	Equal(t, l.Len(), 0)
	Nil(t, l.First())
	Nil(t, l.Last())
	n, ok := l.Pop()
	Nil(t, n)
	False(t, ok)
	n, ok = l.Shift()
	Nil(t, n)
	False(t, ok)

	fillList(t, l)

	i := 0
	c := l.First()

	for i < vLen {
		Equal(t, v[i], c.Value())

		i++
		if i < vLen {
			True(t, c.Next())
		}
	}

	False(t, c.Next())
	Nil(t, c.Value())

	i = vLen - 1
	n, ok = l.Pop()

	for i > -1 && n != nil {
		Equal(t, v[i], n)
		True(t, ok)
		Equal(t, l.Len(), i)

		i--
		n, ok = l.Pop()
	}

	Equal(t, i, -1)
	Nil(t, n)
	False(t, ok)
	Equal(t, l.Len(), 0)

	for i, va := range v {
		l.Unshift(va)

		Equal(t, l.Len(), i+1)
		Equal(t, l.First().Value(), va)
		Equal(t, l.Last().Value(), v[0])
	}

	Equal(t, l.Len(), vLen)

	i = vLen - 1
	n, ok = l.Shift()

	for i > -1 && n != nil {
		Equal(t, v[i], n)
		True(t, ok)
		Equal(t, l.Len(), i)

		i--
		n, ok = l.Shift()
	}

	Equal(t, i, -1)
	Nil(t, n)
	Equal(t, l.Len(), 0)
}

func TestToArray(t *testing.T) {
	l := New(4)
	Equal(t, l.ToArray(), []interface{}{})

	fillList(t, l)
	Equal(t, l.ToArray(), v)

	l.Shift()
	Equal(t, l.ToArray(), v[1:])

	l.Pop()
	Equal(t, l.ToArray(), v[1:len(v)-1])
}

func TestInserts(t *testing.T) {
	// InsertAt
	l1 := newFilledList(t)

	err := l1.InsertAt(0, 0)
	Nil(t, err)
	Equal(t, l1.ToArray(), []interface{}{0, 1, "a", 2, "b", 3, "c", 4, "d"})
	Equal(t, l1.Len(), vLen+1)

	err = l1.InsertAt(l1.Len(), 0)
	Nil(t, err)
	Equal(t, l1.ToArray(), []interface{}{0, 1, "a", 2, "b", 3, "c", 4, "d", 0})
	Equal(t, l1.Len(), vLen+2)

	err = l1.InsertAt(2, 0)
	Nil(t, err)
	Equal(t, l1.ToArray(), []interface{}{0, 1, 0, "a", 2, "b", 3, "c", 4, "d", 0})
	Equal(t, l1.Len(), vLen+3)

	// out of bound
	err = l1.InsertAt(-1, 0)
	NotNil(t, err)
	err = l1.InsertAt(l1.Len()+1, 0)
	NotNil(t, err)
}

func TestRemove(t *testing.T) {
	l := newFilledList(t)

	// out of bound
	_, err := l.RemoveAt(-1)
	NotNil(t, err)
	_, err = l.RemoveAt(l.Len())
	NotNil(t, err)

	// Remove Middle
	n1, err := l.RemoveAt(1)
	Nil(t, err)
	Equal(t, n1, v[1])
	n2, _ := l.Get(1)
	Equal(t, n2, v[2])
	Equal(t, l.Len(), len(v)-1)

	// Remove First
	n1, err = l.RemoveAt(0)
	Nil(t, err)
	Equal(t, n1, v[0])
	n3 := l.First()
	Equal(t, n3.Value(), v[2])
	Equal(t, l.Len(), len(v)-2)

	// Remove Last
	n1, err = l.RemoveAt(l.Len() - 1)
	Nil(t, err)
	Equal(t, n1, v[len(v)-1])
	n3 = l.Last()
	Equal(t, n3.Value(), v[len(v)-2])
	Equal(t, l.Len(), len(v)-3)

	// Remove very last node
	l.Clear()
	l.Push(23)

	n1, err = l.RemoveAt(0)
	Nil(t, err)
	Equal(t, n1, 23)
	Nil(t, l.First())
	Nil(t, l.Last())

	// remove structure
	l = New(7)

	for i := 0; i < 10; i++ {
		l.Push(i % 10)
	}
	Equal(t, l.ToArray(), []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})

	l.RemoveAt(2)
	Equal(t, l.ToArray(), []interface{}{0, 1, 3, 4, 5, 6, 7, 8, 9})
	l.RemoveAt(2)
	Equal(t, l.ToArray(), []interface{}{0, 1, 4, 5, 6, 7, 8, 9})
	l.RemoveAt(2)
	Equal(t, l.ToArray(), []interface{}{0, 1, 5, 6, 7, 8, 9})
	l.RemoveAt(2)
	Equal(t, l.ToArray(), []interface{}{0, 1, 6, 7, 8, 9})
}

func TestRemoveOccurrence(t *testing.T) {
	l := New(4)

	for i := 0; i < 5; i++ {
		l.Push(i % 2)
	}

	ok := l.RemoveFirstOccurrence(0)
	True(t, ok)
	Equal(t, l.Len(), 4)
	Equal(t, l.ToArray(), []interface{}{1, 0, 1, 0})

	ok = l.RemoveFirstOccurrence(0)
	True(t, ok)
	Equal(t, l.Len(), 3)
	Equal(t, l.ToArray(), []interface{}{1, 1, 0})

	ok = l.RemoveFirstOccurrence(0)
	True(t, ok)
	Equal(t, l.Len(), 2)
	Equal(t, l.ToArray(), []interface{}{1, 1})

	ok = l.RemoveFirstOccurrence(0)
	False(t, ok)
	Equal(t, l.Len(), 2)
	Equal(t, l.ToArray(), []interface{}{1, 1})

	l.Clear()

	for i := 0; i < 5; i++ {
		l.Push(i % 2)
	}

	ok = l.RemoveLastOccurrence(0)
	True(t, ok)
	Equal(t, l.Len(), 4)
	Equal(t, l.ToArray(), []interface{}{0, 1, 0, 1})

	ok = l.RemoveLastOccurrence(0)
	True(t, ok)
	Equal(t, l.Len(), 3)
	Equal(t, l.ToArray(), []interface{}{0, 1, 1})

	ok = l.RemoveLastOccurrence(0)
	True(t, ok)
	Equal(t, l.Len(), 2)
	Equal(t, l.ToArray(), []interface{}{1, 1})

	ok = l.RemoveLastOccurrence(0)
	False(t, ok)
	Equal(t, l.Len(), 2)
	Equal(t, l.ToArray(), []interface{}{1, 1})
}

func TestClear(t *testing.T) {
	l := newFilledList(t)

	l.Clear()

	Equal(t, l.Len(), 0)
	Nil(t, l.First())
	Nil(t, l.Last())
	n, ok := l.Pop()
	Nil(t, n)
	False(t, ok)
}

func TestCopy(t *testing.T) {
	l1 := newFilledList(t)

	l2 := l1.Copy()

	Equal(t, l1.Len(), l2.Len())

	n1 := l1.First()
	n2 := l2.First()

	if n1 != nil && n2 != nil {
		for {
			Equal(t, n1.Value(), n2.Value())

			ok1 := n1.Next()
			ok2 := n2.Next()

			Equal(t, ok1, ok2)

			if !ok1 {
				break
			}
		}
	}

	Nil(t, n1.Value())
	Nil(t, n2.Value())
}

func TestFind(t *testing.T) {
	l := New(4)

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

	l.Clear()

	f, ok := l.IndexOf(0)
	Equal(t, f, -1)
	Equal(t, ok, false)

	f, ok = l.LastIndexOf(0)
	Equal(t, f, -1)
	Equal(t, ok, false)

	for i := 0; i < 4; i++ {
		l.Push(0)

		f, ok = l.IndexOf(0)
		Equal(t, f, 0)
		Equal(t, ok, true)

		f, ok = l.LastIndexOf(0)
		Equal(t, f, i)
		Equal(t, ok, true)
	}

	// not found in nonempty list
	f, ok = l.IndexOf(100)
	Equal(t, f, -1)
	Equal(t, ok, false)

	f, ok = l.LastIndexOf(100)
	Equal(t, f, -1)
	Equal(t, ok, false)
}

func TestGetSet(t *testing.T) {
	l := New(4)

	for i := range v {
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

		Equal(t, n, vi)
		Nil(t, err)

		err = l.Set(i, i+10)

		Nil(t, err)

		n, err = l.Get(i)

		Equal(t, n, i+10)
		Nil(t, err)
	}
}

func TestAddLists(t *testing.T) {
	l1 := New(4)
	l1.Push(3)
	l1.Push(4)

	l2 := New(4)
	l2.Push(5)
	l2.Push(6)

	l3 := New(4)
	l3.Push(2)
	l3.Push(1)

	l1.PushList(l2)
	Equal(t, l1.ToArray(), []interface{}{3, 4, 5, 6})

	l1.UnshiftList(l3)
	Equal(t, l1.ToArray(), []interface{}{1, 2, 3, 4, 5, 6})

	// empty lists
	l4 := New(4)

	l1.PushList(l4)
	Equal(t, l1.ToArray(), []interface{}{1, 2, 3, 4, 5, 6})

	l1.UnshiftList(l4)
	Equal(t, l1.ToArray(), []interface{}{1, 2, 3, 4, 5, 6})
}

func TestFunc(t *testing.T) {
	l := newFilledList(t)

	Equal(t, v[1], l.GetFunc(func(v interface{}) bool {
		return v == "a"
	}))
	Nil(t, l.GetFunc(func(v interface{}) bool {
		return v == "z"
	}))

	l.SetFunc(func(v interface{}) bool {
		return v == 2
	}, 3)
	Equal(t, l.ToArray(), []interface{}{1, "a", 3, "b", 3, "c", 4, "d"})
	l.SetFunc(func(v interface{}) bool {
		return v == "z"
	}, 4)
	Equal(t, l.ToArray(), []interface{}{1, "a", 3, "b", 3, "c", 4, "d"})
}
