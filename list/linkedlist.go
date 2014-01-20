package list

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
	return new(LinkedList).Init()
}

// Init initializes or clears the list
func (l *LinkedList) Init() *LinkedList {
	l.first = nil
	l.last = nil
	l.len = 0

	return l
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

	return nil
}

// insertAfter creates a new node from a value, inserts it after a given node and returns the new one
func (l *LinkedList) insertAfter(v interface{}, p *Node) *Node {
	n := l.newNode(v)

	// insert first node
	if p == nil {
		l.first = n
		l.last = n
	} else {
		p.next = n

		if p == l.last {
			l.last = n
		}
	}

	l.len++

	return n
}

// remove removes a given node from the list
func (l *LinkedList) remove(c *Node) *Node {
	if c == nil || c.list != l || l.len == 0 {
		return nil
	}

	r := c

	if c == l.first {
		l.first = c.next

		// c is the last node
		if c == l.last {
			l.last = nil
		}
	} else {
		p := l.findParent(c)

		p.next = c.next

		if c == l.last {
			l.last = p
		}
	}

	r.list = nil
	r.next = nil

	l.len--

	return r
}

// Pop removes and returns the last node or nil
func (l *LinkedList) Pop() *Node {
	return l.remove(l.last)
}

// Push creates a new node from a value, inserts it as the last node and returns it
func (l *LinkedList) Push(v interface{}) *Node {
	return l.insertAfter(v, l.last)
}
