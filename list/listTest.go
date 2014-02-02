package list

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

// V holds the value for basic list tests
var V = []interface{}{1, "a", 2, "b", 3, "c", 4, "d"}

// VLen is the length of V
var VLen = len(V)

// ListTest is the base for all tests of lists
type ListTest struct {
	New func(t *testing.T) List
}

// Run executes all basic list tests
func (lt *ListTest) Run(t *testing.T) {
	lt.NewFilledList(t)

	lt.testBasic(t)
	lt.testToArray(t)
	lt.testInserts(t)
	lt.testRemove(t)
	lt.testRemoveOccurrence(t)
	lt.testClear(t)
	lt.testCopy(t)
	lt.testFind(t)
	lt.testGetSet(t)
	lt.testAddLists(t)
	lt.testFunc(t)
}

// FillList fills up a given list with V
func (lt *ListTest) FillList(t *testing.T, l List) {
	for i, va := range V {
		l.Push(va)

		Equal(t, l.Len(), i+1)
		Equal(t, l.First().Value(), V[0])
		Equal(t, l.Last().Value(), va)
	}

	Equal(t, l.Len(), VLen)
}

// NewFilledList createes a new list and calls FillList on it
func (lt *ListTest) NewFilledList(t *testing.T) List {
	l := lt.New(t)

	lt.FillList(t, l)

	return l
}

func (lt *ListTest) testBasic(t *testing.T) {
	l := lt.New(t)

	Equal(t, l.Len(), 0)
	Nil(t, l.First())
	Nil(t, l.Last())
	n, ok := l.Pop()
	Nil(t, n)
	False(t, ok)
	n, ok = l.Shift()
	Nil(t, n)
	False(t, ok)

	lt.FillList(t, l)

	i := 0
	c := l.First()

	for i < VLen {
		Equal(t, V[i], c.Value())

		i++
		if i < VLen {
			True(t, c.Next())
		}
	}

	False(t, c.Next())
	Nil(t, c.Value())

	i = VLen - 1
	n, ok = l.Pop()

	for i > -1 && n != nil {
		Equal(t, V[i], n)
		True(t, ok)
		Equal(t, l.Len(), i)

		i--
		n, ok = l.Pop()
	}

	Equal(t, i, -1)
	Nil(t, n)
	False(t, ok)
	Equal(t, l.Len(), 0)

	for i, va := range V {
		l.Unshift(va)

		Equal(t, l.Len(), i+1)
		Equal(t, l.First().Value(), va)
		Equal(t, l.Last().Value(), V[0])
	}

	Equal(t, l.Len(), VLen)

	i = VLen - 1
	n, ok = l.Shift()

	for i > -1 && n != nil {
		Equal(t, V[i], n)
		True(t, ok)
		Equal(t, l.Len(), i)

		i--
		n, ok = l.Shift()
	}

	Equal(t, i, -1)
	Nil(t, n)
	Equal(t, l.Len(), 0)
}

func (lt *ListTest) testToArray(t *testing.T) {
	l := lt.New(t)
	Equal(t, l.ToArray(), []interface{}{})

	lt.FillList(t, l)
	Equal(t, l.ToArray(), V)

	l.Shift()
	Equal(t, l.ToArray(), V[1:])

	l.Pop()
	Equal(t, l.ToArray(), V[1:len(V)-1])
}

func (lt *ListTest) testInserts(t *testing.T) {
	// InsertAt
	l1 := lt.NewFilledList(t)

	err := l1.InsertAt(0, 0)
	Nil(t, err)
	Equal(t, l1.ToArray(), []interface{}{0, 1, "a", 2, "b", 3, "c", 4, "d"})
	Equal(t, l1.Len(), VLen+1)

	err = l1.InsertAt(l1.Len(), 0)
	Nil(t, err)
	Equal(t, l1.ToArray(), []interface{}{0, 1, "a", 2, "b", 3, "c", 4, "d", 0})
	Equal(t, l1.Len(), VLen+2)

	err = l1.InsertAt(2, 0)
	Nil(t, err)
	Equal(t, l1.ToArray(), []interface{}{0, 1, 0, "a", 2, "b", 3, "c", 4, "d", 0})
	Equal(t, l1.Len(), VLen+3)

	// out of bound
	err = l1.InsertAt(-1, 0)
	NotNil(t, err)
	err = l1.InsertAt(l1.Len()+1, 0)
	NotNil(t, err)
}

func (lt *ListTest) testRemove(t *testing.T) {
	l := lt.NewFilledList(t)

	// out of bound
	_, err := l.RemoveAt(-1)
	NotNil(t, err)
	_, err = l.RemoveAt(l.Len())
	NotNil(t, err)

	// Remove Middle
	n1, err := l.RemoveAt(1)
	Nil(t, err)
	Equal(t, n1, V[1])
	n2, _ := l.Get(1)
	Equal(t, n2, V[2])
	Equal(t, l.Len(), len(V)-1)

	// Remove First
	n1, err = l.RemoveAt(0)
	Nil(t, err)
	Equal(t, n1, V[0])
	n3 := l.First()
	Equal(t, n3.Value(), V[2])
	Equal(t, l.Len(), len(V)-2)

	// Remove Last
	n1, err = l.RemoveAt(l.Len() - 1)
	Nil(t, err)
	Equal(t, n1, V[len(V)-1])
	n3 = l.Last()
	Equal(t, n3.Value(), V[len(V)-2])
	Equal(t, l.Len(), len(V)-3)

	// Remove very last node
	l.Clear()
	l.Push(23)

	n1, err = l.RemoveAt(0)
	Nil(t, err)
	Equal(t, n1, 23)
	Nil(t, l.First())
	Nil(t, l.Last())

	// remove structure
	l = lt.New(t)

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

func (lt *ListTest) testRemoveOccurrence(t *testing.T) {
	l := lt.New(t)

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

func (lt *ListTest) testClear(t *testing.T) {
	l := lt.NewFilledList(t)

	l.Clear()

	Equal(t, l.Len(), 0)
	Nil(t, l.First())
	Nil(t, l.Last())
	n, ok := l.Pop()
	Nil(t, n)
	False(t, ok)
}

func (lt *ListTest) testCopy(t *testing.T) {
	l1 := lt.NewFilledList(t)

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

func (lt *ListTest) testFind(t *testing.T) {
	l := lt.New(t)

	for _, vi := range V {
		f, ok := l.IndexOf(vi)
		Equal(t, f, -1)
		Equal(t, ok, false)

		ok = l.Contains(vi)
		Equal(t, ok, false)
	}

	lt.FillList(t, l)

	for i, vi := range V {
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

func (lt *ListTest) testGetSet(t *testing.T) {
	l := lt.New(t)

	for i := range V {
		n, err := l.Get(i)

		Nil(t, n)
		NotNil(t, err)

		err = l.Set(i, i+10)

		NotNil(t, err)

		n, err = l.Get(i)

		Nil(t, n)
		NotNil(t, err)
	}

	lt.FillList(t, l)

	for i, vi := range V {
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

func (lt *ListTest) testAddLists(t *testing.T) {
	l1 := lt.New(t)
	l1.Push(3)
	l1.Push(4)

	l2 := lt.New(t)
	l2.Push(5)
	l2.Push(6)

	l3 := lt.New(t)
	l3.Push(2)
	l3.Push(1)

	l1.PushList(l2)
	Equal(t, l1.ToArray(), []interface{}{3, 4, 5, 6})

	l1.UnshiftList(l3)
	Equal(t, l1.ToArray(), []interface{}{1, 2, 3, 4, 5, 6})

	// empty lists
	l4 := lt.New(t)

	l1.PushList(l4)
	Equal(t, l1.ToArray(), []interface{}{1, 2, 3, 4, 5, 6})

	l1.UnshiftList(l4)
	Equal(t, l1.ToArray(), []interface{}{1, 2, 3, 4, 5, 6})
}

func (lt *ListTest) testFunc(t *testing.T) {
	l := lt.NewFilledList(t)

	n, ok := l.GetFunc(func(v interface{}) bool {
		return v == "a"
	})
	Equal(t, V[1], n)
	True(t, ok)
	n, ok = l.GetFunc(func(v interface{}) bool {
		return v == "z"
	})
	Nil(t, nil)
	False(t, ok)

	True(t, l.SetFunc(func(v interface{}) bool {
		return v == 2
	}, 3))
	Equal(t, l.ToArray(), []interface{}{1, "a", 3, "b", 3, "c", 4, "d"})
	False(t, l.SetFunc(func(v interface{}) bool {
		return v == "z"
	}, 4))
	Equal(t, l.ToArray(), []interface{}{1, "a", 3, "b", 3, "c", 4, "d"})
}
