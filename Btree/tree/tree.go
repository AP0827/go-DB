package tree

import (
	"encoding/binary"
	"fmt"

	"BTree/utils"
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
