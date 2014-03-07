package list

// Iterator defines a list iterator
type Iterator interface {
	// Next iterates to the next element in the list and returns the iterator, or nil if there is no next element
	Next() Iterator
	// Previous iterates to the previous element in the list and returns the iterator, or nil if there is no previous element
	Previous() Iterator

	// Get returns the value of the iterator's current element
	Get() interface{}
	// Set sets the value of the iterator's current element
	Set(v interface{})
}

// List defines a list
type List interface {
	// Clear resets the list to zero elements and resets the list's meta data
	Clear()
	// Len returns the current list length
	Len() int
	// Empty returns true if the current list length is zero
	Empty() bool

	// Chan returns a channel which iterates from the front to the back of the list
	Chan(n int) <-chan interface{}
	// ChanBack returns a channel which iterates from the back to the front of the list
	ChanBack(n int) <-chan interface{}

	// Iter returns an iterator which starts at the front of the list, or nil if there are no elements in the list
	Iter() Iterator
	// IterBack returns an iterator which starts at the back of the list, or nil if there are no elements in the list
	IterBack() Iterator

	// First returns the first value of the list and true, or false if there is no value
	First() (interface{}, bool)
	// Last returns the last value of the list and true, or false if there is no value
	Last() (interface{}, bool)
	// Get returns the value of the given index and nil, or an out of bound error if the index is incorrect
	Get(i int) (interface{}, error)
	// GetFunc returns the value of the first element selected by the given function and true, or false if there is no such element
	GetFunc(m func(v interface{}) bool) (interface{}, bool)
	// Set sets the value of the given index and returns nil, or an out of bound error if the index is incorrect
	Set(i int, v interface{}) error
	// SetFunc sets the value of the first element selected by the given function and returns true, or false if there is no such element
	SetFunc(m func(v interface{}) bool, v interface{}) bool
	// Swap swaps the value of index i with the value of index j
	Swap(i, j int)

	// Contains returns true if the value exists in the list, or false if it does not
	Contains(v interface{}) bool
	// IndexOf returns the first index of the given value and true, or false if it does not exists
	IndexOf(v interface{}) (int, bool)
	// LastIndexOf returns the last index of the given value and true, or false if it does not exists
	LastIndexOf(v interface{}) (int, bool)

	// Copy returns an exact copy of the list
	Copy() List
	// Slice returns a copy of the list as a slice
	Slice() []interface{}

	// Insert inserts a value into the list and returns nil, or an out of bound error if the index is incorrect
	Insert(i int, v interface{}) error
	// Remove removes and returns the value with the given index and nil, or an out of bound error if the index is incorrect
	Remove(i int) (interface{}, error)
	// RemoveFirstOccurrence removes the first occurrence of the given value in the list and returns true, or false if there is no such element
	RemoveFirstOccurrence(v interface{}) bool
	// RemoveLastOccurrence removes the last occurrence of the given value in the list and returns true, or false if there is no such element
	RemoveLastOccurrence(v interface{}) bool
	// Pop removes and returns the last element and true, or false if there is no such element
	Pop() (interface{}, bool)
	// Push inserts the given value at the end of the list
	Push(v interface{})
	// PushList pushes the given list
	PushList(l2 List)
	// Shift removes and returns the first element and true, or false if there is no such element
	Shift() (interface{}, bool)
	// Unshift inserts the given value at the beginning of the list
	Unshift(v interface{})
	// UnshiftList unshifts the given list
	UnshiftList(l2 List)

	// MoveAfter moves the element at index i after the element at index m and returns nil, or an out of bound error if an index is incorrect
	MoveAfter(i, m int) error
	// MoveToBack moves the element at index i to the back of the list and returns nil, or an out of bound error if the index is incorrect
	MoveToBack(i int) error
	// MoveBefore moves the element at index i before the element at index m and returns nil, or an out of bound error if an index is incorrect
	MoveBefore(i, m int) error
	// MoveToFront moves the element at index i to the front of the list and returns nil, or an out of bound error if the index is incorrect
	MoveToFront(i int) error
}
