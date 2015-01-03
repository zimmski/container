package list

import (
	"testing"

	. "github.com/zimmski/container/test/assert"
	"github.com/zimmski/go-leak"
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
	Equal(t, 0, leak.GoRoutineLeaks(func() {
		lt.NewFilledList(t)

		lt.TestBasic(t)
		lt.TestIterator(t)
		lt.TestChannels(t)
		lt.TestSlice(t)
		lt.TestInserts(t)
		lt.TestRemove(t)
		lt.TestRemoveOccurrence(t)
		lt.TestClear(t)
		lt.TestCopy(t)
		lt.TestIndexOf(t)
		lt.TestGetSet(t)
		lt.TestAddLists(t)
		lt.TestFuncs(t)
		lt.TestSwap(t)
		lt.TestMoves(t)

		lt.TestLeaks(t)
	}))
}

// FillList fills up a given list with V
func (lt *ListTest) FillList(t *testing.T, l List) {
	for i, va := range V {
		l.Push(va)

		Equal(t, l.Len(), i+1)
		n, ok := l.First()
		True(t, ok)
		Equal(t, n, V[0])
		n, ok = l.Last()
		True(t, ok)
		Equal(t, n, va)
	}

	Equal(t, l.Len(), VLen)
}

// NewFilledList creates a new list and calls FillList on it
func (lt *ListTest) NewFilledList(t *testing.T) List {
	l := lt.New(t)

	lt.FillList(t, l)

	return l
}

// NewDigitList creates a new list and fills it with numbers from 0 to 4
func (lt *ListTest) NewDigitList(t *testing.T) List {
	l := lt.New(t)

	for i := 0; i < 5; i++ {
		l.Push(i)
	}

	return l
}

// TestBasic tests basic list functionality
func (lt *ListTest) TestBasic(t *testing.T) {
	l := lt.New(t)

	Equal(t, l.Len(), 0)
	True(t, l.Empty())
	n, ok := l.First()
	False(t, ok)
	Nil(t, n)
	n, ok = l.Last()
	False(t, ok)
	Nil(t, n)
	n, ok = l.Pop()
	Nil(t, n)
	False(t, ok)
	n, ok = l.Shift()
	Nil(t, n)
	False(t, ok)

	lt.FillList(t, l)

	i := 0
	iter := l.Iter()
	NotNil(t, iter)

	for i < VLen {
		Equal(t, V[i], iter.Get())

		i++

		iter = iter.Next()

		if i < VLen {
			NotNil(t, iter)
		} else {
			Nil(t, iter)
		}
	}

	i = VLen - 1
	n, ok = l.Pop()

	for i > -1 && n != nil {
		Equal(t, V[i], n)
		True(t, ok)
		Equal(t, l.Len(), i)
		if i == 0 {
			True(t, l.Empty())
		} else {
			False(t, l.Empty())
		}

		i--
		n, ok = l.Pop()
	}

	Equal(t, i, -1)
	Nil(t, n)
	False(t, ok)
	Equal(t, l.Len(), 0)
	True(t, l.Empty())

	for i, va := range V {
		l.Unshift(va)

		Equal(t, l.Len(), i+1)
		n, ok := l.First()
		True(t, ok)
		Equal(t, n, va)
		n, ok = l.Last()
		True(t, ok)
		Equal(t, n, V[0])
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

// TestIterator tests list iterators
func (lt *ListTest) TestIterator(t *testing.T) {
	// empty iterators
	l := lt.New(t)

	Nil(t, l.Iter())
	Nil(t, l.IterBack())

	// one element
	l.Push(V[0])

	iter := l.Iter()
	NotNil(t, iter)
	Equal(t, V[0], iter.Get())
	Nil(t, iter.Next())

	iter = l.IterBack()
	NotNil(t, iter)
	Equal(t, V[0], iter.Get())
	Nil(t, iter.Previous())

	// full iterators
	l = lt.NewFilledList(t)

	i := 0

	for iter = l.Iter(); iter != nil; iter = iter.Next() {
		Equal(t, iter.Get(), V[i])

		iter.Set(i)

		Equal(t, iter.Get(), i)

		v, _ := l.Get(i)
		Equal(t, v, i)

		i++
	}

	Equal(t, i, VLen)

	l = lt.NewFilledList(t)

	i = VLen - 1

	for iter = l.IterBack(); iter != nil; iter = iter.Previous() {
		Equal(t, iter.Get(), V[i])

		iter.Set(i)

		Equal(t, iter.Get(), i)

		v, _ := l.Get(i)
		Equal(t, v, i)

		i--
	}

	Equal(t, i, -1)

	// iterate in wrong direction
	iter = l.Iter()
	Nil(t, iter.Previous())

	iter = l.IterBack()
	Nil(t, iter.Next())
}

// TestChannels tests list channels
func (lt *ListTest) TestChannels(t *testing.T) {
	// empty channels
	l := lt.New(t)

	i := 0

	for v := range l.Chan(0) {
		Equal(t, v, V[i])

		i++
	}

	Equal(t, i, 0)

	i = 0

	for v := range l.ChanBack(0) {
		Equal(t, v, V[i])

		i++
	}

	Equal(t, i, 0)

	// one element
	l.Push(1)

	i = 0

	for v := range l.Chan(0) {
		Equal(t, v, V[i])

		i++
	}

	Equal(t, i, 1)

	i = 0

	for v := range l.ChanBack(0) {
		Equal(t, v, V[i])

		i++
	}

	Equal(t, i, 1)

	// full iterators
	l = lt.NewFilledList(t)

	i = 0

	for v := range l.Chan(0) {
		Equal(t, v, V[i])

		i++
	}

	Equal(t, i, VLen)

	i = VLen - 1

	for v := range l.ChanBack(0) {
		Equal(t, v, V[i])

		i--
	}

	Equal(t, i, -1)
}

// TestSlice tests converting the list to slice
func (lt *ListTest) TestSlice(t *testing.T) {
	l := lt.New(t)
	Equal(t, l.Slice(), []interface{}{})

	lt.FillList(t, l)
	Equal(t, l.Slice(), V)

	l.Shift()
	Equal(t, l.Slice(), V[1:])

	l.Pop()
	Equal(t, l.Slice(), V[1:len(V)-1])
}

// TestInserts tests some insert methods
func (lt *ListTest) TestInserts(t *testing.T) {
	// Insert
	l1 := lt.NewFilledList(t)

	err := l1.Insert(0, 0)
	Nil(t, err)
	Equal(t, l1.Slice(), []interface{}{0, 1, "a", 2, "b", 3, "c", 4, "d"})
	Equal(t, l1.Len(), VLen+1)

	err = l1.Insert(l1.Len(), 0)
	Nil(t, err)
	Equal(t, l1.Slice(), []interface{}{0, 1, "a", 2, "b", 3, "c", 4, "d", 0})
	Equal(t, l1.Len(), VLen+2)

	err = l1.Insert(2, 0)
	Nil(t, err)
	Equal(t, l1.Slice(), []interface{}{0, 1, 0, "a", 2, "b", 3, "c", 4, "d", 0})
	Equal(t, l1.Len(), VLen+3)

	// out of bound
	err = l1.Insert(-1, 0)
	NotNil(t, err)
	err = l1.Insert(l1.Len()+1, 0)
	NotNil(t, err)
}

// TestRemove tests some remove methods
func (lt *ListTest) TestRemove(t *testing.T) {
	l := lt.NewFilledList(t)

	// out of bound
	_, err := l.Remove(-1)
	NotNil(t, err)
	_, err = l.Remove(l.Len())
	NotNil(t, err)

	// Remove Middle
	n, err := l.Remove(1)
	Nil(t, err)
	Equal(t, n, V[1])
	n, _ = l.Get(1)
	Equal(t, n, V[2])
	Equal(t, l.Len(), len(V)-1)

	// Remove First
	n, err = l.Remove(0)
	Nil(t, err)
	Equal(t, n, V[0])
	n, _ = l.First()
	Equal(t, n, V[2])
	Equal(t, l.Len(), len(V)-2)

	// Remove Last
	n, err = l.Remove(l.Len() - 1)
	Nil(t, err)
	Equal(t, n, V[len(V)-1])
	n, _ = l.Last()
	Equal(t, n, V[len(V)-2])
	Equal(t, l.Len(), len(V)-3)

	// Remove very last node
	l.Clear()
	l.Push(23)

	n, err = l.Remove(0)
	Nil(t, err)
	Equal(t, n, 23)
	n, ok := l.First()
	False(t, ok)
	Nil(t, n)
	n, ok = l.Last()
	False(t, ok)
	Nil(t, n)

	// remove structure
	l = lt.New(t)

	for i := 0; i < 10; i++ {
		l.Push(i % 10)
	}
	Equal(t, l.Slice(), []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})

	l.Remove(2)
	Equal(t, l.Slice(), []interface{}{0, 1, 3, 4, 5, 6, 7, 8, 9})
	l.Remove(2)
	Equal(t, l.Slice(), []interface{}{0, 1, 4, 5, 6, 7, 8, 9})
	l.Remove(2)
	Equal(t, l.Slice(), []interface{}{0, 1, 5, 6, 7, 8, 9})
	l.Remove(2)
	Equal(t, l.Slice(), []interface{}{0, 1, 6, 7, 8, 9})
}

// TestRemoveOccurrence tests the remove occurrence methods
func (lt *ListTest) TestRemoveOccurrence(t *testing.T) {
	l := lt.New(t)

	for i := 0; i < 5; i++ {
		l.Push(i % 2)
	}

	ok := l.RemoveFirstOccurrence(0)
	True(t, ok)
	Equal(t, l.Len(), 4)
	Equal(t, l.Slice(), []interface{}{1, 0, 1, 0})

	ok = l.RemoveFirstOccurrence(0)
	True(t, ok)
	Equal(t, l.Len(), 3)
	Equal(t, l.Slice(), []interface{}{1, 1, 0})

	ok = l.RemoveFirstOccurrence(0)
	True(t, ok)
	Equal(t, l.Len(), 2)
	Equal(t, l.Slice(), []interface{}{1, 1})

	ok = l.RemoveFirstOccurrence(0)
	False(t, ok)
	Equal(t, l.Len(), 2)
	Equal(t, l.Slice(), []interface{}{1, 1})

	l.Clear()

	for i := 0; i < 5; i++ {
		l.Push(i % 2)
	}

	ok = l.RemoveLastOccurrence(0)
	True(t, ok)
	Equal(t, l.Len(), 4)
	Equal(t, l.Slice(), []interface{}{0, 1, 0, 1})

	ok = l.RemoveLastOccurrence(0)
	True(t, ok)
	Equal(t, l.Len(), 3)
	Equal(t, l.Slice(), []interface{}{0, 1, 1})

	ok = l.RemoveLastOccurrence(0)
	True(t, ok)
	Equal(t, l.Len(), 2)
	Equal(t, l.Slice(), []interface{}{1, 1})

	ok = l.RemoveLastOccurrence(0)
	False(t, ok)
	Equal(t, l.Len(), 2)
	Equal(t, l.Slice(), []interface{}{1, 1})
}

// TestClear tests clearing the list
func (lt *ListTest) TestClear(t *testing.T) {
	l := lt.NewFilledList(t)

	l.Clear()

	Equal(t, l.Len(), 0)
	n, ok := l.First()
	False(t, ok)
	Nil(t, n)
	n, ok = l.Last()
	False(t, ok)
	Nil(t, n)
	n, ok = l.Pop()
	Nil(t, n)
	False(t, ok)
}

// TestCopy tests copying a list
func (lt *ListTest) TestCopy(t *testing.T) {
	l1 := lt.NewFilledList(t)

	l2 := l1.Copy()

	Equal(t, l1.Len(), l2.Len())

	n1 := l1.Iter()
	NotNil(t, n1)
	n2 := l2.Iter()
	NotNil(t, n2)

	if n1 != nil && n2 != nil {
		for {
			Equal(t, n1.Get(), n2.Get())

			n1 = n1.Next()
			n2 = n2.Next()

			if (n1 == nil && n2 != nil) || (n1 != nil && n2 == nil) {
				Fail(t, "n1 not equal to n2")
			}

			if n1 == nil {
				break
			}
		}
	}
}

// TestIndexOf tests the index of methods
func (lt *ListTest) TestIndexOf(t *testing.T) {
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

// TestGetSet tests getters and setters
func (lt *ListTest) TestGetSet(t *testing.T) {
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

	for i := range V {
		n, err := l.Get(i)

		Equal(t, n, V[i])
		Nil(t, err)

		err = l.Set(i, i+10)

		Nil(t, err)

		n, err = l.Get(i)

		Equal(t, n, i+10)
		Nil(t, err)
	}
}

// TestAddLists tests the insert list methods
func (lt *ListTest) TestAddLists(t *testing.T) {
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
	Equal(t, l1.Slice(), []interface{}{3, 4, 5, 6})

	l1.UnshiftList(l3)
	Equal(t, l1.Slice(), []interface{}{1, 2, 3, 4, 5, 6})

	// empty lists
	l4 := lt.New(t)

	l1.PushList(l4)
	Equal(t, l1.Slice(), []interface{}{1, 2, 3, 4, 5, 6})

	l1.UnshiftList(l4)
	Equal(t, l1.Slice(), []interface{}{1, 2, 3, 4, 5, 6})
}

// TestFuncs tests all methods with functions as parameters
func (lt *ListTest) TestFuncs(t *testing.T) {
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
	Equal(t, l.Slice(), []interface{}{1, "a", 3, "b", 3, "c", 4, "d"})
	False(t, l.SetFunc(func(v interface{}) bool {
		return v == "z"
	}, 4))
	Equal(t, l.Slice(), []interface{}{1, "a", 3, "b", 3, "c", 4, "d"})
}

// TestSwap tests swap
func (lt *ListTest) TestSwap(t *testing.T) {
	l := lt.NewFilledList(t)

	l.Swap(0, 0)
	Equal(t, l.Slice(), V)

	l.Swap(0, 1)
	v, _ := l.Get(0)
	Equal(t, v, V[1])
	v, _ = l.Get(1)
	Equal(t, v, V[0])

	l.Swap(0, 1)
	Equal(t, l.Slice(), V)
}

// TestMoves tests all move methods
func (lt *ListTest) TestMoves(t *testing.T) {
	l := lt.NewDigitList(t)
	ll := l.Len()
	lll := ll - 1

	// out of bounds
	err := l.MoveAfter(-1, 0)
	NotNil(t, err)
	err = l.MoveAfter(0, ll)
	NotNil(t, err)
	err = l.MoveBefore(-1, 0)
	NotNil(t, err)
	err = l.MoveBefore(0, ll)
	NotNil(t, err)
	err = l.MoveToBack(-1)
	NotNil(t, err)
	err = l.MoveToBack(ll)
	NotNil(t, err)
	err = l.MoveToFront(-1)
	NotNil(t, err)
	err = l.MoveToFront(ll)
	NotNil(t, err)

	// basics
	l.MoveAfter(0, lll)
	Equal(t, l.Slice(), []interface{}{1, 2, 3, 4, 0})
	Equal(t, l.Len(), ll)

	l.MoveAfter(lll, 0)
	Equal(t, l.Slice(), []interface{}{1, 0, 2, 3, 4})
	Equal(t, l.Len(), ll)

	l.MoveAfter(1, 2)
	Equal(t, l.Slice(), []interface{}{1, 2, 0, 3, 4})
	Equal(t, l.Len(), ll)

	l.MoveAfter(2, 1)
	Equal(t, l.Slice(), []interface{}{1, 2, 0, 3, 4})
	Equal(t, l.Len(), ll)

	l.Push(0)
	Equal(t, l.Slice(), []interface{}{1, 2, 0, 3, 4, 0})
	Equal(t, l.Len(), ll+1)

	l = lt.NewDigitList(t)

	l.MoveBefore(0, lll)
	Equal(t, l.Slice(), []interface{}{1, 2, 3, 0, 4})
	Equal(t, l.Len(), ll)

	l.MoveBefore(lll, 0)
	Equal(t, l.Slice(), []interface{}{4, 1, 2, 3, 0})
	Equal(t, l.Len(), ll)

	l.MoveBefore(1, 2)
	Equal(t, l.Slice(), []interface{}{4, 1, 2, 3, 0})
	Equal(t, l.Len(), ll)

	l.MoveBefore(2, 1)
	Equal(t, l.Slice(), []interface{}{4, 2, 1, 3, 0})
	Equal(t, l.Len(), ll)

	l.Push(0)
	Equal(t, l.Slice(), []interface{}{4, 2, 1, 3, 0, 0})
	Equal(t, l.Len(), ll+1)

	l = lt.NewDigitList(t)

	l.MoveToBack(0)
	Equal(t, l.Slice(), []interface{}{1, 2, 3, 4, 0})
	Equal(t, l.Len(), ll)

	l.MoveToBack(lll)
	Equal(t, l.Slice(), []interface{}{1, 2, 3, 4, 0})
	Equal(t, l.Len(), ll)

	l.MoveToBack(2)
	Equal(t, l.Slice(), []interface{}{1, 2, 4, 0, 3})
	Equal(t, l.Len(), ll)

	l.Push(0)
	Equal(t, l.Slice(), []interface{}{1, 2, 4, 0, 3, 0})
	Equal(t, l.Len(), ll+1)

	l = lt.NewDigitList(t)

	l.MoveToFront(0)
	Equal(t, l.Slice(), []interface{}{0, 1, 2, 3, 4})
	Equal(t, l.Len(), ll)

	l.MoveToFront(lll)
	Equal(t, l.Slice(), []interface{}{4, 0, 1, 2, 3})
	Equal(t, l.Len(), ll)

	l.MoveToFront(2)
	Equal(t, l.Slice(), []interface{}{1, 4, 0, 2, 3})
	Equal(t, l.Len(), ll)

	l.Push(0)
	Equal(t, l.Slice(), []interface{}{1, 4, 0, 2, 3, 0})
	Equal(t, l.Len(), ll+1)
}

// TestLeaks test for leaks
func (lt *ListTest) TestLeaks(t *testing.T) {
	l := lt.New(t)

	Equal(t, 0, leak.MemoryLeaks(func() {
		l.Push(1)
		l.Push(2)
		l.Push(3)

		l.Pop()
		l.Pop()
		l.Pop()
	}))

	Equal(t, 0, leak.MemoryLeaks(func() {
		l.Push(1)
		l.Push(2)
		l.Push(3)

		l.Clear()
	}))
}
