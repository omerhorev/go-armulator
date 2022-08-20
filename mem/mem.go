package mem

import (
	"github.com/pkg/errors"

	"github.com/omerhorev/goarmulator/mem/allocator"
)

type allocBlock struct {
	allocId any
	addr    int64
	data    []byte
	perm    Permissions
}

// Mem is a memory management system that acts as a MMU. Mem translates
// virtual to physical addresses and provide an interface to read and write
// to that virtual memory space
type Mem struct {
	allocator     allocator.Allocator
	blocks        map[int64]*allocBlock
	Reader        *MemReader // Access the virtual memory with read permissions
	ReaderExecute *MemReader // Access the virtual memory with execute permissions
	Writer        *MemWriter // Access the virtual memory with write permissions
}

// MemReader is the gateway to all the memory operations with read/execute permissions
// it implements ReadAt
type MemReader struct {
	perm Permissions
	mem  *Mem
}

// MemWriter is the gateway to all the memory operations with write permissions
// it implements WriteAt
type MemWriter struct {
	mem *Mem
}

// Creates a new memory management
func NewMem(allocator allocator.Allocator) *Mem {
	m := &Mem{
		allocator:     allocator,
		blocks:        map[int64]*allocBlock{},
		Reader:        &MemReader{perm: PermRead},
		ReaderExecute: &MemReader{perm: PermExecute},
		Writer:        &MemWriter{},
	}

	m.Reader.mem = m
	m.ReaderExecute.mem = m
	m.Writer.mem = m

	return m
}

// Creates a new memory management that uses MemoryAllocator
func NewMemFromMemory() *Mem {
	return NewMem(allocator.NewMemory())
}

// Allocates a new memory block and map it to a virtual address
// This method uses the underlying Allocator to allocated memory blocks
func (m *Mem) Alloc(addr int64, size int, permissions Permissions) error {
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
func (m *Mem) Free(address int64) error {
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
	return m.Writer.WriteAt(p, addr)
}

// Writes to a specific virtual address.
// Returns SegmentationFaultError if trying to access unallocated or
// memory without permission
func (m *MemWriter) WriteAt(p []byte, addr int64) (int, error) {
	n := 0
	for n < len(p) {
		b := m.mem.getAllocBlock(addr, PermWrite)

		if b == nil {
			return 0, NewSegmentationFaultError(addr, PermWrite)
		}

		off := addr - b.addr

		n += copy(b.data[off:], p[n:])
		addr += int64(n)
	}

	return n, nil
}

// Reads from a specific virtual address with read permissions.
// Returns SegmentationFaultError if trying to access unallocated or
// memory without permission
func (m *Mem) ReadAt(p []byte, addr int64) (int, error) {
	return m.Reader.ReadAt(p, addr)
}

// Reads from a specific virtual address.
// Returns SegmentationFaultError if trying to access unallocated or
// memory without permission
func (m *MemReader) ReadAt(p []byte, addr int64) (int, error) {
	n := 0
	for n < len(p) {
		b := m.mem.getAllocBlock(addr, m.perm)

		if b == nil {
			return 0, NewSegmentationFaultError(addr, m.perm)
		}

		off := addr - b.addr

		n += copy(p[n:], b.data[off:])
		addr += int64(n)
	}

	return n, nil
}

func (m *Mem) getAllocBlock(addr int64, permissions Permissions) *allocBlock {
	var block *allocBlock = nil

	for k, v := range m.blocks {
		if addr >= k && addr < k+int64(len(v.data)) {
			block = v
		}
	}

	if block == nil || !block.perm.Has(permissions) {
		return nil
	}

	return block
}
