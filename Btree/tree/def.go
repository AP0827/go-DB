package tree

const (
	HEADER         = 4
	MAX_PAGESIZE   = 4096
	MAX_KSIZE      = 1000
	MAX_VSIZE      = 3000
	BNODE_LEAF     = 1
	BNODE_INTERNAL = 2
)

type BNode struct {
	data []byte
}

type BTree struct {
	root uint64

	get func(uint64) BNode
	new func(BNode) uint64
	del func(uint64)
}
