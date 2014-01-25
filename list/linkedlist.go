package list

import (
	"errors"
)

// Node is a node of the list
type Node struct {
	next  *Node       // The node after this node in the list
	list  *LinkedList // The list to which this element belongs
	Value interface{} // The value stored with this node
}

// Next returns the next node or nil
func (n *Node) Next() *Node {
	if i := n.next; n.list != nil {
		return i
	}

	return nil
}

// LinkedList is a single linked list
type LinkedList struct {
	first *Node // The first node of the list
	last  *Node // The last node of the list
	len   int   // The current list length
}

// New returns an initialized list
func New() *LinkedList {
	return new(LinkedList).init()
}

// init initializes or clears the list
func (l *LinkedList) init() *LinkedList {
	l.Clear()

	return l
}

// Clear removes all nodes from the list
func (l *LinkedList) Clear() {
	i := l.first

	for i != nil {
		j := i.Next()

		i.list = nil
		i.next = nil

		i = j
	}

	l.first = nil
	l.last = nil
	l.len = 0
}

// Len returns the curren list length
func (l *LinkedList) Len() int {
	return l.len
}

// First returns the first node of the list or nil
func (l *LinkedList) First() *Node {
	return l.first
}

// Last returns the last node of the list or nil
func (l *LinkedList) Last() *Node {
	return l.last
}

// Get returns the node with the given index or nil
func (l *LinkedList) Get(i int) (*Node, error) {
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
func (l *LinkedList) Set(i int, v interface{}) error {
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
func (l *LinkedList) Copy() *LinkedList {
	n := New()

	for i := l.First(); i != nil; i = i.Next() {
		n.Push(i.Value)
	}

	return n
}

// ToArray returns a copy of the list as slice
func (l *LinkedList) ToArray() []interface{} {
	a := make([]interface{}, l.len)

	j := 0

	for i := l.First(); i != nil; i = i.Next() {
		a[j] = i.Value

		j++
	}

	return a
}

// newNode initializes a new node for the list
func (l *LinkedList) newNode(v interface{}) *Node {
	return &Node{
		list:  l,
		Value: v,
	}
}

// findParent returns the parent to a given node or nil
func (l *LinkedList) findParent(c *Node) *Node {
	if c == nil || c.list != l {
		return nil
	}

	var p *Node

	for i := l.First(); i != nil; i = i.Next() {
		if i == c {
			return p
		}

		p = i
	}

	panic("there is something wrong with the internal structure")
}

// InsertAfter creates a new node from a value, inserts it after a given node and returns the new one
func (l *LinkedList) InsertAfter(v interface{}, p *Node) *Node {
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
		p.next = n

		if p == l.last {
			l.last = n
		}
	}

	l.len++

	return n
}

// InsertBefore creates a new node from a value, inserts it before a given node and returns the new one
func (l *LinkedList) InsertBefore(v interface{}, p *Node) *Node {
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
			pp := l.findParent(p)

			pp.next = n
		}

		n.next = p
	}

	l.len++

	return n
}

// InsertAt creates a new mnode from a value, inserts it at the exact index which must be in range of the list and returns the new node
func (l *LinkedList) InsertAt(i int, v interface{}) (*Node, error) {
	if i < 0 || i > l.len {
		return nil, errors.New("index bounds out of range")
	}

	n := l.newNode(v)

	if i == 0 {
		n.next = l.first
		l.first = n
	} else if i == l.len {
		l.last.next = n
		l.last = n
	} else {
		p, _ := l.Get(i - 1)

		n.next = p.next
		p.next = n
	}

	l.len++

	return n, nil
}

// remove removes a given node from the list using the provided parent p
func (l *LinkedList) remove(c *Node, p *Node) *Node {
	if c == nil || c.list != l || l.len == 0 {
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
			p = l.findParent(c)
		}

		p.next = c.next

		if c == l.last {
			l.last = p
		}
	}

	c.list = nil
	c.next = nil

	l.len--

	return c
}

// Remove removes a given node from the list
func (l *LinkedList) Remove(c *Node) *Node {
	return l.remove(c, nil)
}

// RemoveAt removes a node from the list at the given index
func (l *LinkedList) RemoveAt(i int) (*Node, error) {
	switch {
	case i < 0 || i >= l.len:
		return nil, errors.New("index bounds out of range")
	case i == 0:
		return l.remove(l.first, nil), nil
	default:
		p, _ := l.Get(i - 1)

		return l.remove(p.next, p), nil
	}
}

// RemoveFirstOccurrence removes the first node with the given value from the list and returns it or nil
func (l *LinkedList) RemoveFirstOccurrence(v interface{}) *Node {
	var c, p *Node

	for i := l.First(); i != nil; i = i.Next() {
		if i.Value == v {
			c = i

			break
		}

		p = i
	}

	if c != nil {
		l.remove(c, p)
	}

	return c
}

// RemoveLastOccurrence removes the last node with the given value from the list and returns it or nil
func (l *LinkedList) RemoveLastOccurrence(v interface{}) *Node {
	var c, p, pp *Node

	for i := l.First(); i != nil; i = i.Next() {
		if i.Value == v {
			c = i
			p = pp
		}

		pp = i
	}

	if c != nil {
		l.remove(c, p)
	}

	return c
}

// Pop removes and returns the last node or nil
func (l *LinkedList) Pop() *Node {
	return l.Remove(l.last)
}

// Push creates a new node from a value, inserts it as the last node and returns it
func (l *LinkedList) Push(v interface{}) *Node {
	return l.InsertAfter(v, l.last)
}

// PushList adds the values of a list to the end of the list
func (l *LinkedList) PushList(l2 *LinkedList) {
	for i := l2.First(); i != nil; i = i.Next() {
		l.Push(i.Value)
	}
}

// Shift removes and returns the first node or nil
func (l *LinkedList) Shift() *Node {
	return l.Remove(l.first)
}

// Unshift creates a new node from a value, inserts it as the first node and returns it
func (l *LinkedList) Unshift(v interface{}) *Node {
	return l.InsertBefore(v, l.first)
}

// UnshiftList adds the values of a list to the front of the list
func (l *LinkedList) UnshiftList(l2 *LinkedList) {
	for i := l2.First(); i != nil; i = i.Next() {
		l.Unshift(i.Value)
	}
}

// Contains returns true if the value exists in the list
func (l *LinkedList) Contains(v interface{}) bool {
	_, ok := l.IndexOf(v)

	return ok
}

// IndexOf returns the first index of an occurence of the given value and true or -1 and false if the value does not exist
func (l *LinkedList) IndexOf(v interface{}) (int, bool) {
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
func (l *LinkedList) LastIndexOf(v interface{}) (int, bool) {
	i := 0
	j := -1

	for n := l.First(); n != nil; n = n.Next() {
		if n.Value == v {
			j = i
		}

		i++
	}

	return j, j != -1
}

func (l *LinkedList) MoveAfter(n, p *Node) {
	if n.list != l || p.list != l || n == p {
		return
	}

	l.InsertAfter(l.Remove(n).Value, p)
}

func (l *LinkedList) MoveBefore(n, p *Node) {
	if n.list != l || p.list != l || n == p {
		return
	}

	l.InsertBefore(l.Remove(n).Value, p)
}

func (l *LinkedList) MoveToBack(n *Node) {
	l.MoveAfter(n, l.last)
}

func (l *LinkedList) MoveToFront(n *Node) {
	l.MoveBefore(n, l.first)
}
