package unrolledlinkedlist

import (
	"errors"

	"github.com/zimmski/container/list"
)

// node is holding values in the unrolled linked list
type node struct {
	next     *node         // The node after this node in the list
	previous *node         // The node before this node in the list
	values   []interface{} // The values stored with this node
}

// element holds one value coming from an unrolled linked list
type element struct {
	value interface{}
}

// Value returns the value hold in the element
func (e *element) Value() interface{} {
	return e.value
}

// iterator is an iterator for an unrolled linked list
type iterator struct {
	current *node // The current node in traversal
	i       int   // The current index of the current node
	list    *List // The list to which this iterator belongs
}

// Next moves to the next node in the unrolled linked list and returns true or false if there is no next node
func (iter *iterator) Next() bool {
	iter.i++

	if iter.current != nil && iter.i >= len(iter.current.values) {
		iter.i = 0
		iter.current = iter.current.next
	}

	return iter.current != nil
}

// Set sets a value at the current position of the iterator
func (iter *iterator) Set(v interface{}) {
	if iter.current == nil || iter.i < 0 || iter.i >= len(iter.current.values) {
		return
	}

	iter.current.values[iter.i] = v
}

// Value returns the value at the current position of the iterator
func (iter *iterator) Value() interface{} {
	if iter.current == nil || iter.i < 0 || iter.i >= len(iter.current.values) {
		return nil
	}

	return iter.current.values[iter.i]
}

// List is a unrolled linked list
type List struct {
	first       *node // The first node of the list
	last        *node // The last node of the list
	maxElements int   // Maximum of elements per node
	len         int   // The current list length
}

// New returns an initialized list
func New(maxElements int) *List {
	if maxElements < 1 {
		panic("maxElements must be at least 1")
	}

	l := new(List)

	l.Clear()

	l.maxElements = maxElements

	return l
}

// Clear removes all nodes from the list
func (l *List) Clear() {
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
func (l *List) Len() int {
	return l.len
}

func (l *List) insertElement(v interface{}, c *node, ic int) {
	if c == nil || ic == 0 || len(c.values) == 0 { // begin of node
		n := l.insertNodeBefore(l.first)

		n.values = append(n.values, v)
	} else if len(c.values) == ic { // end of node
		n := c

		if len(n.values) == cap(n.values) {
			n = l.insertNodeAfter(c)

			// move half of the old node if possible
			if l.maxElements > 3 {
				ic = (len(c.values) + 1) / 2
				n.values = append(n.values, c.values[ic:len(c.values)]...)
				c.values = c.values[:ic]
			}
		}

		n.values = append(n.values, v)
	} else { // "middle" of the node
		n := l.insertNodeAfter(c)

		n.values = append(n.values, c.values[ic:len(c.values)]...)
		c.values[ic] = v
		c.values = c.values[:ic+1]
	}

	l.len++
}

func (l *List) removeElement(c *node, ic int) list.Element {
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

	return &element{
		value: v,
	}
}

// newNode initializes a new node for the list
func (l *List) newNode() *node {
	return &node{
		values: make([]interface{}, 0, l.maxElements),
	}
}

func (l *List) getNodeAt(i int) (*node, int) {
	for c := l.first; c != nil; c = c.next {
		if i < len(c.values) {
			return c, i
		}

		i -= len(c.values)
	}

	return nil, -1
}

// insertNodeAfter creates a new node, inserts it after a given node and returns the new one
func (l *List) insertNodeAfter(p *node) *node {
	if p == nil && l.len != 0 {
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
func (l *List) insertNodeBefore(p *node) *node {
	if p == nil && l.len != 0 {
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
func (l *List) removeNode(c *node) *node {
	if c == nil {
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
	c.values = nil

	return c
}

func (l *List) newIterator(current *node, i int) list.Iterator {
	return &iterator{
		i:       i,
		current: current,
		list:    l,
	}
}

// First returns the first node of the list or nil
func (l *List) First() list.Iterator {
	if l.len == 0 {
		return nil
	}

	return l.newIterator(l.first, 0)
}

// Last returns the last node of the list or nil
func (l *List) Last() list.Iterator {
	if l.len == 0 {
		return nil
	}

	return l.newIterator(l.last, len(l.last.values)-1)
}

// Get returns the node with the given index or nil
func (l *List) Get(i int) (list.Element, error) {
	if i < 0 || i >= l.len {
		return nil, errors.New("index bounds out of range")
	}

	for c := l.first; c != nil; c = c.next {
		if i < len(c.values) {
			return &element{
				value: c.values[i],
			}, nil
		}

		i -= len(c.values)
	}

	panic("there is something wrong with the internal structure")
}

// GetFunc returns the first node selected by a given function
func (l *List) GetFunc(m func(v interface{}) bool) list.Element {
	iter := l.First()

	if iter != nil {
		for {
			if m(iter.Value()) {
				return &element{
					value: iter.Value(),
				}
			}

			if !iter.Next() {
				break
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

	for c := l.first; c != nil; c = c.next {
		if i < len(c.values) {
			c.values[i] = v

			return nil
		}

		i -= len(c.values)
	}

	panic("there is something wrong with the internal structure")
}

// SetFunc replaces the value of the first node selected by a given function
func (l *List) SetFunc(m func(v interface{}) bool, v interface{}) {
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
func (l *List) Copy() list.List {
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
func (l *List) ToArray() []interface{} {
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

// InsertAt creates a new mnode from a value, inserts it at the exact index which must be in range of the list and returns the new node
func (l *List) InsertAt(i int, v interface{}) error {
	if i < 0 || i > l.len {
		return errors.New("index bounds out of range")
	}

	if i != l.len {
		c, ic := l.getNodeAt(i)

		l.insertElement(v, c, ic)
	} else { // getNodeAt returns nil for lastIndex + 1
		l.Push(v)
	}

	return nil
}

// RemoveAt removes a node from the list at the given index
func (l *List) RemoveAt(i int) (list.Element, error) {
	if i < 0 || i >= l.len {
		return nil, errors.New("index bounds out of range")
	}

	return l.removeElement(l.getNodeAt(i)), nil
}

// RemoveFirstOccurrence removes the first node with the given value from the list and returns it or nil
func (l *List) RemoveFirstOccurrence(v interface{}) bool {
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

// RemoveLastOccurrence removes the last node with the given value from the list and returns it or nil
func (l *List) RemoveLastOccurrence(v interface{}) bool {
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

// Pop removes and returns the last value and true or nil and false
func (l *List) Pop() (list.Element, bool) {
	r, _ := l.RemoveAt(l.len - 1)

	return r, r != nil
}

// Push creates a new node from a value and inserts it as the last node
func (l *List) Push(v interface{}) {
	if l.last == nil {
		l.insertElement(v, nil, 0)
	} else {
		l.insertElement(v, l.last, len(l.last.values))
	}
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

// Shift removes and returns the first value and true or nil and false
func (l *List) Shift() (list.Element, bool) {
	r, _ := l.RemoveAt(0)

	return r, r != nil
}

// Unshift creates a new node from a value and inserts it as the first node
func (l *List) Unshift(v interface{}) {
	l.insertElement(v, l.first, 0)
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
func (l *List) LastIndexOf(v interface{}) (int, bool) {
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
