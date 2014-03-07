package unrolledlinkedlist

import (
	"errors"

	List "github.com/zimmski/container/list"
)

// node holds a single node with values of a unrolled linked list
type node struct {
	next     *node         // The node after this node in the list
	previous *node         // The node before this node in the list
	values   []interface{} // The values stored with this node
}

// iterator holds the iterator for a doubly linked list
type iterator struct {
	current *node // The current node in traversal
	i       int   // The current index of the current node
}

// Next iterates to the next element in the list and returns the iterator, or nil if there is no next element
func (iter *iterator) Next() List.Iterator {
	iter.i++

	if iter.current != nil && iter.i >= len(iter.current.values) {
		iter.i = 0
		iter.current = iter.current.next
	}

	if iter.current == nil {
		return nil
	}

	return iter
}

// Previous iterates to the previous element in the list and returns the iterator, or nil if there is no previous element
func (iter *iterator) Previous() List.Iterator {
	iter.i--

	if iter.current != nil && iter.i < 0 {
		iter.current = iter.current.previous

		if iter.current != nil {
			iter.i = len(iter.current.values) - 1
		}
	}

	if iter.current == nil {
		return nil
	}

	return iter
}

// Get returns the value of the iterator's current element
func (iter *iterator) Get() interface{} {
	return iter.current.values[iter.i]
}

// Set sets the value of the iterator's current element
func (iter *iterator) Set(v interface{}) {
	iter.current.values[iter.i] = v
}

// list holds a unrolled linked list
type list struct {
	first       *node // The first node of the list
	last        *node // The last node of the list
	maxElements int   // Maximum of elements per node
	len         int   // The current list length
}

// New returns a new unrolled linked list
// @param maxElements defines how many elements should fit in a node
func New(maxElements int) *list {
	if maxElements < 1 {
		panic("maxElements must be at least 1")
	}

	l := new(list)

	l.Clear()

	l.maxElements = maxElements

	return l
}

// Clear resets the list to zero elements and resets the list's meta data
func (l *list) Clear() {
	i := l.first

	for i != nil {
		j := i.next

		i.next = nil
		i.previous = nil
		i.values = nil

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

// insertElement inserts the given value at index ic in the given node
func (l *list) insertElement(v interface{}, c *node, ic int) {
	if c == nil || ic == 0 || len(c.values) == 0 { // begin of node
		n := l.insertNode(c, false)

		n.values = append(n.values, v)
	} else if len(c.values) == ic { // end of node
		n := c

		if len(n.values) == cap(n.values) {
			n = l.insertNode(c, true)

			// move half of the old node if possible
			if l.maxElements > 3 {
				ic = (len(c.values) + 1) / 2
				n.values = append(n.values, c.values[ic:len(c.values)]...)
				c.values = c.values[:ic]
			}
		}

		n.values = append(n.values, v)
	} else { // "middle" of the node
		n := l.insertNode(c, true)

		n.values = append(n.values, c.values[ic:len(c.values)]...)
		c.values[ic] = v
		c.values = c.values[:ic+1]
	}

	l.len++
}

// removeElement removes the value at index ic in the given node
func (l *list) removeElement(c *node, ic int) interface{} {
	v := c.values[ic]

	for ; ic < len(c.values)-1; ic++ {
		c.values[ic] = c.values[ic+1]
	}

	c.values = c.values[:len(c.values)-1]

	l.len--

	if len(c.values) == 0 {
		l.removeNode(c)
	} else if n := c.next; l.maxElements > 3 && n != nil && len(c.values) < l.maxElements/2 {
		if len(n.values)-2 < l.maxElements/2 { // copy the next node into the current node
			for _, v := range n.values {
				c.values = append(c.values, v)
			}

			l.removeNode(n)
		} else { // copy 2 elements of the next node to the current node
			c.values = append(c.values, n.values[0], n.values[1])

			for ic = 2; ic < len(n.values); ic++ {
				n.values[ic-2] = n.values[ic]
			}

			n.values = n.values[:len(n.values)-2]
		}

	}

	return v
}

// newNode returns a new node for the list
func (l *list) newNode() *node {
	return &node{
		values: make([]interface{}, 0, l.maxElements),
	}
}

// getNode returns the node with the given value index and the elements index, or nil and -1 if there is no such element
func (l *list) getNode(i int) (*node, int) {
	for c := l.first; c != nil; c = c.next {
		if i < len(c.values) {
			return c, i
		}

		i -= len(c.values)
	}

	return nil, -1
}

// insertNode creates a new node from a value, inserts it after/before a given node and returns the new one
func (l *list) insertNode(p *node, after bool) *node {
	n := l.newNode()

	if l.len == 0 {
		l.first = n
		l.last = n
	} else if after {
		n.next = p.next
		if p.next != nil {
			p.next.previous = n
		}
		p.next = n
		n.previous = p

		if p == l.last {
			l.last = n
		}
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

	return n
}

// remove removes a given node from the list
func (l *list) removeNode(c *node) *node {
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
	c.values = nil

	return c
}

// newIterator returns a new iterator
func (l *list) newIterator(current *node, i int) *iterator {
	return &iterator{
		i:       i,
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

	return l.newIterator(l.first, 0)
}

// IterBack returns an iterator which starts at the back of the list, or nil if there are no elements in the list
func (l *list) IterBack() List.Iterator {
	if l.len == 0 {
		return nil
	}

	return l.newIterator(l.last, len(l.last.values)-1)
}

// First returns the first value of the list and true, or false if there is no value
func (l *list) First() (interface{}, bool) {
	if l.len == 0 {
		return nil, false
	}

	return l.first.values[0], true
}

// Last returns the last value of the list and true, or false if there is no value
func (l *list) Last() (interface{}, bool) {
	if l.len == 0 {
		return nil, false
	}

	return l.last.values[len(l.last.values)-1], true
}

// Get returns the value of the given index and nil, or an out of bound error if the index is incorrect
func (l *list) Get(i int) (interface{}, error) {
	if i > -1 && i < l.len {
		for c := l.first; c != nil; c = c.next {
			if i < len(c.values) {
				return c.values[i], nil
			}

			i -= len(c.values)
		}
	}

	return nil, errors.New("index bounds out of range")
}

// GetFunc returns the value of the first element selected by the given function and true, or false if there is no such element
func (l *list) GetFunc(m func(v interface{}) bool) (interface{}, bool) {
	for iter := l.Iter(); iter != nil; iter = iter.Next() {
		if m(iter.Get()) {
			return iter.Get(), true
		}
	}

	return nil, false
}

// Set sets the value of the given index and returns nil, or an out of bound error if the index is incorrect
func (l *list) Set(i int, v interface{}) error {
	if i > -1 && i < l.len {
		for c := l.first; c != nil; c = c.next {
			if i < len(c.values) {
				c.values[i] = v

				return nil
			}

			i -= len(c.values)
		}
	}

	return errors.New("index bounds out of range")
}

// SetFunc sets the value of the first element selected by the given function and returns true, or false if there is no such element
func (l *list) SetFunc(m func(v interface{}) bool, v interface{}) bool {
	for iter := l.Iter(); iter != nil; iter = iter.Next() {
		if m(iter.Get()) {
			iter.Set(v)

			return true
		}
	}

	return false
}

// Swap swaps the value of index i with the value of index j
func (l *list) Swap(i, j int) {
	ni, ici := l.getNode(i)
	nj, icj := l.getNode(j)

	if ni != nil && nj != nil {
		ni.values[ici], nj.values[icj] = nj.values[icj], ni.values[ici]
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
		for _, c := range n.values {
			if c == v {
				return i, true
			}

			i++
		}
	}

	return -1, false
}

// LastIndexOf returns the last index of the given value and true, or false if it does not exists
func (l *list) LastIndexOf(v interface{}) (int, bool) {
	i := l.len - 1

	for n := l.last; n != nil; n = n.previous {
		for j := len(n.values) - 1; j > -1; j-- {
			if n.values[j] == v {
				return i, true
			}

			i--
		}
	}

	return -1, false
}

// Copy returns an exact copy of the list
func (l *list) Copy() List.List {
	n := New(l.maxElements)

	for iter := l.Iter(); iter != nil; iter = iter.Next() {
		n.Push(iter.Get())
	}

	return n
}

// Slice returns a copy of the list as slice
func (l *list) Slice() []interface{} {
	a := make([]interface{}, l.len)

	j := 0

	for iter := l.Iter(); iter != nil; iter = iter.Next() {
		a[j] = iter.Get()

		j++
	}

	return a
}

// Insert inserts a value into the list and returns nil, or an out of bound error if the index is incorrect
func (l *list) Insert(i int, v interface{}) error {
	if i < 0 || i > l.len {
		return errors.New("index bounds out of range")
	}

	if i != l.len {
		c, ic := l.getNode(i)

		l.insertElement(v, c, ic)
	} else { // getNode returns nil for lastIndex + 1
		l.Push(v)
	}

	return nil
}

// Remove removes and returns the value with the given index and nil, or an out of bound error if the index is incorrect
func (l *list) Remove(i int) (interface{}, error) {
	if i < 0 || i >= l.len {
		return nil, errors.New("index bounds out of range")
	}

	return l.removeElement(l.getNode(i)), nil
}

// RemoveFirstOccurrence removes the first occurrence of the given value in the list and returns true, or false if there is no such element
func (l *list) RemoveFirstOccurrence(v interface{}) bool {
	for n := l.first; n != nil; n = n.next {
		for ic, c := range n.values {
			if c == v {
				l.removeElement(n, ic)

				return true
			}
		}
	}

	return false
}

// RemoveLastOccurrence removes the last occurrence of the given value in the list and returns true, or false if there is no such element
func (l *list) RemoveLastOccurrence(v interface{}) bool {
	for n := l.last; n != nil; n = n.previous {
		for ic := len(n.values) - 1; ic > -1; ic-- {
			if n.values[ic] == v {
				l.removeElement(n, ic)

				return true
			}
		}
	}

	return false
}

// Pop removes and returns the last element and true, or false if there is no such element
func (l *list) Pop() (interface{}, bool) {
	r, _ := l.Remove(l.len - 1)

	return r, r != nil
}

// Push inserts the given value at the end of the list
func (l *list) Push(v interface{}) {
	if l.last == nil {
		l.insertElement(v, nil, 0)
	} else {
		l.insertElement(v, l.last, len(l.last.values))
	}
}

// PushList pushes the given list
func (l *list) PushList(l2 List.List) {
	for iter := l2.Iter(); iter != nil; iter = iter.Next() {
		l.Push(iter.Get())
	}
}

// Shift removes and returns the first element and true, or false if there is no such element
func (l *list) Shift() (interface{}, bool) {
	r, _ := l.Remove(0)

	return r, r != nil
}

// Unshift inserts the given value at the beginning of the list
func (l *list) Unshift(v interface{}) {
	l.insertElement(v, l.first, 0)
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
