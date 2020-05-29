package stream

import (
	"bytes"
	"crypto/cipher"
	"encoding/binary"
	"fmt"
	"io"
)

type StreamReader struct {
	aead cipher.AEAD

	buf []byte

	r         io.Reader // encrypted conn
	seed      []byte
	chunkSize int
}

func NewStreamReader(seed []byte, chunkSize int, backend io.Reader, aead cipher.AEAD) *StreamReader {
	r := &StreamReader{
		aead:      aead,
		buf:       make([]byte, 0), // save plaintext
		seed:      make([]byte, aead.NonceSize()),
		chunkSize: chunkSize,
		r:         backend,
	}

	copy(r.seed, seed)

	return r
}

func (r *StreamReader) Read(dst []byte) (n int, err error) {
	chunk := make([]byte, (2+r.chunkSize)+r.aead.Overhead())
	fullChunSizeBytes := packUint16LE(uint16(r.chunkSize))

	var l int

	for {
		l, err = r.r.Read(chunk)

		if err == io.ErrUnexpectedEOF {
			chunk = chunk[:l]
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return
		}

		if l > 0 {
			_, err = r.aead.Open(chunk[:0], r.seed, chunk, nil)
			if err != nil {
				fmt.Println(err)
				return
			}

			if bytes.Equal(chunk[0:2], fullChunSizeBytes) {
				r.buf = append(r.buf, chunk[2:2+r.chunkSize]...)
			} else {
				r.buf = append(r.buf, chunk[2:2+binary.LittleEndian.Uint16(chunk[0:2])]...)
			}

		}
	}

	n = copy(dst, r.buf)
	r.buf = r.buf[n:]

	return
}
