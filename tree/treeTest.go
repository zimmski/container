package tree

import (
	"sort"
	"testing"

	. "github.com/zimmski/container/test/assert"
)

// VRaw holds the unsorted value for basic tree tests
var VRaw = []interface{}{5, 3, 1, 4, 6, 2}

// V holds the sorted value for basic tree tests
var V = []interface{}{1, 2, 3, 4, 5, 6}

// VLen is the length of V
var VLen = len(V)

// TreeTest is the base for all tests of trees
type TreeTest struct {
	New func(t *testing.T) Tree
}

// Run executes all basic tree tests
func (tt *TreeTest) Run(t *testing.T) {
	tt.NewFilledTree(t)

	tt.TestBasic(t)
	tt.TestIterator(t)
	tt.TestChannels(t)
	tt.TestSlice(t)
	tt.TestRemove(t)
	tt.TestClear(t)
	tt.TestCopy(t)
	tt.TestContains(t)
	tt.TestGetSet(t)
	tt.TestFuncs(t)
}

// FillTree fills up a given tree with V
func (tt *TreeTest) FillTree(t *testing.T, tr Tree) {
	for i, va := range VRaw {
		tr.Insert(va)

		Equal(t, tr.Len(), i+1)

		vr, ok := tr.Get(va)
		True(t, ok)
		Equal(t, vr, va)
	}

	Equal(t, tr.Len(), VLen)

	n, ok := tr.First()
	True(t, ok)
	Equal(t, n, V[0])
	n, ok = tr.Last()
	True(t, ok)
	Equal(t, n, V[VLen-1])
}

// NewFilledTree creates a new tree and calls FillTree on it
func (tt *TreeTest) NewFilledTree(t *testing.T) Tree {
	tr := tt.New(t)

	tt.FillTree(t, tr)

	return tr
}

// TestBasic tests basic tree functionality
func (tt *TreeTest) TestBasic(t *testing.T) {
	tr := tt.New(t)

	Equal(t, tr.Len(), 0)
	True(t, tr.Empty())
	n, ok := tr.First()
	False(t, ok)
	Nil(t, n)
	n, ok = tr.Last()
	False(t, ok)
	Nil(t, n)
	n, ok = tr.Pop()
	Nil(t, n)
	False(t, ok)
	n, ok = tr.Shift()
	Nil(t, n)
	False(t, ok)

	tt.FillTree(t, tr)

	i := 0
	iter := tr.Iter()
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
	n, ok = tr.Pop()

	for i > -1 && n != nil {
		Equal(t, V[i], n)
		True(t, ok)
		Equal(t, tr.Len(), i)
		if i == 0 {
			True(t, tr.Empty())
		} else {
			False(t, tr.Empty())
		}

		i--
		n, ok = tr.Pop()
	}

	Equal(t, i, -1)
	Nil(t, n)
	False(t, ok)
	Equal(t, tr.Len(), 0)
	True(t, tr.Empty())

	tt.FillTree(t, tr)

	i = 0
	n, ok = tr.Shift()

	for i < VLen && n != nil {
		Equal(t, V[i], n)
		True(t, ok)
		Equal(t, tr.Len(), VLen-i-1)

		i++
		n, ok = tr.Shift()
	}

	Equal(t, i, VLen)
	Nil(t, n)
	Equal(t, tr.Len(), 0)
}

// TestIterator tests tree iterators
func (tt *TreeTest) TestIterator(t *testing.T) {
	// empty iterators
	tr := tt.New(t)

	Nil(t, tr.Iter())
	Nil(t, tr.IterBack())

	// one element
	tr.Insert(V[0])

	iter := tr.Iter()
	NotNil(t, iter)
	Equal(t, V[0], iter.Get())
	Nil(t, iter.Next())

	iter = tr.IterBack()
	NotNil(t, iter)
	Equal(t, V[0], iter.Get())
	Nil(t, iter.Previous())

	// full iterators
	tr = tt.NewFilledTree(t)

	i := 0

	for iter = tr.Iter(); iter != nil; iter = iter.Next() {
		Equal(t, iter.Get(), V[i])

		i++
	}

	Equal(t, i, VLen)

	tr = tt.NewFilledTree(t)

	i = VLen - 1

	for iter = tr.IterBack(); iter != nil; iter = iter.Previous() {
		Equal(t, iter.Get(), V[i])

		i--
	}

	Equal(t, i, -1)

	// iterate in wrong direction
	iter = tr.Iter()
	Nil(t, iter.Previous())

	iter = tr.IterBack()
	Nil(t, iter.Next())

	// iterate only within the left lane
	tr = tt.New(t)

	for i := 6; i > -1; i-- {
		tr.Insert(i)
	}

	iter = tr.Iter()

	for i := 0; i <= 6; i++ {
		Equal(t, iter.Get(), i)

		iter = iter.Next()
	}
	Nil(t, iter)

	// iterate only within the right lane
	tr = tt.New(t)

	for i := 0; i < 6; i++ {
		tr.Insert(i)
	}

	iter = tr.Iter()

	for i := 0; i < 6; i++ {
		Equal(t, iter.Get(), i)

		iter = iter.Next()
	}
	Nil(t, iter)

	// full tree
	testFullTree := func(cV []int) {
		max := len(cV) - 1

		tr = tt.New(t)

		for _, v := range cV {
			tr.Insert(v)
		}

		sort.Ints(cV)

		iter = tr.Iter()

		for _, v := range cV {
			Equal(t, iter.Get(), v)

			iter = iter.Next()
		}
		Nil(t, iter)

		// traverse back and forth
		iter = tr.Iter()
		i = 0

		for iter.Get() != max {
			Equal(t, iter.Get(), cV[i])

			iter = iter.Next()
			i++
		}

		for iter.Get() != 0 {
			Equal(t, iter.Get(), cV[i])

			iter = iter.Previous()
			i--
		}

		for iter.Get() != max {
			Equal(t, iter.Get(), cV[i])

			iter = iter.Next()
			i++
		}
		Nil(t, iter.Next())

		iter = tr.IterBack()
		i = len(cV) - 1

		for iter.Get() != 0 {
			Equal(t, iter.Get(), cV[i])

			iter = iter.Previous()
			i--
		}

		for iter.Get() != max {
			Equal(t, iter.Get(), cV[i])

			iter = iter.Next()
			i++
		}

		for iter.Get() != 0 {
			Equal(t, iter.Get(), cV[i])

			iter = iter.Previous()
			i--
		}
		Nil(t, iter.Previous())
	}

	testFullTree([]int{7, 3, 2, 0, 1, 5, 4, 6, 11, 9, 8, 10, 13, 12, 14})
	testFullTree([]int{8, 3, 1, 0, 2, 6, 5, 4, 7, 13, 10, 9, 11, 12, 15, 14, 16})
	testFullTree([]int{5, 1, 0, 4, 3, 2})

	// change direction in the middle of the tree
	tr = tt.NewFilledTree(t)

	iter = tr.Iter()

	for i := 1; i <= 3; i++ {
		Equal(t, iter.Get(), V[i-1])
		iter = iter.Next()
	}
	Equal(t, iter.Get(), V[3])

	for i := 4; i > 0; i-- {
		Equal(t, iter.Get(), V[i-1])
		iter = iter.Previous()
	}
	Nil(t, iter)

	iter = tr.IterBack()

	for i := 6; i > 3; i-- {
		Equal(t, iter.Get(), V[i-1])
		iter = iter.Previous()
	}
	Equal(t, iter.Get(), V[2])

	for i := 3; i <= 6; i++ {
		Equal(t, iter.Get(), V[i-1])
		iter = iter.Next()
	}
	Nil(t, iter)
}

// TestChannels tests tree channels
func (tt *TreeTest) TestChannels(t *testing.T) {
	// empty channels
	tr := tt.New(t)

	i := 0

	for v := range tr.Chan(0) {
		Equal(t, v, V[i])

		i++
	}

	Equal(t, i, 0)

	i = 0

	for v := range tr.ChanBack(0) {
		Equal(t, v, V[i])

		i++
	}

	Equal(t, i, 0)

	// one element
	tr.Insert(1)

	i = 0

	for v := range tr.Chan(0) {
		Equal(t, v, V[i])

		i++
	}

	Equal(t, i, 1)

	i = 0

	for v := range tr.ChanBack(0) {
		Equal(t, v, V[i])

		i++
	}

	Equal(t, i, 1)

	// full iterators
	tr = tt.NewFilledTree(t)

	i = 0

	for v := range tr.Chan(0) {
		Equal(t, v, V[i])

		i++
	}

	Equal(t, i, VLen)

	i = VLen - 1

	for v := range tr.ChanBack(0) {
		Equal(t, v, V[i])

		i--
	}

	Equal(t, i, -1)
}

// TestSlice tests converting the tree to slice
func (tt *TreeTest) TestSlice(t *testing.T) {
	tr := tt.New(t)
	Equal(t, tr.Slice(), []interface{}{})

	tt.FillTree(t, tr)
	Equal(t, tr.Slice(), V)

	tr.Shift()
	Equal(t, tr.Slice(), V[1:])

	tr.Pop()
	Equal(t, tr.Slice(), V[1:len(V)-1])
}

// TestRemove tests some remove methods
func (tt *TreeTest) TestRemove(t *testing.T) {
	tr := tt.NewFilledTree(t)

	// remove leaf
	v, ok := tr.Remove(4)
	True(t, ok)
	Equal(t, v, 4)
	Equal(t, tr.Slice(), []interface{}{1, 2, 3, 5, 6})

	// remove parent with left child
	v, ok = tr.Remove(3)
	True(t, ok)
	Equal(t, v, 3)
	Equal(t, tr.Slice(), []interface{}{1, 2, 5, 6})

	// remove parent with right child
	v, ok = tr.Remove(1)
	True(t, ok)
	Equal(t, v, 1)
	Equal(t, tr.Slice(), []interface{}{2, 5, 6})

	// remove parent with both childs
	v, ok = tr.Remove(5)
	True(t, ok)
	Equal(t, v, 5)
	Equal(t, tr.Slice(), []interface{}{2, 6})

	// remove last
	v, ok = tr.Remove(2)
	True(t, ok)
	Equal(t, v, 2)
	Equal(t, tr.Slice(), []interface{}{6})

	v, ok = tr.Remove(6)
	True(t, ok)
	Equal(t, v, 6)
	Equal(t, tr.Slice(), []interface{}{})

	// remove nothing
	v, ok = tr.Remove(-100)
	False(t, ok)
	v, ok = tr.Remove(100)
	False(t, ok)

	tr = tt.New(t)

	v, ok = tr.Remove(-100)
	False(t, ok)
	v, ok = tr.Remove(100)
	False(t, ok)

	tr = tt.NewFilledTree(t)

	v, ok = tr.Remove(-100)
	False(t, ok)
	v, ok = tr.Remove(100)
	False(t, ok)

	// prepare special cases
	tr = tt.New(t)

	for _, v := range []interface{}{4, 1, 3, 2, 5, 6, 12, 10, 7, 8, 9, 11} {
		tr.Insert(v)
	}

	// remove right child with left child
	v, ok = tr.Remove(3)
	True(t, ok)
	Equal(t, v, 3)
	Equal(t, tr.Slice(), []interface{}{1, 2, 4, 5, 6, 7, 8, 9, 10, 11, 12})

	// remove right child with right child
	v, ok = tr.Remove(6)
	True(t, ok)
	Equal(t, v, 6)
	Equal(t, tr.Slice(), []interface{}{1, 2, 4, 5, 7, 8, 9, 10, 11, 12})

	// remove with two children put removed right children at the end of left right children
	v, ok = tr.Remove(10)
	True(t, ok)
	Equal(t, v, 10)
	Equal(t, tr.Slice(), []interface{}{1, 2, 4, 5, 7, 8, 9, 11, 12})
}

// TestClear tests clearing the list
func (tt *TreeTest) TestClear(t *testing.T) {
	tr := tt.NewFilledTree(t)

	tr.Clear()

	Equal(t, tr.Len(), 0)
	n, ok := tr.First()
	False(t, ok)
	Nil(t, n)
	n, ok = tr.Last()
	False(t, ok)
	Nil(t, n)
	n, ok = tr.Pop()
	Nil(t, n)
	False(t, ok)
}

// TestCopy tests copying a list
func (tt *TreeTest) TestCopy(t *testing.T) {
	l1 := tt.NewFilledTree(t)

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

// TestContains tests contains methods
func (tt *TreeTest) TestContains(t *testing.T) {
	tr := tt.New(t)

	for _, vi := range V {
		ok := tr.Contains(vi)
		Equal(t, ok, false)
	}

	tr = tt.NewFilledTree(t)

	for _, vi := range V {
		ok := tr.Contains(vi)
		Equal(t, ok, true)
	}
}

// TestGetSet tests getters and setters
func (tt *TreeTest) TestGetSet(t *testing.T) {
	tr := tt.New(t)

	for i := range V {
		n, ok := tr.Get(V[i])

		False(t, ok)
		Nil(t, n)

		ok = tr.Set(V[i], i+10)

		False(t, ok)

		n, ok = tr.Get(i + 10)

		False(t, ok)
		Nil(t, n)
	}

	tt.FillTree(t, tr)

	for i := range V {
		n, ok := tr.Get(V[i])

		True(t, ok)
		Equal(t, n, V[i])

		ok = tr.Set(V[i], i+10)

		True(t, ok)

		n, ok = tr.Get(i + 10)

		True(t, ok)
		Equal(t, n, i+10)
	}
}

// TestFuncs tests all methods with functions as parameters
func (tt *TreeTest) TestFuncs(t *testing.T) {
	tr := tt.NewFilledTree(t)

	n, ok := tr.GetFunc(func(v interface{}) bool {
		return v == 2
	})
	Equal(t, V[1], n)
	True(t, ok)
	n, ok = tr.GetFunc(func(v interface{}) bool {
		return v == 100
	})
	Nil(t, nil)
	False(t, ok)

	True(t, tr.SetFunc(func(v interface{}) bool {
		return v == 4
	}, 99))
	Equal(t, tr.Slice(), []interface{}{1, 2, 3, 5, 6, 99})
	False(t, tr.SetFunc(func(v interface{}) bool {
		return v == 100
	}, 100))
	Equal(t, tr.Slice(), []interface{}{1, 2, 3, 5, 6, 99})
}
