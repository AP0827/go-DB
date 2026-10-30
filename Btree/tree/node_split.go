package tree

import (
	"BTree/utils"
)

// Split an oversized node into two nodes.
// The right node always fits within one page.
func nodeSplit2(left BNode, right BNode, old BNode) {
	utils.Assert(old.nkeys() >= 2)

	// Initial guess: split the keys approximately in half.
	nleft := old.nkeys() / 2

	// Calculate the size of the left node.
	leftBytes := func() uint16 {
		return HEADER +
			8*nleft +
			2*nleft +
			old.getOffset(nleft)
	}

	// Move the split point left until the left node fits.
	for leftBytes() > MAX_PAGESIZE {
		nleft--
	}
	utils.Assert(nleft >= 1)

	// Calculate the size of the right node.
	rightBytes := func() uint16 {
		return old.nbytes() - leftBytes() + HEADER
	}

	// Move the split point right until the right node fits.
	for rightBytes() > MAX_PAGESIZE {
		nleft++
	}
	utils.Assert(nleft < old.nkeys())

	nright := old.nkeys() - nleft

	// Initialize the destination nodes.
	left.set(old.btype(), nleft)
	right.set(old.btype(), nright)

	// Copy the entries into the two nodes.
	nodeAppendRange(left, old, 0, 0, nleft)
	nodeAppendRange(right, old, 0, nleft, nright)

	// The right node must fit within a page.
	utils.Assert(right.nbytes() <= MAX_PAGESIZE)
}
