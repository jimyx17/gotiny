package bst

// this only implementes a way to search pointers faster

type Node struct {
	Key   uint64
	Index uint64
	Left  *Node
	Right *Node
}

func (n *Node) Search(key uint64) *Node {
	// This is our base case. If n == nil, `key`
	// doesn't exist in our binary search tree.
	if n == nil {
		return nil
	}

	if n.Key < key { // move right
		return n.Right.Search(key)
	} else if n.Key > key { // move left
		return n.Left.Search(key)
	}

	// n.Key == key, we found it!
	return n
}

func (n *Node) Insert(key uint64, index uint64) {
	if n.Key < key {
		if n.Right == nil { // we found an empty spot, done!
			n.Right = &Node{Key: key, Index: index}
		} else { // look right
			n.Right.Insert(key, index)
		}
	} else if n.Key > key {
		if n.Left == nil { // we found an empty spot, done!
			n.Left = &Node{Key: key, Index: index}
		} else { // look left
			n.Left.Insert(key, index)
		}
	}
	// n.Key == key, don't need to do anything
}

func (n *Node) Delete(key uint64) *Node {
	// search for `key`
	if n.Key < key {
		n.Right = n.Right.Delete(key)
	} else if n.Key > key {
		n.Left = n.Left.Delete(key)
		// n.Key == `key`
	} else {
		if n.Left == nil { // just pouint64 to opposite node
			return n.Right
		} else if n.Right == nil { // just point to opposite node
			return n.Left
		}

		// if `n` has two children, you need to
		// find the next highest number that
		// should go in `n`'s position so that
		// the BST stays correct
		min := n.Right.Min()

		// we only update `n`'s key with min
		// instead of replacing n with the min
		// node so n's immediate children aren't orphaned
		n.Key = min.Key
		n.Right = n.Right.Delete(min.Key)
	}
	return n
}

func (n *Node) Min() *Node {
	if n.Left == nil {
		return n
	}
	return n.Left.Min()
}

func (n *Node) Max() *Node {
	if n.Right == nil {
		return n
	}
	return n.Right.Max()
}
