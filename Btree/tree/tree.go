package tree

import (
	"BTree/utils"
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	HEADER         = 4
	MAX_PAGESIZE   = 4096
	MAX_KSIZE      = 1000
	MAX_VSIZE      = 3000
	BNODE_LEAF     = 1
	BNODE_INTERNAL = 2
)

// | type | nkeys | pointers   | offsets    | key-values
// | 2B   | 2B    | nkeys * 8B | nkeys * 2B | ...
type BNode struct {
	data []byte
}

type BTree struct {
	root uint64

	get func(uint64) BNode
	new func(BNode) uint64
	del func(uint64)
}

func init() {
	size := HEADER + 8 + 2 + 4 + MAX_KSIZE + MAX_VSIZE
	utils.Assert(size <= MAX_PAGESIZE)
	if size > MAX_PAGESIZE {
		fmt.Printf("B-Tree Node size (%d) exceeds the allowed page size (%d)", size, MAX_PAGESIZE)
		return
	}
}

func (node BNode) btype() uint16 {
	return binary.LittleEndian.Uint16(node.data)
}
func (node BNode) nkeys() uint16 {
	return binary.LittleEndian.Uint16(node.data[2:4])
}
func (node BNode) set(btype, nkeys uint16) {
	binary.LittleEndian.PutUint16(node.data[0:2], btype)
	binary.LittleEndian.PutUint16(node.data[2:4], nkeys)
}

func (node BNode) getPtr(idx uint16) uint64 {
	utils.Assert(idx < node.nkeys())
	pos := HEADER + 8*idx
	return binary.LittleEndian.Uint64(node.data[pos:])
}
func (node BNode) setPtr(idx uint16, val uint64) {
	utils.Assert(idx < node.nkeys())
	pos := HEADER + 8*idx
	binary.LittleEndian.PutUint64(node.data[pos:], val)
}

func offsetPos(node BNode, idx uint16) uint16 {
	utils.Assert(1 <= idx && idx <= node.nkeys())
	return HEADER + 8*node.nkeys() + 2*(idx-1)
}
func (node BNode) getOffset(idx uint16) uint16 {
	if idx == 0 {
		return 0
	}
	return binary.LittleEndian.Uint16(node.data[offsetPos(node, idx):])
}

func (node BNode) setOffset(idx, offset uint16) {
	binary.LittleEndian.PutUint16(node.data[offsetPos(node, idx):], offset)
}

func (node BNode) kvPos(idx uint16) uint16 {
	utils.Assert(idx <= node.nkeys())
	return HEADER + 8*node.nkeys() + 2*node.nkeys() + node.getOffset(idx)
}

func (node BNode) getKey(idx uint16) []byte {
	pos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node.data[pos:])

	return node.data[pos+4:][:klen]
}

func (node BNode) getVal(idx uint16) []byte {
	utils.Assert(idx < node.nkeys())
	pos := node.kvPos(idx)
	klen := binary.LittleEndian.Uint16(node.data[pos+0:])
	vlen := binary.LittleEndian.Uint16(node.data[pos+2:])
	return node.data[pos+4+klen:][:vlen]
}

// node size in bytes
func (node BNode) nbytes() uint16 {
	return node.kvPos(node.nkeys()) //simply get the value of last node's index using node.nkeys()
}

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
