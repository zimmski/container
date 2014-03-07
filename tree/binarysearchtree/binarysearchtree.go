package binarysearchtree

import (
	List "github.com/zimmski/container/list"
	dll "github.com/zimmski/container/list/doublylinkedlist"
	Tree "github.com/zimmski/container/tree"
)

// node holds a single node of a binary search tree
type node struct {
	parent *node       // The parent of this node
	left   *node       // The left child of this node
	right  *node       // The right child of this node
	value  interface{} // The value stored with this node
}

// iterator holds the iterator for a binary search tree
type iterator struct {
	tree    *tree     // The tree of this iterator
	current *node     // The current node in traversal
	stack   List.List // The stack holds the current state in the traversal
}

// Next iterates to the next node in the tree and returns the iterator, or nil if there is no next node
func (iter *iterator) Next() Tree.Iterator {
	if iter.current != nil {
		if iter.current.right != nil {
			iter.current = iter.current.right
			iter.stack.Push(iter.current)

			for iter.current.left != nil {
				iter.current = iter.current.left

				iter.stack.Push(iter.current)
			}
		}

		if iter.stack.Len() == 0 {
			iter.current = nil
		} else {
			c, _ := iter.stack.Pop()
			iter.current = c.(*node)
		}
	}

	if iter.current == nil {
		return nil
	}

	return iter
}

// Previous iterates to the previous node in the tree and returns the iterator, or nil if there is no previous node
func (iter *iterator) Previous() Tree.Iterator {
	if iter.current != nil {
		if iter.current.left != nil {
			iter.stack.Push(iter.current)

			iter.current = iter.current.left

			for iter.current.right != nil {
				iter.stack.Push(iter.current)

				iter.current = iter.current.right
			}
		} else {
			if iter.stack.Len() != 0 {
				// check if we stopped at the first element
				if c, _ := iter.stack.First(); c.(*node).parent == nil {
					var last, p *node

					it := iter.stack.Iter()

					for ; it != nil; it = it.Next() {
						last = it.Get().(*node)

						if last.parent != p || (p != nil && p.left != last) {
							break
						}

						p = last
					}

					// the stack contains only the left lane of the tree
					// and we stopped at the first element
					if it == nil && last.left == iter.current {
						iter.stack.Clear()
						iter.current = nil

						return nil
					}
				}

				if iter.tree.compare(iter.current.parent.value, iter.current.value) < 0 {
					iter.current = iter.current.parent

					if last, _ := iter.stack.Last(); iter.tree.compare(last.(*node).value, iter.current.value) == 0 {
						iter.stack.Pop()
					}

				} else {
					c, _ := iter.stack.Pop()

					if iter.stack.Len() != 0 {
						// we stopped at a leaf and we have to go one lane left
						// so we go up until we are at the junction of the two lanes
						if iter.tree.compare(c.(*node).value, iter.current.value) > 0 {
							for iter.tree.compare(c.(*node).parent.value, iter.current.value) > 0 {
								c, _ = iter.stack.Pop()
							}
							iter.stack.Pop()
						}
					}
					iter.current = c.(*node).parent
				}
			} else {
				iter.current = iter.current.parent
			}
		}
	}

	if iter.current == nil {
		return nil
	}

	return iter
}

// Get returns the value of the iterator's current node
func (iter *iterator) Get() interface{} {
	return iter.current.value
}

// tree holds a binary search tree
type tree struct {
	root    *node                      // The root node of the tree
	len     int                        // The current node count
	compare func(a, b interface{}) int // Compare two values for the tree node order
}

// New returns a new binary search tree
func New(compare func(a, b interface{}) int) *tree {
	t := new(tree)

	t.compare = compare

	t.Clear()

	return t
}

// Clear resets the tree to zero nodes and resets the tree's meta data
func (t *tree) Clear() {
	if t.len != 0 {
		stack := dll.New()

		stack.Push(t.root)

		for stack.Len() != 0 {
			cr, _ := stack.Pop()
			c := cr.(*node)

			c.parent = nil

			if c.right != nil {
				stack.Push(c.right)
				c.right = nil
			}
			if c.left != nil {
				stack.Push(c.left)
				c.left = nil
			}
		}
	}

	t.root = nil
	t.len = 0
}

// Len returns the current node count
func (t *tree) Len() int {
	return t.len
}

// Empty returns true if the current node count is zero
func (t *tree) Empty() bool {
	return t.len == 0
}

// newNode returns a new node for the tree
func (t *tree) newNode(v interface{}) *node {
	return &node{
		value: v,
	}
}

// getNode returns the node identified by the given id value, or nil if there is no such node
func (t *tree) getNode(id interface{}) *node {
	if t.len == 0 {
		return nil
	}

	c := t.root

	for {
		r := t.compare(id, c.value)

		if r == 0 {
			break
		} else if r < 0 {
			if c.left != nil {
				c = c.left
			} else {
				return nil
			}
		} else {
			if c.right != nil {
				c = c.right
			} else {
				return nil
			}
		}
	}

	return c
}

// getNodeFunc returns the first node selected by the given function, or nil if there is no such node
func (t *tree) getNodeFunc(m func(v interface{}) bool) *node {
	stack := dll.New()

	stack.Push(t.root)

	for stack.Len() != 0 {
		cr, _ := stack.Pop()
		c := cr.(*node)

		if m(c.value) {
			return c
		}

		if c.right != nil {
			stack.Push(c.right)
		}
		if c.left != nil {
			stack.Push(c.left)
		}
	}

	return nil
}

// getFirstNode returns the node with the first value of the tree
func (t *tree) getFirstNode() *node {
	if t.len == 0 {
		return nil
	}

	c := t.root

	for c.left != nil {
		c = c.left
	}

	return c
}

// getLastNode returns the node with the last value of the tree
func (t *tree) getLastNode() *node {
	if t.len == 0 {
		return nil
	}

	c := t.root

	for c.right != nil {
		c = c.right
	}

	return c
}

// insert creates a new node with the given value and adds the node accordingly to the tree
func (t *tree) insert(v interface{}) *node {
	n := t.newNode(v)

	if t.len == 0 {
		t.root = n
	} else {
		c := t.root

		for {
			if t.compare(n.value, c.value) <= 0 {
				if c.left != nil {
					c = c.left
				} else {
					c.left = n
					n.parent = c

					break
				}
			} else {
				if c.right != nil {
					c = c.right
				} else {
					c.right = n
					n.parent = c

					break
				}
			}
		}
	}

	t.len++

	return nil
}

// removeNode removes the given node from the tree
func (t *tree) removeNode(c *node) interface{} {
	if c == nil {
		return nil
	}

	if c.left == nil && c.right == nil {
		// no children
		if c.parent != nil {
			if c.parent.left == c {
				c.parent.left = nil
			} else {
				c.parent.right = nil
			}
		}

		if t.root == c {
			t.root = nil
		}
	} else if c.left != nil {
		// left child
		c.left.parent = c.parent

		if c.parent != nil {
			if c.parent.left == c {
				c.parent.left = c.left
			} else {
				c.parent.right = c.left
			}
		}

		if c.right != nil {
			// two children
			// put the right child as the rightest child of the left child
			r := c.left

			for {
				if r.right == nil {
					r.right = c.right
					c.right.parent = r

					break
				}

				r = r.right
			}
		}

		if t.root == c {
			t.root = c.left
		}
	} else {
		// right child
		c.right.parent = c.parent

		if c.parent != nil {
			if c.parent.left == c {
				c.parent.left = c.right
			} else {
				c.parent.right = c.right
			}
		}

		if t.root == c {
			t.root = c.right
		}
	}

	c.parent = nil
	c.left = nil
	c.right = nil

	t.len--

	return c.value
}

// Chan returns a channel which iterates from the front to the back of the tree
func (t *tree) Chan(n int) <-chan interface{} {
	ch := make(chan interface{})

	go func() {
		for iter := t.Iter(); iter != nil; iter = iter.Next() {
			ch <- iter.Get()
		}

		close(ch)
	}()

	return ch
}

// ChanBack returns a channel which iterates from the back to the front of the tree
func (t *tree) ChanBack(n int) <-chan interface{} {
	ch := make(chan interface{})

	go func() {
		for iter := t.IterBack(); iter != nil; iter = iter.Previous() {
			ch <- iter.Get()
		}

		close(ch)
	}()

	return ch
}

// Iter returns an iterator which starts at the front of the tree, or nil if there are no nodes in the tree
func (t *tree) Iter() Tree.Iterator {
	if t.len == 0 {
		return nil
	}

	iter := &iterator{
		tree:    t,
		current: t.root,
		stack:   dll.New(),
	}

	iter.stack.Push(iter.current)

	for iter.current.left != nil {
		iter.current = iter.current.left

		iter.stack.Push(iter.current)
	}

	c, _ := iter.stack.Pop()
	iter.current = c.(*node)

	return iter
}

// IterBack returns an iterator which starts at the back of the tree, or nil if there are no nodes in the tree
func (t *tree) IterBack() Tree.Iterator {
	if t.len == 0 {
		return nil
	}

	iter := &iterator{
		tree:    t,
		current: t.getLastNode(),
		stack:   dll.New(),
	}

	return iter
}

// First returns the first value of the tree and true, or false if there is no value
func (t *tree) First() (interface{}, bool) {
	if t.len == 0 {
		return nil, false
	}

	n := t.getFirstNode()

	return n.value, true
}

// Last returns the last value of the tree and true, or false if there is no value
func (t *tree) Last() (interface{}, bool) {
	if t.len == 0 {
		return nil, false
	}

	n := t.getLastNode()

	return n.value, true
}

// Get returns the value of the node identified by the given id value and true, or false if there is no such node
func (t *tree) Get(id interface{}) (interface{}, bool) {
	n := t.getNode(id)

	if n == nil {
		return nil, false
	}

	return n.value, true
}

// GetFunc returns the value of the first node selected by the given function and true, or false if there is no such node
func (t *tree) GetFunc(m func(v interface{}) bool) (interface{}, bool) {
	n := t.getNodeFunc(m)

	if n == nil {
		return nil, false
	}

	return n.value, true
}

// Set sets the value of the node identified by the given id value and returns true, or false if there is no such node
func (t *tree) Set(id interface{}, v interface{}) bool {
	n := t.getNode(id)

	if n == nil {
		return false
	}

	t.removeNode(n)

	t.Insert(v)

	return true
}

// SetFunc sets the value of the first node selected by the given function and returns true, or false if there is no such node
func (t *tree) SetFunc(m func(v interface{}) bool, v interface{}) bool {
	n := t.getNodeFunc(m)

	if n == nil {
		return false
	}

	t.removeNode(n)

	t.Insert(v)

	return true
}

// Contains returns true if a node identified by the given id value exists in the tree, or false if it does not
func (t *tree) Contains(id interface{}) bool {
	return t.getNode(id) != nil
}

// Copy returns an exact copy of the tree
func (t *tree) Copy() Tree.Tree {
	l2 := New(t.compare)

	stack := dll.New()

	stack.Push([3]*node{t.root, nil, nil})

	for stack.Len() != 0 {
		cr, _ := stack.Pop()
		c := cr.([3]*node)

		n := &node{
			value: c[0].value,
		}

		l2.len++

		if c[1] != nil {
			n.parent = c[1]
			c[1].left = n
		} else if c[2] != nil {
			n.parent = c[2]
			c[2].right = n
		} else {
			l2.root = n
		}

		if c[0].left != nil {
			stack.Push([3]*node{c[0].left, n, nil})
		}
		if c[0].right != nil {
			stack.Push([3]*node{c[0].right, nil, n})
		}
	}

	return l2
}

// Slice returns a copy of the tree as a slice
func (t *tree) Slice() []interface{} {
	a := make([]interface{}, t.len)

	j := 0

	for iter := t.Iter(); iter != nil; iter = iter.Next() {
		a[j] = iter.Get()

		j++
	}

	return a
}

// Insert inserts a new node into the tree with the given value
func (t *tree) Insert(v interface{}) {
	t.insert(v)
}

// Remove removes the node identified by the given id value and returns its value and true, or false if there is no such node
func (t *tree) Remove(id interface{}) (interface{}, bool) {
	n := t.getNode(id)

	if n == nil {
		return nil, false
	}

	return t.removeNode(n), true
}

// Pop removes the last node and returns its value and true, or false if there is no such node
func (t *tree) Pop() (interface{}, bool) {
	r := t.removeNode(t.getLastNode())

	return r, r != nil
}

// Shift removes the first node and returns its value and true, or false if there is no such node
func (t *tree) Shift() (interface{}, bool) {
	r := t.removeNode(t.getFirstNode())

	return r, r != nil
}
