package handler

import "sync"

type bytesBuf []byte

var bufPool = sync.Pool{
	New: func() any {
		// We rarely see log lines longer than 300 bytes.
		b := make([]byte, 0, 512)
		return (*bytesBuf)(&b)
	},
}

func newBuf() *bytesBuf {
	return bufPool.Get().(*bytesBuf)
}

func (b *bytesBuf) free() {
	// Don't return buffers larger than 1K to the pool.
	if cap(*b) <= 1024 {
		*b = (*b)[:0]
		bufPool.Put(b)
	}
}

func (b *bytesBuf) write(data []byte) {
	*b = append(*b, data...)
}

func (b *bytesBuf) Write(data []byte) (int, error) {
	*b = append(*b, data...)
	return len(data), nil
}

func (b *bytesBuf) writeString(s string) {
	*b = append(*b, s...)
}

func (b *bytesBuf) writeByte(c byte) {
	*b = append(*b, c)
}
