package allocator

import (
	"unsafe"

	"github.com/pkg/errors"
)

// Memory is an allocator that uses go garbage collector for block allocation
type Memory struct {
	allocation map[uintptr][]byte
}

func (m Memory) Alloc(size int) (any, []byte, error) {
	data := make([]byte, size)
	p := uintptr(unsafe.Pointer(&data))
	m.allocation[p] = data

	return p, data, nil
}

func (m Memory) Free(id any) error {
	p, ok := id.(uintptr)
	if !ok {
		return errors.Errorf("id is not uintptr")
	}

	if _, exists := m.allocation[p]; !exists {
		return errors.Errorf("allocation %d does not exist", id)
	}

	delete(m.allocation, p)

	return nil
}

// Creates a new allocator that uses golang memory management
func NewMemory() *Memory {
	return &Memory{
		allocation: map[uintptr][]byte{},
	}
}
