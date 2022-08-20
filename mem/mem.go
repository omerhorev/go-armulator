package mem

import (
	"io"

	"github.com/pkg/errors"

	"github.com/omerhorev/goarmulator/mem/allocator"
)

type allocBlock struct {
	allocId any
	addr    uint64
	data    []byte
	perm    Permissions
}

// Mem is a memory management system that acts as a MMU. Mem translates
// virtual to physical addresses and provide an interface to read and write
// to that virtual memory space
type Mem struct {
	allocator allocator.Allocator
	blocks    map[uint64]*allocBlock
	raw       *MemRW
	r         *MemRW // Access the virtual memory with read permissions
	rx        *MemRW // Access the virtual memory with execute permissions
	w         *MemRW // Access the virtual memory with write permissions
}

// MemRW is the gateway to all the memory operations with read/execute permissions
// it implements ReadAt
type MemRW struct {
	perm Permissions
	mem  *Mem
}

// Creates a new memory management
func NewMem(allocator allocator.Allocator) *Mem {
	m := &Mem{
		allocator: allocator,
		blocks:    map[uint64]*allocBlock{},
		r:         &MemRW{perm: PermRead},
		rx:        &MemRW{perm: PermExecute},
		w:         &MemRW{perm: PermWrite},
		raw:       &MemRW{perm: PermRead | PermWrite | PermExecute},
	}

	m.r.mem = m
	m.rx.mem = m
	m.w.mem = m
	m.raw.mem = m

	return m
}

// Creates a new memory management that uses MemoryAllocator
func NewMemFromMemory() *Mem {
	return NewMem(allocator.NewMemory())
}

// Allocates a new memory block and map it to a virtual address
// This method uses the underlying Allocator to allocated memory blocks
func (m *Mem) Alloc(addr uint64, size int, permissions Permissions) error {
	id, data, err := m.allocator.Alloc(size)
	if err != nil {
		return errors.Wrap(err, "alloc")
	}

	m.blocks[addr] = &allocBlock{
		allocId: id,
		data:    data,
		addr:    addr,
		perm:    permissions,
	}

	return nil
}

// Free a memory block and from both the underlying Allocator and
// the memory manager.
// address must be the address used in the Alloc method.
func (m *Mem) Free(address uint64) error {
	b, exists := m.blocks[address]
	if !exists {
		return errors.Errorf("no block was allocated at address 0x%x", address)
	}

	if err := m.allocator.Free(b.allocId); err != nil {
		return err
	}

	delete(m.blocks, address)

	return nil
}

// Writes to a specific virtual address.
// Returns SegmentationFaultError if trying to access unallocated or
// memory without permission
func (m *Mem) WriteAt(p []byte, addr int64) (int, error) {
	return m.w.WriteAt(p, addr)
}

// Reads from a specific virtual address with read permissions.
// Returns SegmentationFaultError if trying to access unallocated or
// memory without permission
func (m *Mem) ReadAt(p []byte, addr int64) (int, error) {
	return m.r.ReadAt(p, addr)
}

func (m *Mem) getAllocBlock(addr uint64, permissions Permissions) *allocBlock {
	var block *allocBlock = nil

	for k, v := range m.blocks {
		if uint64(addr) >= k && uint64(addr) < k+uint64(len(v.data)) {
			block = v
		}
	}

	if block == nil || !block.perm.Has(permissions) {
		return nil
	}

	return block
}

// Read operations using read permission
func (m *Mem) Reader() io.ReaderAt {
	return m.r
}

// Read operations using execute permission
func (m *Mem) ReaderX() io.ReaderAt {
	return m.rx
}

// Read operations using execute permission
func (m *Mem) Writer() io.WriterAt {
	return m.w
}

// Read operations using execute permission
func (m *Mem) Raw() *MemRW {
	return m.raw
}

// Closes the memory manager, thus freeing all allocated resources
func (m *Mem) Close() error {
	for _, b := range m.blocks {
		m.Free(b.addr)
	}

	return nil
}

// Reads from a specific virtual address.
// Returns SegmentationFaultError if trying to access unallocated or
// memory without permission
func (m *MemRW) ReadAt(p []byte, addr int64) (int, error) {
	n := 0
	for n < len(p) {
		b := m.mem.getAllocBlock(uint64(addr), m.perm)

		if b == nil {
			return 0, NewSegmentationFaultError(addr, m.perm)
		}

		off := uint64(addr) - b.addr

		n += copy(p[n:], b.data[off:])
		addr += int64(n)
	}

	return n, nil
}

// Writes to a specific virtual address.
// Returns SegmentationFaultError if trying to access unallocated or
// memory without permission
func (m *MemRW) WriteAt(p []byte, addr int64) (int, error) {
	n := 0
	for n < len(p) {
		b := m.mem.getAllocBlock(uint64(addr), m.perm)

		if b == nil {
			return 0, NewSegmentationFaultError(addr, m.perm)
		}

		off := uint64(addr) - b.addr

		n += copy(b.data[off:], p[n:])
		addr += int64(n)
	}

	return n, nil
}
