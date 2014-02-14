package tree

import (
	"testing"

	. "github.com/stretchr/testify/assert"
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

		i--
		n, ok = tr.Pop()
	}

	Equal(t, i, -1)
	Nil(t, n)
	False(t, ok)
	Equal(t, tr.Len(), 0)

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
	// TODO
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

// TestGetSet tests getters and setters
func (tt *TreeTest) TestGetSet(t *testing.T) {
	tr := tt.New(t)

	for i := range V {
		n, err := tr.Get(i)

		Nil(t, n)
		NotNil(t, err)

		err = tr.Set(i, i+10)

		NotNil(t, err)

		n, err = tr.Get(i)

		Nil(t, n)
		NotNil(t, err)
	}

	tt.FillTree(t, tr)

	for i, vi := range V {
		n, err := tr.Get(i)

		Equal(t, n, vi)
		Nil(t, err)

		err = tr.Set(i, i+10)

		Nil(t, err)

		n, err = tr.Get(i)

		Equal(t, n, i+10)
		Nil(t, err)
	}
}

// TestFuncs tests all methods with functions as parameters
func (tt *TreeTest) TestFuncs(t *testing.T) {
	tr := tt.NewFilledTree(t)

	n, ok := tr.GetFunc(func(v interface{}) bool {
		return v == "a"
	})
	Equal(t, V[1], n)
	True(t, ok)
	n, ok = tr.GetFunc(func(v interface{}) bool {
		return v == "z"
	})
	Nil(t, nil)
	False(t, ok)

	True(t, tr.SetFunc(func(v interface{}) bool {
		return v == 2
	}, 3))
	Equal(t, tr.Slice(), []interface{}{1, "a", 3, "b", 3, "c", 4, "d"})
	False(t, tr.SetFunc(func(v interface{}) bool {
		return v == "z"
	}, 4))
	Equal(t, tr.Slice(), []interface{}{1, "a", 3, "b", 3, "c", 4, "d"})
}
