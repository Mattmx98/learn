package main

import "fmt"

// Tree contains a Root node of a binary search tree.
type Tree struct {
	Root *Node
}

// New returns a new Tree with its root Node.
func New(root *Node) *Tree {
	tr := &Tree{}
	root.Black = true
	tr.Root = root
	return tr
}

// Interface represents a single object in the tree.
type Interface interface {
	// Less returns true when the receiver item(key)
	// is less than the given(than) argument.
	Less(than Interface) bool
}

// Node is a Node and a Tree itself.
type Node struct {
	// Left is a left child Node.
	Left *Node

	Key   Interface
	Black bool // True when the color of parent link is black.
	// In Left-Leaning Red-Black tree, new nodes are always red
	// because the zero boolean value is false.
	// Null links are black.

	// Right is a right child Node.
	Right *Node
}

// NewNode returns a new Node.
func NewNode(key Interface) *Node {
	nd := &Node{}
	nd.Key = key
	nd.Black = false
	return nd
}

func (tr *Tree) String() string {
	return tr.Root.String()
}

func (nd *Node) String() string {
	if nd == nil {
		return "[]"
	}
	s := ""
	if nd.Left != nil {
		s += nd.Left.String() + " "
	}
	s += fmt.Sprintf("%v(%v)", nd.Key, nd.Black)
	if nd.Right != nil {
		s += " " + nd.Right.String()
	}
	return "[" + s + "]"
}

func isRed(nd *Node) bool {
	if nd == nil {
		return false
	}
	return !nd.Black
}

// insert inserts nd2 with nd1 as a root.
func (nd1 *Node) insert(nd2 *Node) *Node {
	if nd1 == nil {
		return nd2
	}
	if nd1.Key.Less(nd2.Key) {
		// nd1 is smaller than nd2
		// nd1 < nd2
		nd1.Right = nd1.Right.insert(nd2)
	} else {
		// nd1 is greater than nd2
		// nd1 >= nd2
		nd1.Left = nd1.Left.insert(nd2)
	}
	// Balance from nd1
	return Balance(nd1)
}

// Insert inserts a Node to a Tree without replacement.
// It does standard BST insert and colors the new link red.
// If the new red link is a right link, rotate left.
// If two left red links in a row, rotate to right and flip color.
// (https://youtu.be/lKmLBOJXZHI?t=20m43s)
//
// Note that it recursively balances from its parent nodes
// to the root node at the top.
//
// And make sure paint the Root black(not-red).
func (tr *Tree) Insert(nd *Node) {
	if tr.Root == nd {
		return
	}
	tr.Root = tr.Root.insert(nd)

	// Root node must be always black.
	tr.Root.Black = true
}

// RotateToLeft runs when there is a right-leaning link.
// tr.Root = RotateToLeft(tr.Root) overwrite the Root
// with the new top Node.
func RotateToLeft(nd *Node) *Node {
	if nd.Right.Black {
		panic("Can't rotate a black link")
	}

	// exchange x and nd
	// nd is parent node, x is Right child
	x := nd.Right
	nd.Right = x.Left
	x.Left = nd

	x.Black = nd.Black
	nd.Black = false

	return x
}

// RotateToRight runs when there are two left red links in a row.
// tr.Root = RotateToRight(tr.Root) overwrite the Root
// with the new top Node.
func RotateToRight(nd *Node) *Node {
	if nd.Left.Black {
		panic("Can't rotate a black link")
	}

	// exchange x and nd
	// nd is parent node, x is Left child
	x := nd.Left
	nd.Left = x.Right
	x.Right = nd

	x.Black = nd.Black
	nd.Black = false

	return x
}

// FlipColor flips the color.
// Left and Right children must be present
func FlipColor(nd *Node) {
	// nd is parent node
	nd.Black = !nd.Black
	nd.Left.Black = !nd.Left.Black
	nd.Right.Black = !nd.Right.Black
}

// MoveRedFromRightToLeft moves Red Node
// from Right sub-Tree to Left sub-Tree.
// Left and Right children must be present
func MoveRedFromRightToLeft(nd *Node) *Node {
	FlipColor(nd)
	if isRed(nd.Right.Left) {
		nd.Right = RotateToRight(nd.Right)
		nd = RotateToLeft(nd)
		FlipColor(nd)
	}
	return nd
}

// MoveRedFromLeftToRight moves Red Node
// from Left sub-Tree to Right sub-Tree.
// Left and Right children must be present
func MoveRedFromLeftToRight(nd *Node) *Node {
	FlipColor(nd)
	if isRed(nd.Left.Left) {
		nd = RotateToRight(nd)
		FlipColor(nd)
	}
	return nd
}

// Balance balances the Node.
func Balance(nd *Node) *Node {
	// nd is parent node
	if isRed(nd.Right) && !isRed(nd.Left) {
		nd = RotateToLeft(nd)
	}
	if isRed(nd.Left) && isRed(nd.Left.Left) {
		nd = RotateToRight(nd)
	}
	if isRed(nd.Left) && isRed(nd.Right) {
		FlipColor(nd)
	}
	return nd
}

// FixUp fixes the balances of the Node.
func FixUp(nd *Node) *Node {
	if isRed(nd.Right) {
		nd = RotateToLeft(nd)
	}
	if isRed(nd.Left) && isRed(nd.Left.Left) {
		nd = RotateToRight(nd)
	}
	if isRed(nd.Left) && isRed(nd.Right) {
		FlipColor(nd)
	}
	return nd
}

type Float64 float64

// Less returns true if float64(a) < float64(b).
func (a Float64) Less(b Interface) bool {
	return a < b.(Float64)
}

// Min returns the minimum key Node in the tree.
func (tr Tree) Min() *Node {
	nd := tr.Root
	if nd == nil {
		return nil
	}
	for nd.Left != nil {
		nd = nd.Left
	}
	return nd
}

// Max returns the maximum key Node in the tree.
func (tr *Tree) Max() *Node {
	nd := tr.Root
	if nd == nil {
		return nil
	}
	for nd.Right != nil {
		nd = nd.Right
	}
	return nd
}

// Search does binary-search on a given key and returns the first Node with the key.
func (tr Tree) Search(key Interface) *Node {
	nd := tr.Root
	// just updating the pointer value (address)
	for nd != nil {
		if nd.Key == nil {
			break
		}
		switch {
		case nd.Key.Less(key):
			nd = nd.Right
		case key.Less(nd.Key):
			nd = nd.Left
		default:
			return nd
		}
	}
	return nil
}

// SearchChan does binary-search on a given key and return the first Node with the key.
func (tr Tree) SearchChan(key Interface, ch chan *Node) {
	searchChan(tr.Root, key, ch)
	close(ch)
}

func searchChan(nd *Node, key Interface, ch chan *Node) {
	// leaf node
	if nd == nil {
		return
	}
	// when equal
	if !nd.Key.Less(key) && !key.Less(nd.Key) {
		ch <- nd
		return
	}
	searchChan(nd.Left, key, ch)  // left
	searchChan(nd.Right, key, ch) // right
}

// SearchParent does binary-search on a given key and returns the parent Node.
func (tr Tree) SearchParent(key Interface) *Node {
	nd := tr.Root
	parent := new(Node)
	parent = nil
	// just updating the pointer value (address)
	for nd != nil {
		if nd.Key == nil {
			break
		}
		switch {
		case nd.Key.Less(key):
			parent = nd // copy the pointer(address)
			nd = nd.Right
		case key.Less(nd.Key):
			parent = nd // copy the pointer(address)
			nd = nd.Left
		default:
			return parent
		}
	}
	return nil
}

func main() {
	root := NewNode(Float64(1))
	tr := New(root)
	nums := []float64{3, 9, 13, 17, 20, 25, 39, 16, 15, 2, 2.5}
	for _, num := range nums {
		tr.Insert(NewNode(Float64(num)))
	}

	fmt.Println("Deleting start!")
	fmt.Println("Deleted", tr.Delete(Float64(39)))
	fmt.Println(tr.Root.Left.Key)
	fmt.Println(tr.Root.Key)
	fmt.Println(tr.Root.Right.Key)
	fmt.Println()

	fmt.Println("Deleted", tr.Delete(Float64(20)))
	fmt.Println(tr.Root.Left.Key)
	fmt.Println(tr.Root.Key)
	fmt.Println(tr.Root.Right.Key)
	fmt.Println()

	/*
	   Deleting start!
	   Deleted 39
	   3
	   13
	   20

	   Deleted 20
	   3
	   13
	   16
	*/
}

// DeleteMin deletes the minimum-Key Node of the sub-Tree.
func DeleteMin(nd *Node) (*Node, Interface) {
	if nd == nil {
		return nil, nil
	}
	if nd.Left == nil {
		return nil, nd.Key
	}
	if !isRed(nd.Left) && !isRed(nd.Left.Left) {
		nd = MoveRedFromRightToLeft(nd)
	}
	var deleted Interface
	nd.Left, deleted = DeleteMin(nd.Left)
	return FixUp(nd), deleted
}

// DeleteMin deletes the minimum-Key Node of the Tree.
// It returns the minimum Key or nil.
func (tr *Tree) DeleteMin() Interface {
	var deleted Interface
	tr.Root, deleted = DeleteMin(tr.Root)
	if tr.Root != nil {
		tr.Root.Black = true
	}
	return deleted
}

// Delete deletes the node with the Key and returns the Key Interface.
// It returns nil if the Key does not exist in the tree.
//
//
//	Delete Algorithm:
//	1. Start 'delete' from tree Root.
//
//	2. Call 'delete' method recursively on each Node from binary search path.
//		- e.g. If the key to delete is greater than Root's key
//			, call 'delete' on Right Node.
//
//
//	# start
//
//	3. Recursive 'tree.delete(nd, key)'
//
//		if key < nd.Key:
//
//			if nd.Left is empty:
//				# then nothing to delete, so return nil
//				return nd, nil
//
//			if (nd.Left is Black) and (nd.Left.Left is Black):
//				# then move Red from Right to Left to update nd
//				nd = MoveRedFromRightToLeft(nd)
//
//			# recursively call 'delete' to update nd.Left
//			nd.Left, deleted = tr.delete(nd.Left, key)
//
//		else if key >= nd.Key:
//
//			if nd.Left is Red:
//				# RotateToRight(nd) to update nd
//				nd = RotateToRight(nd)
//
//			if (key == nd.Key) and nd.Right is empty:
//				# then return nil, nd.Key to recursively update nd
//				return nil, nd.Key
//
//			if (nd.Right is not empty)
//			and (nd.Right is Black)
//			and (nd.Right.Left is Black):
//				# then move Red from Left to Right to update nd
//				nd = MoveRedFromLeftToRight(nd)
//
//			if (key == nd.Key):
//				# then DeleteMin of nd.Right to update nd.Right
//				nd.Right, subDeleted = DeleteMin(nd.Right)
//
//				# and then update nd.Key with DeleteMin(nd.Right)
//				deleted, nd.Key = nd.Key, subDeleted
//
//			else if key != nd.Key:
//				# recursively call 'delete' to update nd.Right
//				nd.Right, deleted = tr.delete(nd.Right, key)
//
//		# recursively FixUp upwards to Root
//		return FixUp(nd), deleted
//
//	# end
//
//
//	4. If the tree's Root is not nil, set Root Black.
//
//	5. Return the Interface(nil if the key does not exist.)
//
func (tr *Tree) Delete(key Interface) Interface {
	var deleted Interface
	tr.Root, deleted = tr.delete(tr.Root, key)
	if tr.Root != nil {
		tr.Root.Black = true
	}
	return deleted
}

func (tr *Tree) delete(nd *Node, key Interface) (*Node, Interface) {
	if nd == nil {
		return nil, nil
	}

	var deleted Interface

	// if key is Less than nd.Key
	if key.Less(nd.Key) {

		// if key is Less than nd.Key
		// and nd.Left is empty
		if nd.Left == nil {

			// then nothing to delete
			// so return the nil
			return nd, nil
		}

		// if key is Less than nd.Key
		// and nd.Left is Black
		// and nd.Left.Left is Black
		if !isRed(nd.Left) && !isRed(nd.Left.Left) {

			// then MoveRedFromRightToLeft(nd)
			nd = MoveRedFromRightToLeft(nd)
		}

		// and recursively call tr.delete(nd.Left, key)
		nd.Left, deleted = tr.delete(nd.Left, key)

	} else {
		// if key is not Less than nd.Key
		//(or key is greater than or equal to nd.Key)
		//(or key >= nd.Key)

		// and nd.Left is Red
		if isRed(nd.Left) {

			// then RotateToRight(nd)
			nd = RotateToRight(nd)
		}

		// and nd.Key is not Less than key
		// (or nd.Key >= key)
		// (or key == nd.Key)
		// and nd.Right is empty
		if !nd.Key.Less(key) && nd.Right == nil {
			// then return nil to delete the key
			return nil, nd.Key
		}

		// and nd.Right is not empty
		// and nd.Right is Black
		// and nd.Right.Left is Black
		if nd.Right != nil && !isRed(nd.Right) && !isRed(nd.Right.Left) {
			// then MoveRedFromLeftToRight(nd)
			nd = MoveRedFromLeftToRight(nd)
		}

		// and key == nd.Key
		if !nd.Key.Less(key) {

			var subDeleted Interface

			// then DeleteMin(nd.Right)
			nd.Right, subDeleted = DeleteMin(nd.Right)
			if subDeleted == nil {
				panic("Unexpected nil value")
			}

			// and update nd.Key with DeleteMin(nd.Right)
			deleted, nd.Key = nd.Key, subDeleted

		} else {
			// if updated nd.Key is Less than key (nd.Key < key) to update nd.Right
			nd.Right, deleted = tr.delete(nd.Right, key)
		}
	}

	// recursively FixUp upwards to Root
	return FixUp(nd), deleted
}
