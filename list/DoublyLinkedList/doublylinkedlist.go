package list

import (
	"errors"
)

// Node is a node of the list
type Node struct {
	next     *Node             // The node after this node in the list
	previous *Node             // The node before this node in the list
	list     *DoublyLinkedList // The list to which this element belongs
	Value    interface{}       // The value stored with this node
}

// Next returns the next node or nil
func (n *Node) Next() *Node {
	if i := n.next; n.list != nil {
		return i
	}

	return nil
}

// Previous returns the previous node or nil
func (n *Node) Previous() *Node {
	if i := n.previous; n.list != nil {
		return i
	}

	return nil
}

// DoublyLinkedList is a single linked list
type DoublyLinkedList struct {
	first *Node // The first node of the list
	last  *Node // The last node of the list
	len   int   // The current list length
}

// New returns an initialized list
func New() *DoublyLinkedList {
	return new(DoublyLinkedList).init()
}

// init initializes or clears the list
func (l *DoublyLinkedList) init() *DoublyLinkedList {
	l.Clear()

	return l
}

// Clear removes all nodes from the list
func (l *DoublyLinkedList) Clear() {
	i := l.first

	for i != nil {
		j := i.Next()

		i.list = nil
		i.next = nil
		i.previous = nil

		i = j
	}

	l.first = nil
	l.last = nil
	l.len = 0
}

// Len returns the curren list length
func (l *DoublyLinkedList) Len() int {
	return l.len
}

// First returns the first node of the list or nil
func (l *DoublyLinkedList) First() *Node {
	return l.first
}

// Last returns the last node of the list or nil
func (l *DoublyLinkedList) Last() *Node {
	return l.last
}

// Get returns the node with the given index or nil
func (l *DoublyLinkedList) Get(i int) (*Node, error) {
	if i < 0 || i >= l.len {
		return nil, errors.New("index bounds out of range")
	}

	j := 0

	for n := l.First(); n != nil; n = n.Next() {
		if i == j {
			return n, nil
		}

		j++
	}

	panic("there is something wrong with the internal structure")
}

// Set replaces the value in the list with the given value
func (l *DoublyLinkedList) Set(i int, v interface{}) error {
	if i < 0 || i >= l.len {
		return errors.New("index bounds out of range")
	}

	j := 0

	for n := l.First(); n != nil; n = n.Next() {
		if i == j {
			n.Value = v

			return nil
		}

		j++
	}

	panic("there is something wrong with the internal structure")
}

// Copy returns an exact copy of the list
func (l *DoublyLinkedList) Copy() *DoublyLinkedList {
	n := New()

	for i := l.First(); i != nil; i = i.Next() {
		n.Push(i.Value)
	}

	return n
}

// ToArray returns a copy of the list as slice
func (l *DoublyLinkedList) ToArray() []interface{} {
	a := make([]interface{}, l.len)

	j := 0

	for i := l.First(); i != nil; i = i.Next() {
		a[j] = i.Value

		j++
	}

	return a
}

// newNode initializes a new node for the list
func (l *DoublyLinkedList) newNode(v interface{}) *Node {
	return &Node{
		list:  l,
		Value: v,
	}
}

// InsertAfter creates a new node from a value, inserts it after a given node and returns the new one
func (l *DoublyLinkedList) InsertAfter(v interface{}, p *Node) *Node {
	if (p == nil && l.len != 0) || (p != nil && p.list != l) {
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

// InsertBefore creates a new node from a value, inserts it before a given node and returns the new one
func (l *DoublyLinkedList) InsertBefore(v interface{}, p *Node) *Node {
	if (p == nil && l.len != 0) || (p != nil && p.list != l) {
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

// InsertAt creates a new mnode from a value, inserts it at the exact index which must be in range of the list and returns the new node
func (l *DoublyLinkedList) InsertAt(i int, v interface{}) (*Node, error) {
	if i < 0 || i > l.len {
		return nil, errors.New("index bounds out of range")
	}

	if i == 0 {
		return l.Unshift(v), nil
	} else if i == l.len {
		return l.Push(v), nil
	}

	p, _ := l.Get(i)

	return l.InsertBefore(v, p), nil
}

// Remove removes a given node from the list
func (l *DoublyLinkedList) Remove(c *Node) *Node {
	if c == nil || c.list != l || l.len == 0 {
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

	c.list = nil
	c.next = nil
	c.previous = nil

	l.len--

	return c
}

// RemoveAt removes a node from the list at the given index
func (l *DoublyLinkedList) RemoveAt(i int) (*Node, error) {
	if i < 0 || i >= l.len {
		return nil, errors.New("index bounds out of range")
	}

	c, _ := l.Get(i)

	return l.Remove(c), nil
}

// RemoveFirstOccurrence removes the first node with the given value from the list and returns it or nil
func (l *DoublyLinkedList) RemoveFirstOccurrence(v interface{}) *Node {
	for i := l.First(); i != nil; i = i.Next() {
		if i.Value == v {
			return l.Remove(i)
		}
	}

	return nil
}

// RemoveLastOccurrence removes the last node with the given value from the list and returns it or nil
func (l *DoublyLinkedList) RemoveLastOccurrence(v interface{}) *Node {
	for i := l.Last(); i != nil; i = i.Previous() {
		if i.Value == v {
			return l.Remove(i)
		}
	}

	return nil
}

// Pop removes and returns the last node or nil
func (l *DoublyLinkedList) Pop() *Node {
	return l.Remove(l.last)
}

// Push creates a new node from a value, inserts it as the last node and returns it
func (l *DoublyLinkedList) Push(v interface{}) *Node {
	return l.InsertAfter(v, l.last)
}

// PushList adds the values of a list to the end of the list
func (l *DoublyLinkedList) PushList(l2 *DoublyLinkedList) {
	for i := l2.First(); i != nil; i = i.Next() {
		l.Push(i.Value)
	}
}

// Shift removes and returns the first node or nil
func (l *DoublyLinkedList) Shift() *Node {
	return l.Remove(l.first)
}

// Unshift creates a new node from a value, inserts it as the first node and returns it
func (l *DoublyLinkedList) Unshift(v interface{}) *Node {
	return l.InsertBefore(v, l.first)
}

// UnshiftList adds the values of a list to the front of the list
func (l *DoublyLinkedList) UnshiftList(l2 *DoublyLinkedList) {
	for i := l2.First(); i != nil; i = i.Next() {
		l.Unshift(i.Value)
	}
}

// Contains returns true if the value exists in the list
func (l *DoublyLinkedList) Contains(v interface{}) bool {
	_, ok := l.IndexOf(v)

	return ok
}

// IndexOf returns the first index of an occurence of the given value and true or -1 and false if the value does not exist
func (l *DoublyLinkedList) IndexOf(v interface{}) (int, bool) {
	i := 0

	for n := l.First(); n != nil; n = n.Next() {
		if n.Value == v {
			return i, true
		}

		i++
	}

	return -1, false
}

// LastIndexOf returns the last index of an occurence of the given value and true or -1 and false if the value does not exist
func (l *DoublyLinkedList) LastIndexOf(v interface{}) (int, bool) {
	i := l.len - 1

	for n := l.Last(); n != nil; n = n.Previous() {
		if n.Value == v {
			return i, true
		}

		i--
	}

	return -1, false
}

func (l *DoublyLinkedList) MoveAfter(n, p *Node) {
	if n.list != l || p.list != l || n == p {
		return
	}

	l.InsertAfter(l.Remove(n).Value, p)
}

func (l *DoublyLinkedList) MoveBefore(n, p *Node) {
	if n.list != l || p.list != l || n == p {
		return
	}

	l.InsertBefore(l.Remove(n).Value, p)
}

func (l *DoublyLinkedList) MoveToBack(n *Node) {
	l.MoveAfter(n, l.last)
}

func (l *DoublyLinkedList) MoveToFront(n *Node) {
	l.MoveBefore(n, l.first)
}
