package utils

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type ReadAtWriteAtMock struct {
	b []byte
}

func (m ReadAtWriteAtMock) ReadAt(p []byte, off int64) (int, error) {
	if off >= int64(len(m.b)) {
		return 0, errors.Errorf("out of range")
	}

	return copy(p, m.b[off:]), nil
}

func (m ReadAtWriteAtMock) WriteAt(p []byte, off int64) (int, error) {
	if off >= int64(len(m.b)) {
		return 0, errors.Errorf("out of range")
	}

	return copy(m.b[off:], p), nil
}

func TestOffsetReader(t *testing.T) {
	a := make([]byte, 10)
	b := []byte("0123456789")

	c := NewOffsetReader(ReadAtWriteAtMock{b: b}, 2)

	n, err := c.Read(a[:1])
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
	assert.Equal(t, byte('2'), a[0])
	assert.Equal(t, int64(3), c.Offset())

	n, err = c.Read(a[:7])
	assert.NoError(t, err)
	assert.Equal(t, 7, n)
	assert.Equal(t, b[3:10], a[:7])
	assert.Equal(t, int64(10), c.Offset())

	n, err = c.Read(a[:1])
	assert.Error(t, err)
	assert.Equal(t, 0, n)
	assert.Equal(t, int64(10), c.Offset())
}

func TestOffsetWriter(t *testing.T) {
	a := []byte("0123456789")
	b := make([]byte, 10)

	c := NewOffsetWriter(ReadAtWriteAtMock{b: b}, 2)

	n, err := c.Write(a[:1])
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
	assert.Equal(t, byte('0'), b[2])
	assert.Equal(t, int64(3), c.Offset())

	n, err = c.Write(a[:7])
	assert.NoError(t, err)
	assert.Equal(t, 7, n)
	assert.Equal(t, b[3:10], a[:7])
	assert.Equal(t, int64(10), c.Offset())

	n, err = c.Write(a[:1])
	assert.Error(t, err)
	assert.Equal(t, 0, n)
	assert.Equal(t, int64(10), c.Offset())
}
