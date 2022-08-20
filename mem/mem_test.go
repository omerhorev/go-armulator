package mem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testingBuffer = make([]byte, 1000)
	data10bytes1  = []byte("0123456789")
	data10bytes2  = []byte("abcdefgijk")
	data10bytes3  = []byte("kdjasdyf6s")
)

func TestMemRead(t *testing.T) {
	m := NewMemFromMemory()
	m.Alloc(10, 10, PermReadWrite)
	m.Alloc(30, 10, PermReadWrite)
	m.Alloc(40, 10, PermReadWrite)
	m.WriteAt(data10bytes1, 10)
	m.WriteAt(data10bytes2, 30)
	m.WriteAt(data10bytes3, 40)

	// whole block
	n, err := m.ReadAt(testingBuffer[:10], 10)
	assert.NoError(t, err)
	assert.Equal(t, 10, n)
	cleanTestingBuffer()

	// start of block
	n, err = m.ReadAt(testingBuffer[:3], 10)
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, data10bytes1[:3], testingBuffer[:3])
	cleanTestingBuffer()

	// end of block
	n, err = m.ReadAt(testingBuffer[:3], 17)
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, data10bytes1[7:10], testingBuffer[:3])
	cleanTestingBuffer()

	// middle of block
	n, err = m.ReadAt(testingBuffer[:3], 16)
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, data10bytes1[6:9], testingBuffer[:3])
	cleanTestingBuffer()

	// read more that a block, with unallocated parts
	_, err = m.ReadAt(testingBuffer[:11], 10)
	assert.Error(t, err)
	cleanTestingBuffer()

	// read more that a block, with allocated parts
	_, err = m.ReadAt(testingBuffer[:20], 30)
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, data10bytes2, testingBuffer[00:10])
	assert.Equal(t, data10bytes3, testingBuffer[10:20])
	cleanTestingBuffer()
}

func TestErrors(t *testing.T) {
	m := NewMemFromMemory()
	m.Alloc(10, 10, PermRead)

	_, err := m.Reader().ReadAt(testingBuffer[:1], 9)
	assert.ErrorContains(t, err, "r--")
	assert.ErrorContains(t, err, "9")

	_, err = m.ReaderX().ReadAt(testingBuffer[:1], 26)
	assert.ErrorContains(t, err, "--x")
	assert.ErrorContains(t, err, "0x1a")

	_, err = m.Writer().WriteAt(testingBuffer[:1], 9)
	assert.ErrorContains(t, err, "-w-")
	assert.ErrorContains(t, err, "9")
}

func TestPermissions(t *testing.T) {
	m := NewMemFromMemory()
	m.Alloc(10, 10, PermRead)

	_, err := m.Reader().ReadAt(testingBuffer[:1], 10)
	assert.NoError(t, err)

	_, err = m.Writer().WriteAt(testingBuffer[:1], 10)
	assert.Error(t, err)

	_, err = m.ReaderX().ReadAt(testingBuffer[:1], 10)
	assert.Error(t, err)
}

func TestFree(t *testing.T) {
	m := NewMemFromMemory()
	m.Alloc(10, 10, PermRead)
	m.Alloc(20, 10, PermRead)

	assert.Error(t, m.Free(11))
	assert.Contains(t, m.blocks, uint64(10))
	assert.Contains(t, m.blocks, uint64(20))

	assert.NoError(t, m.Free(10))
	assert.NotContains(t, m.blocks, 10)
	assert.Contains(t, m.blocks, uint64(20))

	assert.NoError(t, m.Close())
	assert.NotContains(t, m.blocks, 10)
	assert.NotContains(t, m.blocks, 20)
}

func cleanTestingBuffer() {
	for i := range testingBuffer {
		testingBuffer[i] = 0
	}
}
