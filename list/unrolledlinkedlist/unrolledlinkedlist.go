package unrolledlinkedlist

import (
	"errors"
)

type Node struct {
	next     *Node               // The node after this node in the list
	previous *Node               // The node before this node in the list
	list     *UnrolledLinkedList // The list to which this node belongs
	values   []interface{}       // The values stored with this node
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

type Iterator struct {
	current *Node // The current node in traversal
	i       int   // The current index of the current node
}

func newIterator(current *Node, i int) *Iterator {
	iter := new(Iterator)

	iter.i = i
	iter.current = current

	return iter
}

func (iter *Iterator) Next() bool {
	iter.i++

	if iter.current != nil && iter.i >= len(iter.current.values) {
		iter.i = 0
		iter.current = iter.current.Next()
	}

	return iter.current != nil
}

func (iter *Iterator) Set(v interface{}) {
	if iter.current == nil || iter.i < 0 || iter.i >= len(iter.current.values) {
		return
	}

	iter.current.values[iter.i] = v
}

func (iter *Iterator) Value() interface{} {
	if iter.current == nil || iter.i < 0 || iter.i >= len(iter.current.values) {
		return nil
	}

	return iter.current.values[iter.i]
}

type UnrolledLinkedList struct {
	first       *Node // The first node of the list
	last        *Node // The last node of the list
	maxElements int   // Maximum of elements per node
	len         int   // The current list length
}

// New returns an initialized list
func New(maxElements int) *UnrolledLinkedList {
	if maxElements < 1 {
		panic("maxElements must be at least 1")
	}

	l := new(UnrolledLinkedList)

	l.Clear()

	l.maxElements = maxElements

	return l
}

// Clear removes all nodes from the list
func (l *UnrolledLinkedList) Clear() {
	i := l.first

	for i != nil {
		j := i.Next()

		i.list = nil
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
func (l *UnrolledLinkedList) Len() int {
	return l.len
}

// newNode initializes a new node for the list
func (l *UnrolledLinkedList) newNode() *Node {
	return &Node{
		list:   l,
		values: make([]interface{}, 0, l.maxElements),
	}
}

func (l *UnrolledLinkedList) getNodeAt(i int) (*Node, int) {
	for c := l.first; c != nil; c = c.Next() {
		if i < len(c.values) {
			return c, i
		}

		i -= len(c.values)
	}

	return nil, -1
}

// insertNodeAfter creates a new node, inserts it after a given node and returns the new one
func (l *UnrolledLinkedList) insertNodeAfter(p *Node) *Node {
	if (p == nil && l.len != 0) || (p != nil && p.list != l) {
		return nil
	}

	n := l.newNode()

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

	return n
}

// InsertNodeBefore creates a new node, inserts it before a given node and returns the new one
func (l *UnrolledLinkedList) insertNodeBefore(p *Node) *Node {
	if (p == nil && l.len != 0) || (p != nil && p.list != l) {
		return nil
	}

	n := l.newNode()

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

	return n
}

// removeNode removes a given node from the list
func (l *UnrolledLinkedList) removeNode(c *Node) *Node {
	if c == nil || c.list != l {
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
	c.values = nil

	return c
}

// First returns the first node of the list or nil
func (l *UnrolledLinkedList) First() *Iterator {
	if l.len == 0 {
		return nil
	}

	return newIterator(l.first, 0)
}

// Last returns the last node of the list or nil
func (l *UnrolledLinkedList) Last() *Iterator {
	if l.len == 0 {
		return nil
	}

	return newIterator(l.last, len(l.last.values)-1)
}

// Get returns the node with the given index or nil
func (l *UnrolledLinkedList) Get(i int) (interface{}, error) {
	if i < 0 || i >= l.len {
		return nil, errors.New("index bounds out of range")
	}

	for c := l.first; c != nil; c = c.Next() {
		if i < len(c.values) {
			return c.values[i], nil
		}

		i -= len(c.values)
	}

	panic("there is something wrong with the internal structure")
}

// GetFunc returns the first node selected by a given function
func (l *UnrolledLinkedList) GetFunc(m func(v interface{}) bool) interface{} {
	iter := l.First()

	if iter != nil {
		for {
			if m(iter.Value()) {
				return iter.Value()
			}

			if !iter.Next() {
				break
			}
		}
	}

	return nil
}

// Set replaces the value in the list with the given value
func (l *UnrolledLinkedList) Set(i int, v interface{}) error {
	if i < 0 || i >= l.len {
		return errors.New("index bounds out of range")
	}

	for c := l.first; c != nil; c = c.Next() {
		if i < len(c.values) {
			c.values[i] = v

			return nil
		}

		i -= len(c.values)
	}

	panic("there is something wrong with the internal structure")
}

// SetFunc replaces the value of the first node selected by a given function
func (l *UnrolledLinkedList) SetFunc(m func(v interface{}) bool, v interface{}) {
	iter := l.First()

	if iter != nil {
		for {
			if m(iter.Value()) {
				iter.Set(v)

				return
			}

			if !iter.Next() {
				break
			}
		}
	}
}

// Copy returns an exact copy of the list
func (l *UnrolledLinkedList) Copy() *UnrolledLinkedList {
	n := New(l.maxElements)

	iter := l.First()

	if iter != nil {
		for {
			n.Push(iter.Value())

			if !iter.Next() {
				break
			}
		}
	}

	return n
}

// ToArray returns a copy of the list as slice
func (l *UnrolledLinkedList) ToArray() []interface{} {
	a := make([]interface{}, l.len)

	j := 0

	iter := l.First()

	if iter != nil {
		for {
			a[j] = iter.Value()

			if !iter.Next() {
				break
			}

			j++
		}
	}

	return a
}

/*

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

*/

// RemoveAt removes a node from the list at the given index
func (l *UnrolledLinkedList) RemoveAt(i int) (interface{}, error) {
	if i < 0 || i >= l.len {
		return nil, errors.New("index bounds out of range")
	}

	c, ic := l.getNodeAt(i)

	v := c.values[ic]

	for ; ic < len(c.values)-1; ic++ {
		c.values[ic] = c.values[ic+1]
	}

	c.values = c.values[:len(c.values)-1]

	l.len--

	if len(c.values) == 0 {
		l.removeNode(c)
	}

	return v, nil
}

/*

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

*/

// Pop removes and returns the last value and true or nil and false
func (l *UnrolledLinkedList) Pop() (interface{}, bool) {
	if l.len == 0 {
		return nil, false
	}

	v, _ := l.RemoveAt(l.len - 1)

	return v, true
}

// Push creates a new node from a value and inserts it as the last node
func (l *UnrolledLinkedList) Push(v interface{}) {
	if l.last == nil || len(l.last.values) == cap(l.last.values) {
		l.insertNodeAfter(l.last)
	}

	l.last.values = append(l.last.values, v)

	l.len++
}

// PushList adds the values of a list to the end of the list
func (l *UnrolledLinkedList) PushList(l2 *UnrolledLinkedList) {
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

// Shift removes and returns the first value and true or nil and false
func (l *UnrolledLinkedList) Shift() (interface{}, bool) {
	if l.len == 0 {
		return nil, false
	}

	v, _ := l.RemoveAt(0)

	return v, true
}

// Unshift creates a new node from a value and inserts it as the first node
func (l *UnrolledLinkedList) Unshift(v interface{}) {
	l.insertNodeBefore(l.first)

	l.first.values = append(l.first.values, v)

	l.len++
}

// UnshiftList adds the values of a list to the front of the list
func (l *UnrolledLinkedList) UnshiftList(l2 *UnrolledLinkedList) {
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
func (l *UnrolledLinkedList) Contains(v interface{}) bool {
	_, ok := l.IndexOf(v)

	return ok
}

// IndexOf returns the first index of an occurence of the given value and true or -1 and false if the value does not exist
func (l *UnrolledLinkedList) IndexOf(v interface{}) (int, bool) {
	i := 0

	for n := l.first; n != nil; n = n.Next() {
		for _, c := range n.values {
			if c == v {
				return i, true
			}

			i++
		}
	}

	return -1, false
}

// LastIndexOf returns the last index of an occurence of the given value and true or -1 and false if the value does not exist
func (l *UnrolledLinkedList) LastIndexOf(v interface{}) (int, bool) {
	i := l.len - 1

	for n := l.last; n != nil; n = n.Previous() {
		for j := len(n.values) - 1; j > -1; j-- {
			if n.values[j] == v {
				return i, true
			}

			i--
		}
	}

	return -1, false
}
