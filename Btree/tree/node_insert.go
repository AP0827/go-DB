package tree

import (
	"BTree/utils"
	"bytes"
	"encoding/binary"
)

//B-Tree Insertion

// 1. Find the appropiate insertion point
func nodeLookupLE(node BNode, key []byte) uint16 {
	found := uint16(0) // Store the answer in found and return found
	nkeys := node.nkeys()

	// the 1st node (root, 0 index) is always less than or equal to the key
	for i := uint16(1); i < nkeys; i++ {
		cmp := bytes.Compare(key, node.getKey(i))

		if cmp <= 0 {
			found = i
		}
		if cmp >= 0 {
			break
		}
	}
	return found
}

//2. Leaf insert
// After finding the relevant node position, we create a copy in which we insert the K-V at the appropriate position

func InsertLeaf(new, old BNode, idx uint16, key, value []byte) {
	new.set(BNODE_LEAF, old.nkeys()+1)
	//we'll get idx from the nodeLookupLE function
	nodeAppendRange(new, old, 0, 0, idx)
	nodeAppendKV(new, idx, 0, key, value)
	nodeAppendRange(new, old, idx+1, idx, old.nkeys()-idx)
}
func UpdateLeaf(new, old BNode, idx uint16, key, value []byte) {
	new.set(BNODE_LEAF, old.nkeys())

	//we'll get idx from the nodeLookupLE function
	nodeAppendRange(new, old, 0, 0, idx-1)
	nodeAppendKV(new, idx, 0, key, value)
	nodeAppendRange(new, old, idx+1, idx+1, old.nkeys()-idx-1)

}

func nodeAppendKV(node BNode, idx uint16, ptr uint64, key, val []byte) {
	node.setPtr(idx, ptr)

	pos := node.kvPos(idx)

	// Append the lenght of the key and value to the node
	binary.LittleEndian.PutUint16(node.data[pos:], uint16(len(key)))
	binary.LittleEndian.PutUint16(node.data[pos+2:], uint16(len(val)))

	// Use copy function of golang
	copy(node.data[pos+4:], key)
	copy(node.data[pos+4+uint16(len(key)):], val)

	// Now set the offset
	node.setOffset(idx+1, node.getOffset(idx)+4+uint16(len(key)+len(val)))
}

func nodeAppendRange(new, old BNode, startIdx_new, startIdx_old, count uint16) {

	//Pointers
	utils.Assert(startIdx_new <= new.nkeys())
	utils.Assert(startIdx_old <= old.nkeys())

	for i := uint16(0); i < count; i++ {
		new.setPtr(startIdx_new+i, old.getPtr(startIdx_old+i))
	}

	// Offsets
	newBegin := new.getOffset(startIdx_new)
	oldBegin := old.getOffset(startIdx_old)

	for i := uint16(1); i < count; i++ {
		offset := newBegin + old.getOffset(startIdx_old+i) - oldBegin
		new.setOffset(startIdx_new+i, offset)
	}

	//KVs
	//Simply dump the data

	oldstartKeyVal := old.kvPos(startIdx_old)
	newstartKeyVal := old.kvPos(startIdx_old)

	copy(new.data[newstartKeyVal:], old.data[oldstartKeyVal:])
}

//3. Recursive Insertion

// insert a KV into a node, the result might be split into 2 nodes.
// the caller is responsible for deallocating the input node
// and splitting and allocating result nodes.

func treeInsert(tree *BTree, node BNode, key, val []byte) BNode {
	new := BNode{data: make([]byte, 2*MAX_PAGESIZE)}

	idx := nodeLookupLE(node, key)

	switch node.btype() {
	case BNODE_LEAF:
		if bytes.Equal(node.getKey(idx), key) {
			UpdateLeaf(new, node, idx, key, val)
		} else {
			InsertLeaf(new, node, idx+1, key, val)
		}
	case BNODE_INTERNAL:
		nodeInsert(tree, new, node, idx, key, val)
	default:
		panic("Bad Node!")
	}

	return new
}

func nodeInsert(tree *BTree, new, node BNode, idx uint16, key, val []byte) {
	kptr := node.getPtr(idx)
	knode := tree.get(kptr)
	tree.del(kptr)

	//recursive insertion to kid node
	knode = treeInsert(tree, knode, key, val)

	nsplit, splited := nodeSplit3(knode)

	//update the kid links
	nodeReplaceKidN(tree, new, node, idx, splited[:nsplit])

}
