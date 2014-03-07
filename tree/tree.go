package tree

// Iterator defines a tree iterator
type Iterator interface {
	// Next iterates to the next node in the tree and returns the iterator, or nil if there is no next node
	Next() Iterator
	// Previous iterates to the previous node in the tree and returns the iterator, or nil if there is no previous node
	Previous() Iterator

	// Get returns the value of the iterator's current node
	Get() interface{}
}

// Tree defines a tree
// Trees consists of nodes which are not exposed to the user. Only the values of each node is exposed.
// Trees are sorted by a compare function which also helps to identify nodes in the tree. This compare function makes use of the values of each node.
type Tree interface {
	// Clear resets the tree to zero nodes and resets the tree's meta data
	Clear()
	// Len returns the current node count
	Len() int
	// Empty returns true if the current node count is zero
	Empty() bool

	// Chan returns a channel which iterates from the front to the back of the tree
	Chan(n int) <-chan interface{}
	// ChanBack returns a channel which iterates from the back to the front of the tree
	ChanBack(n int) <-chan interface{}

	// Iter returns an iterator which starts at the front of the tree, or nil if there are no nodes in the tree
	Iter() Iterator
	// IterBack returns an iterator which starts at the back of the tree, or nil if there are no nodes in the tree
	IterBack() Iterator

	// First returns the first value of the tree and true, or false if there is no value
	First() (interface{}, bool)
	// Last returns the last value of the tree and true, or false if there is no value
	Last() (interface{}, bool)
	// Get returns the value of the node identified by the given id value and true, or false if there is no such node
	Get(id interface{}) (interface{}, bool)
	// GetFunc returns the value of the first node selected by the given function and true, or false if there is no such node
	GetFunc(m func(v interface{}) bool) (interface{}, bool)
	// Set sets the value of the node identified by the given id value and returns true, or false if there is no such node
	Set(id interface{}, v interface{}) bool
	// SetFunc sets the value of the first node selected by the given function and returns true, or false if there is no such node
	SetFunc(m func(v interface{}) bool, v interface{}) bool

	// Contains returns true if a node identified by the given id value exists in the tree, or false if it does not
	Contains(id interface{}) bool

	// Copy returns an exact copy of the tree
	Copy() Tree
	// Slice returns a copy of the tree as a slice
	Slice() []interface{}

	// Insert inserts a new node into the tree with the given value
	Insert(v interface{})
	// Remove removes the node identified by the given id value and returns its value and true, or false if there is no such node
	Remove(id interface{}) (interface{}, bool)
	// Pop removes the last node and returns its value and true, or false if there is no such node
	Pop() (interface{}, bool)
	// Shift removes the first node and returns its value and true, or false if there is no such node
	Shift() (interface{}, bool)
}
