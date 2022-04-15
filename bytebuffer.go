package bytebufferpool

import "io"

// Buffer provides byte buffer, which can be used for minimizing
// memory allocations.
//
// Buffer may be used with functions appending data to the given []byte
// slice. See example code for details.
//
// Use Get for obtaining an empty byte buffer.
type Buffer[T any] struct {

	// B is a byte buffer to use in append-like workloads.
	// See example code for details.
	B []T
}

type ByteBuffer struct {
	*Buffer[byte]
}

// Len returns the size of the byte buffer.
func (b *Buffer[T]) Len() int {
	return len(b.B)
}

// ReadFrom implements io.ReaderFrom.
//
// The function appends all the data read from r to b.
func (b ByteBuffer) ReadFrom(r io.Reader) (int64, error) {
	p := b.B
	nStart := int64(len(p))
	nMax := int64(cap(p))
	n := nStart
	if nMax == 0 {
		nMax = 64
		p = make([]byte, nMax)
	} else {
		p = p[:nMax]
	}
	for {
		if n == nMax {
			nMax *= 2
			bNew := make([]byte, nMax)
			copy(bNew, p)
			p = bNew
		}
		nn, err := r.Read(p[n:])
		n += int64(nn)
		if err != nil {
			b.B = p[:n]
			n -= nStart
			if err == io.EOF {
				return n, nil
			}
			return n, err
		}
	}
}

// WriteTo implements io.WriterTo.
func (b ByteBuffer) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(b.B)
	return int64(n), err
}

// Bytes returns b.B, i.e. all the bytes accumulated in the buffer.
//
// The purpose of this function is bytes.Buffer compatibility.
func (b ByteBuffer) Bytes() []byte {
	return b.B
}

// Write implements io.Writer - it appends p to Buffer.B
func (b *Buffer[T]) Write(p []T) (int, error) {
	b.B = append(b.B, p...)
	return len(p), nil
}

// WriteByte appends the byte c to the buffer.
//
// The purpose of this function is bytes.Buffer compatibility.
//
// The function always returns nil.
func (b ByteBuffer) WriteByte(c byte) error {
	b.B = append(b.B, c)
	return nil
}

// WriteString appends s to Buffer.B.
func (b ByteBuffer) WriteString(s string) (int, error) {
	b.B = append(b.B, s...)
	return len(s), nil
}

// Set sets Buffer.B to p.
func (b *Buffer[T]) Set(p []T) {
	b.B = append(b.B[:0], p...)
}

// SetString sets Buffer.B to s.
func (b ByteBuffer) SetString(s string) {
	b.B = append(b.B[:0], s...)
}

// String returns string representation of Buffer.B.
func (b ByteBuffer) String() string {
	return string(b.B)
}

// Reset makes Buffer.B empty.
func (b *Buffer[T]) Reset() {
	b.B = b.B[:0]
}
