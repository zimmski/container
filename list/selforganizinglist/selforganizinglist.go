package selforganizinglist

import (
	"errors"

	List "github.com/zimmski/container/list"
)

// node holds a single node of a self organizing list
type node struct {
	next     *node       // The node after this node in the list
	previous *node       // The node before this node in the list
	value    interface{} // The value stored with this node
	meta     interface{} // Holds meta data of the node for sorting
}

// iterator holds the iterator for a self organizing list
type iterator struct {
	current *node // The current node in traversal
}

// Next iterates to the next element in the list and returns the iterator, or nil if there is no next element
func (iter *iterator) Next() List.Iterator {
	if iter.current != nil {
		iter.current = iter.current.next
	}

	if iter.current == nil {
		return nil
	}

	return iter
}

// Previous iterates to the previous element in the list and returns the iterator, or nil if there is no previous element
func (iter *iterator) Previous() List.Iterator {
	if iter.current != nil {
		iter.current = iter.current.previous
	}

	if iter.current == nil {
		return nil
	}

	return iter
}

// Get returns the value of the iterator's current element
func (iter *iterator) Get() interface{} {
	return iter.current.value
}

// Set sets the value of the iterator's current element
func (iter *iterator) Set(v interface{}) {
	iter.current.value = v
}

// list holds a self organizing list
type list struct {
	first *node // The first node of the list
	last  *node // The last node of the list
	len   int   // The current list length

	insertNode func(c *node) *node // Is called when a new node is created
	accessNode func(c *node) *node // is called when a node gets accessed
	copyList   func() *list        // is called for copying the list skeletal
}

// new returns a new self organizing list skeletal
func newList() *list {
	l := new(list)

	l.Clear()

	return l
}

// NewCount returns a new self organizing list with "count" method
// The "count" method increments an access counter if a node gets accessed.
// Afterwards the list is sorted based on the nodes counters.
// To make this method less prone to burst accesses only nodes who are not
// the first node will get an increased.
func NewCount() *list {
	l := newList()

	l.insertNode = func(c *node) *node {
		c.meta = 0

		return c
	}
	l.accessNode = func(c *node) *node {
		if c != l.first {
			c.meta = c.meta.(int) + 1

			for c != l.first && c.meta.(int) >= c.previous.meta.(int) {
				// it is cheaper to just swap values
				c.previous.value, c.value = c.value, c.previous.value
				c.previous.meta, c.meta = c.meta, c.previous.meta

				c = c.previous
			}
		}

		return c
	}
	l.copyList = func() *list {
		return NewCount()
	}

	return l
}

// NewMoveToFront returns a new self organizing list with "move to front" method
// The "move to front" method puts a node to the front if it gets accessed.
func NewMoveToFront() *list {
	l := newList()

	l.insertNode = func(c *node) *node {
		return c
	}
	l.accessNode = func(c *node) *node {
		if c != l.first {
			l.removeNode(c)
			return l.insertNodeBefore(c.value, l.first)
		}

		return c
	}
	l.copyList = func() *list {
		return NewMoveToFront()
	}

	return l
}

// NewTranspose returns a new self organizing list with "transpose" method
// The "transpose" method swaps a node with its parent if it gets accessed.
func NewTranspose() *list {
	l := newList()

	l.insertNode = func(c *node) *node {
		return c
	}
	l.accessNode = func(c *node) *node {
		if c.previous != nil {
			// it is cheaper to just swap values
			c.previous.value, c.value = c.value, c.previous.value

			return c.previous
		}

		return c
	}
	l.copyList = func() *list {
		return NewTranspose()
	}

	return l
}

// Clear resets the list to zero elements and resets the list's meta data
func (l *list) Clear() {
	i := l.first

	for i != nil {
		j := i.next

		i.next = nil
		i.previous = nil

		i = j
	}

	l.first = nil
	l.last = nil
	l.len = 0
}

// Len returns the current list length
func (l *list) Len() int {
	return l.len
}

// Empty returns true if the current list length is zero
func (l *list) Empty() bool {
	return l.len == 0
}

// newNode returns a new node for the list
func (l *list) newNode(v interface{}) *node {
	c := &node{
		value: v,
	}

	return l.insertNode(c)
}

// getNode returns the node with the given index or nil
func (l *list) getNode(i int) (*node, error) {
	if i > -1 && i < l.len {
		j := 0

		for n := l.first; n != nil; n = n.next {
			if i == j {
				return n, nil
			}

			j++
		}
	}

	return nil, errors.New("index bounds out of range")
}

// insertNodeBefore creates a new node from a value, inserts it before a given node and returns the new one
func (l *list) insertNodeBefore(v interface{}, p *node) *node {
	n := l.newNode(v)

	if l.len == 0 {
		l.first = n
		l.last = n
	} else {
		if p == l.first {
			l.first = n
		} else {
			if p.previous != nil {
				p.previous.next = n
				n.previous = p.previous
			}
		}

		n.next = p
		p.previous = n
	}

	l.len++

	return n
}

// remove removes a given node from the list
func (l *list) removeNode(c *node) interface{} {
	if c == nil || l.len == 0 {
		return nil
	}

	if c == l.first {
		l.first = c.next
		if c.next != nil {
			c.next.previous = nil
		}

		// c is the last node
		if c == l.last {
			l.last = nil
		}
	} else {
		if c.previous != nil {
			c.previous.next = c.next

			if c.next != nil {
				c.next.previous = c.previous
			} else if c == l.last {
				l.last = c.previous
			}
		}
	}

	c.next = nil
	c.previous = nil

	l.len--

	return c.value
}

// newIterator returns a new iterator
func (l *list) newIterator(current *node) *iterator {
	return &iterator{
		current: current,
	}
}

// Chan returns a channel which iterates from the front to the back of the list
func (l *list) Chan(n int) <-chan interface{} {
	ch := make(chan interface{})

	go func() {
		for iter := l.Iter(); iter != nil; iter = iter.Next() {
			ch <- iter.Get()
		}

		close(ch)
	}()

	return ch
}

// ChanBack returns a channel which iterates from the back to the front of the list
func (l *list) ChanBack(n int) <-chan interface{} {
	ch := make(chan interface{})

	go func() {
		for iter := l.IterBack(); iter != nil; iter = iter.Previous() {
			ch <- iter.Get()
		}

		close(ch)
	}()

	return ch
}

// Iter returns an iterator which starts at the front of the list, or nil if there are no elements in the list
func (l *list) Iter() List.Iterator {
	if l.len == 0 {
		return nil
	}

	return l.newIterator(l.first)
}

// IterBack returns an iterator which starts at the back of the list, or nil if there are no elements in the list
func (l *list) IterBack() List.Iterator {
	if l.len == 0 {
		return nil
	}

	return l.newIterator(l.last)
}

// First returns the first value of the list and true, or false if there is no value
func (l *list) First() (interface{}, bool) {
	if l.len == 0 {
		return nil, false
	}

	return l.first.value, true
}

// Last returns the last value of the list and true, or false if there is no value
func (l *list) Last() (interface{}, bool) {
	if l.len == 0 {
		return nil, false
	}

	return l.last.value, true
}

// Get returns the value of the given index and nil, or an out of bound error if the index is incorrect
func (l *list) Get(i int) (interface{}, error) {
	n, err := l.getNode(i)

	if err != nil {
		return nil, err
	}

	return n.value, nil
}

// GetFunc returns the value of the first element selected by the given function and true, or false if there is no such element
func (l *list) GetFunc(m func(v interface{}) bool) (interface{}, bool) {
	for n := l.first; n != nil; n = n.next {
		if m(n.value) {
			c := l.accessNode(n)

			return c.value, true
		}
	}

	return nil, false
}

// Set sets the value of the given index and returns nil, or an out of bound error if the index is incorrect
func (l *list) Set(i int, v interface{}) error {
	n, err := l.getNode(i)

	if err != nil {
		return err
	}

	n.value = v

	return nil
}

// SetFunc sets the value of the first element selected by the given function and returns true, or false if there is no such element
func (l *list) SetFunc(m func(v interface{}) bool, v interface{}) bool {
	for n := l.first; n != nil; n = n.next {
		if m(n.value) {
			n.value = v

			l.accessNode(n)

			return true
		}
	}

	return false
}

// Swap swaps the value of index i with the value of index j
func (l *list) Swap(i, j int) {
	ni, erri := l.getNode(i)
	nj, errj := l.getNode(j)

	if erri == nil && errj == nil {
		ni.value, nj.value = nj.value, ni.value
	}
}

// Contains returns true if the value exists in the list, or false if it does not
func (l *list) Contains(v interface{}) bool {
	_, ok := l.IndexOf(v)

	return ok
}

// IndexOf returns the first index of the given value and true, or false if it does not exists
func (l *list) IndexOf(v interface{}) (int, bool) {
	i := 0

	for n := l.first; n != nil; n = n.next {
		if n.value == v {
			return i, true
		}

		i++
	}

	return -1, false
}

// LastIndexOf returns the last index of the given value and true, or false if it does not exists
func (l *list) LastIndexOf(v interface{}) (int, bool) {
	i := l.len - 1

	for n := l.last; n != nil; n = n.previous {
		if n.value == v {
			return i, true
		}

		i--
	}

	return -1, false
}

// Copy returns an exact copy of the list
func (l *list) Copy() List.List {
	n := l.copyList()

	for i := l.first; i != nil; i = i.next {
		n.Push(i.value)
	}

	return n
}

// Slice returns a copy of the list as slice
func (l *list) Slice() []interface{} {
	a := make([]interface{}, l.len)

	j := 0

	for i := l.first; i != nil; i = i.next {
		a[j] = i.value

		j++
	}

	return a
}

// Insert inserts a value into the list and returns nil, or an out of bound error if the index is incorrect
func (l *list) Insert(i int, v interface{}) error {
	if i < 0 || i > l.len {
		return errors.New("index bounds out of range")
	}

	if i == 0 {
		l.Unshift(v)
	} else if i == l.len {
		l.Push(v)
	} else {
		p, _ := l.getNode(i)

		l.insertNodeBefore(v, p)
	}

	return nil
}

// Remove removes and returns the value with the given index and nil, or an out of bound error if the index is incorrect
func (l *list) Remove(i int) (interface{}, error) {
	if i < 0 || i >= l.len {
		return nil, errors.New("index bounds out of range")
	}

	c, _ := l.getNode(i)

	return l.removeNode(c), nil
}

// RemoveFirstOccurrence removes the first occurrence of the given value in the list and returns true, or false if there is no such element
func (l *list) RemoveFirstOccurrence(v interface{}) bool {
	for i := l.first; i != nil; i = i.next {
		if i.value == v {
			l.removeNode(i)

			return true
		}
	}

	return false
}

// RemoveLastOccurrence removes the last occurrence of the given value in the list and returns true, or false if there is no such element
func (l *list) RemoveLastOccurrence(v interface{}) bool {
	for i := l.last; i != nil; i = i.previous {
		if i.value == v {
			l.removeNode(i)

			return true
		}
	}

	return false
}

// Pop removes and returns the last element and true, or false if there is no such element
func (l *list) Pop() (interface{}, bool) {
	r := l.removeNode(l.last)

	return r, r != nil
}

// Push inserts the given value at the end of the list
func (l *list) Push(v interface{}) {
	n := l.newNode(v)

	if l.len == 0 {
		l.first = n
	} else {
		n.previous = l.last
		l.last.next = n
	}

	l.last = n

	l.len++
}

// PushList pushes the given list
func (l *list) PushList(l2 List.List) {
	for iter := l2.Iter(); iter != nil; iter = iter.Next() {
		l.Push(iter.Get())
	}
}

// Shift removes and returns the first element and true, or false if there is no such element
func (l *list) Shift() (interface{}, bool) {
	r := l.removeNode(l.first)

	return r, r != nil
}

// Unshift inserts the given value at the beginning of the list
func (l *list) Unshift(v interface{}) {
	l.insertNodeBefore(v, l.first)
}

// UnshiftList unshifts the given list
func (l *list) UnshiftList(l2 List.List) {
	for iter := l2.Iter(); iter != nil; iter = iter.Next() {
		l.Unshift(iter.Get())
	}
}

// MoveAfter moves the element at index i after the element at index m and returns nil, or an out of bound error if an index is incorrect
func (l *list) MoveAfter(i, m int) error {
	if i < 0 || i >= l.len {
		return errors.New("i bounds out of range")
	} else if m < 0 || m >= l.len {
		return errors.New("m bounds out of range")
	}

	if i == m || i-1 == m {
		return nil
	}

	v, _ := l.Remove(i)

	if i < m {
		m--
	}

	l.Insert(m+1, v)

	return nil
}

// MoveToBack moves the element at index i to the back of the list and returns nil, or an out of bound error if the index is incorrect
func (l *list) MoveToBack(i int) error {
	return l.MoveAfter(i, l.len-1)
}

// MoveBefore moves the element at index i before the element at index m and returns nil, or an out of bound error if an index is incorrect
func (l *list) MoveBefore(i, m int) error {
	if i < 0 || i >= l.len {
		return errors.New("i bounds out of range")
	} else if m < 0 || m >= l.len {
		return errors.New("m bounds out of range")
	}

	if i == m || i == m-1 {
		return nil
	}

	v, _ := l.Remove(i)

	if i < m {
		m--
	}

	l.Insert(m, v)

	return nil
}

// MoveToFront moves the element at index i to the front of the list and returns nil, or an out of bound error if the index is incorrect
func (l *list) MoveToFront(i int) error {
	return l.MoveBefore(i, 0)
}
