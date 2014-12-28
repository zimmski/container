package selforganizinglist

import (
	"testing"

	. "github.com/zimmski/container/test/assert"

	List "github.com/zimmski/container/list"
)

func TestAll(t *testing.T) {
	lt := &List.ListTest{
		New: func(t *testing.T) List.List {
			return NewTranspose()
		},
	}

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
	lt.TestSwap(t)
	lt.TestMoves(t)

	// This test methods are affected by the rearranging methods
	//lt.TestFuncs(t)
}

func TestCount(t *testing.T) {
	// GetFunc
	l := NewCount()

	for i := 0; i < 5; i++ {
		l.Push(i)
	}

	Equal(t, l.Slice(), []interface{}{0, 1, 2, 3, 4})

	n, ok := l.GetFunc(func(v interface{}) bool {
		return v == 0
	})
	Equal(t, n, 0)
	True(t, ok)
	Equal(t, l.Slice(), []interface{}{0, 1, 2, 3, 4})

	n, ok = l.GetFunc(func(v interface{}) bool {
		return v == 4
	})
	Equal(t, n, 4)
	True(t, ok)
	Equal(t, l.Slice(), []interface{}{4, 0, 1, 2, 3})

	n, ok = l.GetFunc(func(v interface{}) bool {
		return v == 2
	})
	Equal(t, n, 2)
	True(t, ok)
	Equal(t, l.Slice(), []interface{}{2, 4, 0, 1, 3})

	n, ok = l.GetFunc(func(v interface{}) bool {
		return v == 2
	})
	Equal(t, n, 2)
	True(t, ok)
	Equal(t, l.Slice(), []interface{}{2, 4, 0, 1, 3})

	n, ok = l.GetFunc(func(v interface{}) bool {
		return v == "z"
	})
	Nil(t, n)
	False(t, ok)
	Equal(t, l.Slice(), []interface{}{2, 4, 0, 1, 3})

	// SetFunc
	l = NewCount()

	for i := 0; i < 5; i++ {
		l.Push(i)
	}

	Equal(t, l.Slice(), []interface{}{0, 1, 2, 3, 4})

	True(t, l.SetFunc(func(v interface{}) bool {
		return v == 0
	}, "null"))
	Equal(t, l.Slice(), []interface{}{"null", 1, 2, 3, 4})

	True(t, l.SetFunc(func(v interface{}) bool {
		return v == 4
	}, "vier"))
	Equal(t, l.Slice(), []interface{}{"vier", "null", 1, 2, 3})

	True(t, l.SetFunc(func(v interface{}) bool {
		return v == 2
	}, "zwei"))
	Equal(t, l.Slice(), []interface{}{"zwei", "vier", "null", 1, 3})

	True(t, l.SetFunc(func(v interface{}) bool {
		return v == "zwei"
	}, "zweihai"))
	Equal(t, l.Slice(), []interface{}{"zweihai", "vier", "null", 1, 3})

	False(t, l.SetFunc(func(v interface{}) bool {
		return v == "z"
	}, 4))
	Equal(t, l.Slice(), []interface{}{"zweihai", "vier", "null", 1, 3})

	// Copy
	l2 := l.Copy()
	Equal(t, l2.Slice(), []interface{}{"zweihai", "vier", "null", 1, 3})
}

func TestMoveToFront(t *testing.T) {
	// GetFunc
	l := NewMoveToFront()

	for i := 0; i < 5; i++ {
		l.Push(i)
	}

	Equal(t, l.Slice(), []interface{}{0, 1, 2, 3, 4})

	n, ok := l.GetFunc(func(v interface{}) bool {
		return v == 0
	})
	Equal(t, n, 0)
	True(t, ok)
	Equal(t, l.Slice(), []interface{}{0, 1, 2, 3, 4})

	n, ok = l.GetFunc(func(v interface{}) bool {
		return v == 4
	})
	Equal(t, n, 4)
	True(t, ok)
	Equal(t, l.Slice(), []interface{}{4, 0, 1, 2, 3})

	n, ok = l.GetFunc(func(v interface{}) bool {
		return v == 2
	})
	Equal(t, n, 2)
	True(t, ok)
	Equal(t, l.Slice(), []interface{}{2, 4, 0, 1, 3})

	n, ok = l.GetFunc(func(v interface{}) bool {
		return v == 2
	})
	Equal(t, n, 2)
	True(t, ok)
	Equal(t, l.Slice(), []interface{}{2, 4, 0, 1, 3})

	n, ok = l.GetFunc(func(v interface{}) bool {
		return v == "z"
	})
	Nil(t, n)
	False(t, ok)
	Equal(t, l.Slice(), []interface{}{2, 4, 0, 1, 3})

	// SetFunc
	l = NewMoveToFront()

	for i := 0; i < 5; i++ {
		l.Push(i)
	}

	Equal(t, l.Slice(), []interface{}{0, 1, 2, 3, 4})

	True(t, l.SetFunc(func(v interface{}) bool {
		return v == 0
	}, "null"))
	Equal(t, l.Slice(), []interface{}{"null", 1, 2, 3, 4})

	True(t, l.SetFunc(func(v interface{}) bool {
		return v == 4
	}, "vier"))
	Equal(t, l.Slice(), []interface{}{"vier", "null", 1, 2, 3})

	True(t, l.SetFunc(func(v interface{}) bool {
		return v == 2
	}, "zwei"))
	Equal(t, l.Slice(), []interface{}{"zwei", "vier", "null", 1, 3})

	True(t, l.SetFunc(func(v interface{}) bool {
		return v == "zwei"
	}, "zweihai"))
	Equal(t, l.Slice(), []interface{}{"zweihai", "vier", "null", 1, 3})

	False(t, l.SetFunc(func(v interface{}) bool {
		return v == "z"
	}, 4))
	Equal(t, l.Slice(), []interface{}{"zweihai", "vier", "null", 1, 3})

	// Copy
	l2 := l.Copy()
	Equal(t, l2.Slice(), []interface{}{"zweihai", "vier", "null", 1, 3})
}

func TestTranspose(t *testing.T) {
	// GetFunc
	l := NewTranspose()

	for i := 0; i < 5; i++ {
		l.Push(i)
	}

	Equal(t, l.Slice(), []interface{}{0, 1, 2, 3, 4})

	n, ok := l.GetFunc(func(v interface{}) bool {
		return v == 0
	})
	Equal(t, n, 0)
	True(t, ok)
	Equal(t, l.Slice(), []interface{}{0, 1, 2, 3, 4})

	n, ok = l.GetFunc(func(v interface{}) bool {
		return v == 4
	})
	Equal(t, n, 4)
	True(t, ok)
	Equal(t, l.Slice(), []interface{}{0, 1, 2, 4, 3})

	n, ok = l.GetFunc(func(v interface{}) bool {
		return v == 2
	})
	Equal(t, n, 2)
	True(t, ok)
	Equal(t, l.Slice(), []interface{}{0, 2, 1, 4, 3})

	n, ok = l.GetFunc(func(v interface{}) bool {
		return v == 2
	})
	Equal(t, n, 2)
	True(t, ok)
	Equal(t, l.Slice(), []interface{}{2, 0, 1, 4, 3})

	n, ok = l.GetFunc(func(v interface{}) bool {
		return v == "z"
	})
	Nil(t, n)
	False(t, ok)
	Equal(t, l.Slice(), []interface{}{2, 0, 1, 4, 3})

	// SetFunc
	l = NewTranspose()

	for i := 0; i < 5; i++ {
		l.Push(i)
	}

	Equal(t, l.Slice(), []interface{}{0, 1, 2, 3, 4})

	True(t, l.SetFunc(func(v interface{}) bool {
		return v == 0
	}, "null"))
	Equal(t, l.Slice(), []interface{}{"null", 1, 2, 3, 4})

	True(t, l.SetFunc(func(v interface{}) bool {
		return v == 4
	}, "vier"))
	Equal(t, l.Slice(), []interface{}{"null", 1, 2, "vier", 3})

	True(t, l.SetFunc(func(v interface{}) bool {
		return v == 2
	}, "zwei"))
	Equal(t, l.Slice(), []interface{}{"null", "zwei", 1, "vier", 3})

	True(t, l.SetFunc(func(v interface{}) bool {
		return v == "zwei"
	}, "zweihai"))
	Equal(t, l.Slice(), []interface{}{"zweihai", "null", 1, "vier", 3})

	False(t, l.SetFunc(func(v interface{}) bool {
		return v == "z"
	}, 4))
	Equal(t, l.Slice(), []interface{}{"zweihai", "null", 1, "vier", 3})

	// Copy
	l2 := l.Copy()
	Equal(t, l2.Slice(), []interface{}{"zweihai", "null", 1, "vier", 3})
}

func BenchmarkPushSequentiel(b *testing.B) {
	lb := &List.ListBenchmark{
		New: func(b *testing.B) List.List {
			return NewTranspose()
		},
	}

	lb.BenchmarkPushSequentiel(b)
}

func BenchmarkUnshiftSequentiel(b *testing.B) {
	lb := &List.ListBenchmark{
		New: func(b *testing.B) List.List {
			return NewTranspose()
		},
	}

	lb.BenchmarkUnshiftSequentiel(b)
}
