package socket

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

const err_overflow_cap = "Overflow capacity"

type BufferReader struct {
	buf *bytes.Buffer
	src io.Reader
	lim int64
}

func NewBufferReader(lim int64, src io.Reader) *BufferReader {
	return &BufferReader{buf: bytes.NewBuffer(make([]byte, 0, lim)), src: src, lim: lim}
}

func (b *BufferReader) Next(sep byte) ([]byte, error) {
	var totRead = int64(b.buf.Len())
	for totRead < b.lim {
		if b.buf.Len() != 0 {
			for _, buf := range b.buf.Bytes() {
				if buf == sep {
					return b.buf.ReadBytes(sep)
				}
			}
		}

		n, err := b.readNext()
		if err != nil {
			return nil, err
		}
		totRead += n
	}
	return nil, errors.New(err_overflow_cap)
}

func (b *BufferReader) readNext() (int64, error) {
	n, err := io.Copy(b.buf, b.src)
	if err != nil {
		return -1, err
	}
	return n, nil
}

func ReadWhile(reader io.Reader, sep byte, buffer *bytes.Buffer) (int, error) {
	_, err := io.Copy(buffer, reader)
	if err != nil {
		return -1, err
	}
	for idx, b := range buffer.Bytes() {
		if b == sep {
			return idx, nil
		}
	}
	return -1, fmt.Errorf("No separator found")
}
