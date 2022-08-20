package allocator

// Allocator is an interface that wraps the basic allocation interface
//
// The memory allocated using this method will act as a physical memory
// and will be used by other logic with a page managenet system
type Allocator interface {
	// Allocate a new memory block
	//
	// Returns the block id (any), the block data and a possible error
	// The block id will be used to free the block
	Alloc(size int) (any, []byte, error)

	// Free a memory block by an id
	Free(id any) error
}
