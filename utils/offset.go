package utils

import "io"

// OffsetReader translates ReaderAt to Reader given a starting offset
type OffsetReader struct {
	off int64
	r   io.ReaderAt
}

func NewOffsetReader(r io.ReaderAt, off int64) *OffsetReader {
	return &OffsetReader{
		r:   r,
		off: off,
	}
}

// Uses underlying ReaderAt's ReadAt to read data, increment the offset
// by the amount of bytes read and returns any underlying error
func (r *OffsetReader) Read(p []byte) (int, error) {
	n, err := r.r.ReadAt(p, r.off)
	r.off += int64(n)

	return n, err
}

// Returns the current offset
func (r OffsetReader) Offset() int64 {
	return r.off
}

// OffsetWriter translates WriterAt to Writer given a starting offset
type OffsetWriter struct {
	off int64
	w   io.WriterAt
}

func NewOffsetWriter(w io.WriterAt, off int64) *OffsetWriter {
	return &OffsetWriter{
		w:   w,
		off: off,
	}
}

// Uses underlying WriterAt's WriteAt to write data, increment the offset
// by the amount of bytes written and returns any underlying error
func (w *OffsetWriter) Write(p []byte) (int, error) {
	n, err := w.w.WriteAt(p, w.off)
	w.off += int64(n)

	return n, err
}

// Returns the current offset
func (r OffsetWriter) Offset() int64 {
	return r.off
}
