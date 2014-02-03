package doublylinkedlist

import (
	"errors"

	"github.com/zimmski/container/list"
)

// node holds a single node of a doubly linked list
type node struct {
	next     *node       // The node after this node in the list
	previous *node       // The node before this node in the list
	value    interface{} // The value stored with this node
}

// iterator holds the iterator for a doubly linked list
type iterator struct {
	current *node // The current node in traversal
}

// Next iterates to the next element in the list and returns the iterator, or nil if there is no next element
func (iter *iterator) Next() list.Iterator {
	if iter.current != nil {
		iter.current = iter.current.next
	}

	if iter.current == nil {
		return nil
	}

	return iter
}

// Previous iterates to the previous element in the list and returns the iterator, or nil if there is no previous element
func (iter *iterator) Previous() list.Iterator {
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

// List holds a doubly linked list
type List struct {
	first *node // The first node of the list
	last  *node // The last node of the list
	len   int   // The current list length
}

// New returns a new doubly linked list
func New() *List {
	l := new(List)

	l.Clear()

	return l
}

// Clear resets the list to zero elements and resets the list's meta data
func (l *List) Clear() {
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
func (l *List) Len() int {
	return l.len
}

// newNode returns a new node for the list
func (l *List) newNode(v interface{}) *node {
	return &node{
		value: v,
	}
}

// getNode returns the node with the given index or nil
func (l *List) getNode(i int) (*node, error) {
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
func (l *List) insertNodeBefore(v interface{}, p *node) *node {
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
func (l *List) removeNode(c *node) interface{} {
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
func (l *List) newIterator(current *node) *iterator {
	return &iterator{
		current: current,
	}
}

// Chan returns a channel which iterates from the front to the back of the list
func (l *List) Chan(n int) <-chan interface{} {
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
func (l *List) ChanBack(n int) <-chan interface{} {
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
func (l *List) Iter() list.Iterator {
	if l.len == 0 {
		return nil
	}

	return l.newIterator(l.first)
}

// IterBack returns an iterator which starts at the back of the list, or nil if there are no elements in the list
func (l *List) IterBack() list.Iterator {
	if l.len == 0 {
		return nil
	}

	return l.newIterator(l.last)
}

// First returns the first value of the list and true, or false if there is no value
func (l *List) First() (interface{}, bool) {
	if l.len == 0 {
		return nil, false
	}

	return l.first.value, true
}

// Last returns the last value of the list and true, or false if there is no value
func (l *List) Last() (interface{}, bool) {
	if l.len == 0 {
		return nil, false
	}

	return l.last.value, true
}

// Get returns the value of the given index and nil, or an out of bound error if the index is incorrect
func (l *List) Get(i int) (interface{}, error) {
	n, err := l.getNode(i)

	if err != nil {
		return nil, err
	}

	return n.value, nil
}

// GetFunc returns the value of the first element selected by the given function and true, or false if there is no such element
func (l *List) GetFunc(m func(v interface{}) bool) (interface{}, bool) {
	for n := l.first; n != nil; n = n.next {
		if m(n.value) {
			return n.value, true
		}
	}

	return nil, false
}

// Set sets the value of the given index and returns nil, or an out of bound error if the index is incorrect
func (l *List) Set(i int, v interface{}) error {
	n, err := l.getNode(i)

	if err != nil {
		return err
	}

	n.value = v

	return nil
}

// SetFunc sets the value of the first element selected by the given function and returns true, or false if there is no such element
func (l *List) SetFunc(m func(v interface{}) bool, v interface{}) bool {
	for n := l.first; n != nil; n = n.next {
		if m(n.value) {
			n.value = v

			return true
		}
	}

	return false
}

// Swap swaps the value of index i with the value of index j
func (l *List) Swap(i, j int) {
	ni, erri := l.getNode(i)
	nj, errj := l.getNode(j)

	if erri == nil && errj == nil {
		ni.value, nj.value = nj.value, ni.value
	}
}

// Contains returns true if the value exists in the list, or false if it does not
func (l *List) Contains(v interface{}) bool {
	_, ok := l.IndexOf(v)

	return ok
}

// IndexOf returns the first index of the given value and true, or false if it does not exists
func (l *List) IndexOf(v interface{}) (int, bool) {
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
func (l *List) LastIndexOf(v interface{}) (int, bool) {
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
func (l *List) Copy() list.List {
	n := New()

	for i := l.first; i != nil; i = i.next {
		n.Push(i.value)
	}

	return n
}

// ToArray returns a copy of the list as slice
func (l *List) Slice() []interface{} {
	a := make([]interface{}, l.len)

	j := 0

	for i := l.first; i != nil; i = i.next {
		a[j] = i.value

		j++
	}

	return a
}

// InsertAt inserts a value into the list and returns nil, or an out of bound error if the index is incorrect
func (l *List) InsertAt(i int, v interface{}) error {
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

// RemoveAt removes and returns the value with the given index and nil, or an out of bound error if the index is incorrect
func (l *List) RemoveAt(i int) (interface{}, error) {
	if i < 0 || i >= l.len {
		return nil, errors.New("index bounds out of range")
	}

	c, _ := l.getNode(i)

	return l.removeNode(c), nil
}

// RemoveFirstOccurrence removes the first occurrence of the given value in the list and returns true, or false if there is no such element
func (l *List) RemoveFirstOccurrence(v interface{}) bool {
	for i := l.first; i != nil; i = i.next {
		if i.value == v {
			l.removeNode(i)

			return true
		}
	}

	return false
}

// RemoveLastOccurrence removes the last occurrence of the given value in the list and returns true, or false if there is no such element
func (l *List) RemoveLastOccurrence(v interface{}) bool {
	for i := l.last; i != nil; i = i.previous {
		if i.value == v {
			l.removeNode(i)

			return true
		}
	}

	return false
}

// Pop removes and returns the last element and true, or false if there is no such element
func (l *List) Pop() (interface{}, bool) {
	r := l.removeNode(l.last)

	return r, r != nil
}

// Push inserts the given value at the end of the list
func (l *List) Push(v interface{}) {
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
func (l *List) PushList(l2 list.List) {
	for iter := l2.Iter(); iter != nil; iter = iter.Next() {
		l.Push(iter.Get())
	}
}

// Shift removes and returns the first element and true, or false if there is no such element
func (l *List) Shift() (interface{}, bool) {
	r := l.removeNode(l.first)

	return r, r != nil
}

// Unshift inserts the given value at the beginning of the list
func (l *List) Unshift(v interface{}) {
	l.insertNodeBefore(v, l.first)
}

// UnshiftList unshifts the given list
func (l *List) UnshiftList(l2 list.List) {
	for iter := l2.Iter(); iter != nil; iter = iter.Next() {
		l.Unshift(iter.Get())
	}
}

// MoveAfter moves the element at index i after the element at index m and returns nil, or an out of bound error if an index is incorrect
func (l *List) MoveAfter(i, m int) error {
	if i < 0 || i >= l.len {
		return errors.New("i bounds out of range")
	} else if m < 0 || m >= l.len {
		return errors.New("m bounds out of range")
	}

	if i == m || i-1 == m {
		return nil
	}

	v, _ := l.RemoveAt(i)

	if i < m {
		m--
	}

	l.InsertAt(m+1, v)

	return nil
}

// MoveToBack moves the element at index i to the back of the list and returns nil, or an out of bound error if the index is incorrect
func (l *List) MoveToBack(i int) error {
	return l.MoveAfter(i, l.len-1)
}

// MoveBefore moves the element at index i before the element at index m and returns nil, or an out of bound error if an index is incorrect
func (l *List) MoveBefore(i, m int) error {
	if i < 0 || i >= l.len {
		return errors.New("i bounds out of range")
	} else if m < 0 || m >= l.len {
		return errors.New("m bounds out of range")
	}

	if i == m || i == m-1 {
		return nil
	}

	v, _ := l.RemoveAt(i)

	if i < m {
		m--
	}

	l.InsertAt(m, v)

	return nil
}

// MoveToFront moves the element at index i to the front of the list and returns nil, or an out of bound error if the index is incorrect
func (l *List) MoveToFront(i int) error {
	return l.MoveBefore(i, 0)
}
