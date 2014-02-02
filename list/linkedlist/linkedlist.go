package linkedlist

import (
	"errors"

	"github.com/zimmski/container/list"
)

// node is a node of the list
type node struct {
	next  *node       // The node after this node in the list
	value interface{} // The value stored with this node
}

// element holds one value coming from an linked list
type element struct {
	value interface{}
}

// Value returns the value hold in the element
func (e *element) Value() interface{} {
	return e.value
}

// iterator is an iterator for an linked list
type iterator struct {
	current *node // The current node in traversal
	list    *List // The list to which this iterator belongs
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

// List is a single linked list
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

// findParentNode returns the parent to a given node or nil
func (l *List) findParentNode(c *node) *node {
	if c == nil {
		return nil
	}

	var p *node

	for i := l.first; i != nil; i = i.next {
		if i == c {
			return p
		}

		p = i
	}

	panic("there is something wrong with the internal structure")
}

// getNode returns the node with the given index or nil
func (l *List) getNode(i int) (*node, error) {
	if i < 0 || i >= l.len {
		return nil, errors.New("index bounds out of range")
	}

	j := 0

	for n := l.first; n != nil; n = n.next {
		if i == j {
			return n, nil
		}

		j++
	}

	panic("there is something wrong with the internal structure")
}

// remove removes a given node from the list using the provided parent p
func (l *List) removeNode(c *node, p *node) list.Element {
	if c == nil || l.len == 0 {
		return nil
	}

	if c == l.first {
		l.first = c.next

		// c is the last node
		if c == l.last {
			l.last = nil
		}
	} else {
		if p == nil {
			p = l.findParentNode(c)
		}

		p.next = c.next

		if c == l.last {
			l.last = p
		}
	}

	c.next = nil

	l.len--

	return &element{
		value: c.value,
	}
}

func (l *List) newiterator(current *node) list.Iterator {
	return &iterator{
		current: current,
		list:    l,
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
func (l *List) Get(i int) (list.Element, error) {
	n, err := l.getNode(i)

	if err != nil {
		return nil, err
	}

	return &element{
		value: n.value,
	}, nil
}

// GetFunc returns the first node selected by a given function
func (l *List) GetFunc(m func(v interface{}) bool) list.Element {
	for n := l.first; n != nil; n = n.next {
		if m(n.value) {
			return &element{
				value: n.value,
			}
		}
	}

	return nil
}

// Set replaces the value in the list with the given value
func (l *List) Set(i int, v interface{}) error {
	if i < 0 || i >= l.len {
		return errors.New("index bounds out of range")
	}

	j := 0

	for n := l.first; n != nil; n = n.next {
		if i == j {
			n.value = v

			return nil
		}

		j++
	}

	panic("there is something wrong with the internal structure")
}

// SetFunc replaces the value of the first node selected by a given function
func (l *List) SetFunc(m func(v interface{}) bool, v interface{}) {
	for n := l.first; n != nil; n = n.next {
		if m(n.value) {
			n.value = v

			return
		}
	}
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

// InsertAfter creates a new node from a value, inserts it after a given node and returns the new one
func (l *List) InsertAfter(v interface{}, p *node) *node {
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
		p.next = n

		if p == l.last {
			l.last = n
		}
	}

	l.len++

	return n
}

// InsertBefore creates a new node from a value, inserts it before a given node and returns the new one
func (l *List) InsertBefore(v interface{}, p *node) *node {
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
			pp := l.findParentNode(p)

			pp.next = n
		}

		n.next = p
	}

	l.len++

	return n
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

		l.InsertBefore(v, p)
	}

	return nil
}

/* TODO implement with Element
// Remove removes a given node from the list
func (l *List) Remove(c *node) list.Element {
	return l.removeNode(c, nil)
}
*/

// RemoveAt removes a node from the list at the given index
func (l *List) RemoveAt(i int) (list.Element, error) {
	switch {
	case i < 0 || i >= l.len:
		return nil, errors.New("index bounds out of range")
	case i == 0:
		return l.removeNode(l.first, nil), nil
	default:
		p, _ := l.getNode(i - 1)

		return l.removeNode(p.next, p), nil
	}
}

// RemoveFirstOccurrence removes the first node with the given value from the list and returns it or nil
func (l *List) RemoveFirstOccurrence(v interface{}) bool {
	var p *node

	for i := l.first; i != nil; i = i.next {
		if i.value == v {
			l.removeNode(i, p)

			return true
		}

		p = i
	}

	return false
}

// RemoveLastOccurrence removes the last node with the given value from the list and returns it or nil
func (l *List) RemoveLastOccurrence(v interface{}) bool {
	var c, p, pp *node

	for i := l.first; i != nil; i = i.next {
		if i.value == v {
			c = i
			p = pp
		}

		pp = i
	}

	if c != nil {
		l.removeNode(c, p)

		return true
	}

	return false
}

// Pop removes and returns the last node or nil
func (l *List) Pop() (list.Element, bool) {
	r := l.removeNode(l.last, nil)

	return r, r != nil
}

// Push creates a new node from a value, inserts it as the last node and returns it
func (l *List) Push(v interface{}) {
	l.InsertAfter(v, l.last)
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
func (l *List) Shift() (list.Element, bool) {
	r := l.removeNode(l.first, nil)

	return r, r != nil
}

// Unshift creates a new node from a value, inserts it as the first node and returns it
func (l *List) Unshift(v interface{}) {
	l.InsertBefore(v, l.first)
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
	i := 0
	j := -1

	for n := l.first; n != nil; n = n.next {
		if n.value == v {
			j = i
		}

		i++
	}

	return j, j != -1
}

/* TODO implement this with Element
// MoveAfter moves node n after node p
func (l *List) MoveAfter(n, p *node) {
	if n.list != l || p.list != l || n == p {
		return
	}

	l.InsertAfter(l.Remove(n).Value, p)
}

// MoveBefore moves node n before node p
func (l *List) MoveBefore(n, p *node) {
	if n.list != l || p.list != l || n == p {
		return
	}

	l.InsertBefore(l.Remove(n).Value, p)
}

// MoveToBack moves the given node after the last node of the list
func (l *List) MoveToBack(n *node) {
	l.MoveAfter(n, l.last)
}

// MoveToFront moves the given node before the first node of the list
func (l *List) MoveToFront(n *node) {
	l.MoveBefore(n, l.first)
}
*/
