package container

// Iterator defines a container iterator
type Iterator interface {
	// Next iterates to the next element in the container and returns the iterator, or nil if there is no next element
	Next() Iterator
	// Previous iterates to the previous element in the container and returns the iterator, or nil if there is no previous element
	Previous() Iterator

	// Get returns the value of the iterator's current element
	Get() interface{}
}

// Container defines a container
type Container interface {
	// Clear resets the container to zero elements and resets the container's meta data
	Clear()
	// Len returns the current count of container elements
	Len() int
	// Empty returns true if the current count of container elements is zero
	Empty() bool

	// Chan returns a channel which iterates from the front to the back of the container
	Chan(n int) <-chan interface{}
	// ChanBack returns a channel which iterates from the back to the front of the container
	ChanBack(n int) <-chan interface{}

	// Iter returns an iterator which starts at the front of the container, or nil if there are no elements in the container
	Iter() Iterator
	// IterBack returns an iterator which starts at the back of the container, or nil if there are no elements in the container
	IterBack() Iterator

	// Contains returns true if an element identified by the given id value exists in the container, or false if it does not
	Contains(id interface{}) bool

	// Slice returns a copy of the container as a slice
	Slice() []interface{}
}
