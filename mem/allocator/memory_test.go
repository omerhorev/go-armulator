package allocator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	// "github.com/omerhorev/goarmulator/mem/allocator"
)

func TestMemoryAllocator(t *testing.T) {
	m := NewMemory()

	i1, _, err := m.Alloc(10)

	assert.NoError(t, err)
	assert.Contains(t, m.allocation, i1)

	i2, _, err := m.Alloc(20)
	assert.NoError(t, err)
	assert.Contains(t, m.allocation, i1)
	assert.Contains(t, m.allocation, i2)

	assert.NoError(t, m.Free(i1))
	assert.NotContains(t, m.allocation, i1)
	assert.Contains(t, m.allocation, i2)

	assert.NoError(t, m.Free(i2))
	assert.NotContains(t, m.allocation, i1)
	assert.NotContains(t, m.allocation, i2)

	assert.Error(t, m.Free(uintptr(0)))
	assert.Error(t, m.Free("a"))

	m.Alloc(20)
	assert.NoError(t, m.Close())

}
