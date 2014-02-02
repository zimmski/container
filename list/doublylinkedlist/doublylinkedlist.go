package doublylinkedlist

import (
	"errors"

	"github.com/zimmski/container/list"
)

// node is a node of the list
type node struct {
	next     *node       // The node after this node in the list
	previous *node       // The node before this node in the list
	value    interface{} // The value stored with this node
}

// iterator is an iterator for an linked list
type iterator struct {
	current *node // The current node in traversal
}

// Next moves to the next node in the linked list and returns true or false if there is no next node
func (iter *iterator) Next() bool {
	if iter.current != nil {
		iter.current = iter.current.next
	}

	return iter.current != nil
}

// Set sets a value at the current position of the iterator
func (iter *iterator) Set(v interface{}) {
	if iter.current == nil {
		return
	}

	iter.current.value = v
}

// Value returns the value at the current position of the iterator
func (iter *iterator) Value() interface{} {
	if iter.current == nil {
		return nil
	}

	return iter.current.value
}

// List is a doubly linked list
type List struct {
	first *node // The first node of the list
	last  *node // The last node of the list
	len   int   // The current list length
}

// New returns an initialized list
func New() *List {
	l := new(List)

	l.Clear()

	return l
}

// Clear removes all nodes from the list
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

// newNode initializes a new node for the list
func (l *List) newNode(v interface{}) *node {
	return &node{
		value: v,
	}
}

// getNode returns the node with the given index or nil
func (l *List) getNode(i int) (*node, error) {
	if i > -1 || i < l.len {
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

// insertNodeAfter creates a new node from a value, inserts it after a given node and returns the new one
func (l *List) insertNodeAfter(v interface{}, p *node) *node {
	if p == nil && l.len != 0 {
		return nil
	}

	n := l.newNode(v)

	// insert first node
	if p == nil {
		l.first = n
		l.last = n
	} else {
		n.next = p.next
		if p.next != nil {
			p.next.previous = n
		}
		p.next = n
		n.previous = p

		if p == l.last {
			l.last = n
		}
	}

	l.len++

	return n
}

// insertNodeBefore creates a new node from a value, inserts it before a given node and returns the new one
func (l *List) insertNodeBefore(v interface{}, p *node) *node {
	if p == nil && l.len != 0 {
		return nil
	}

	n := l.newNode(v)

	// insert first node
	if p == nil {
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

// remove removes a given node from the list using the provided parent p
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

func (l *List) newiterator(current *node) *iterator {
	return &iterator{
		current: current,
	}
}

// First returns the first node of the list or nil
func (l *List) First() list.Iterator {
	if l.len == 0 {
		return nil
	}

	return l.newiterator(l.first)
}

// Last returns the last node of the list or nil
func (l *List) Last() list.Iterator {
	if l.len == 0 {
		return nil
	}

	return l.newiterator(l.last)
}

// Get returns the node with the given index or nil
func (l *List) Get(i int) (interface{}, error) {
	n, err := l.getNode(i)

	if err != nil {
		return nil, err
	}

	return n.value, nil
}

// GetFunc returns the first node selected by a given function
func (l *List) GetFunc(m func(v interface{}) bool) (interface{}, bool) {
	for n := l.first; n != nil; n = n.next {
		if m(n.value) {
			return n.value, true
		}
	}

	return nil, false
}

// Set replaces the value in the list with the given value
func (l *List) Set(i int, v interface{}) error {
	n, err := l.getNode(i)

	if err != nil {
		return err
	}

	n.value = v

	return nil
}

// SetFunc replaces the value of the first node selected by a given function
func (l *List) SetFunc(m func(v interface{}) bool, v interface{}) bool {
	for n := l.first; n != nil; n = n.next {
		if m(n.value) {
			n.value = v

			return true
		}
	}

	return false
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
func (l *List) ToArray() []interface{} {
	a := make([]interface{}, l.len)

	j := 0

	for i := l.first; i != nil; i = i.next {
		a[j] = i.value

		j++
	}

	return a
}

// InsertAt creates a new mnode from a value, inserts it at the exact index which must be in range of the list and returns the new node
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

// RemoveAt removes a node from the list at the given index
func (l *List) RemoveAt(i int) (interface{}, error) {
	if i < 0 || i >= l.len {
		return nil, errors.New("index bounds out of range")
	}

	c, _ := l.getNode(i)

	return l.removeNode(c), nil
}

// RemoveFirstOccurrence removes the first node with the given value from the list and returns it or nil
func (l *List) RemoveFirstOccurrence(v interface{}) bool {
	for i := l.first; i != nil; i = i.next {
		if i.value == v {
			l.removeNode(i)

			return true
		}
	}

	return false
}

// RemoveLastOccurrence removes the last node with the given value from the list and returns it or nil
func (l *List) RemoveLastOccurrence(v interface{}) bool {
	for i := l.last; i != nil; i = i.previous {
		if i.value == v {
			l.removeNode(i)

			return true
		}
	}

	return false
}

// Pop removes and returns the last node or nil
func (l *List) Pop() (interface{}, bool) {
	r := l.removeNode(l.last)

	return r, r != nil
}

// Push creates a new node from a value, inserts it as the last node and returns it
func (l *List) Push(v interface{}) {
	l.insertNodeAfter(v, l.last)
}

// PushList adds the values of a list to the end of the list
func (l *List) PushList(l2 list.List) {
	iter := l2.First()

	if iter == nil {
		return
	}

	for {
		l.Push(iter.Value())

		if !iter.Next() {
			break
		}
	}
}

// Shift removes and returns the first node or nil
func (l *List) Shift() (interface{}, bool) {
	r := l.removeNode(l.first)

	return r, r != nil
}

// Unshift creates a new node from a value, inserts it as the first node and returns it
func (l *List) Unshift(v interface{}) {
	l.insertNodeBefore(v, l.first)
}

// UnshiftList adds the values of a list to the front of the list
func (l *List) UnshiftList(l2 list.List) {
	iter := l2.First()

	if iter == nil {
		return
	}

	for {
		l.Unshift(iter.Value())

		if !iter.Next() {
			break
		}
	}
}

// Contains returns true if the value exists in the list
func (l *List) Contains(v interface{}) bool {
	_, ok := l.IndexOf(v)

	return ok
}

// IndexOf returns the first index of an occurence of the given value and true or -1 and false if the value does not exist
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

// LastIndexOf returns the last index of an occurence of the given value and true or -1 and false if the value does not exist
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
