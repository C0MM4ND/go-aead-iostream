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
	if len(r.buf) > 0 {
		n := copy(dst, r.buf)
		r.buf = r.buf[n:]
		return n, nil
	}

	nn, err := r.read()
	n = copy(dst, r.buf[:nn])
	r.buf = r.buf[n:]

	return
}

func (r *StreamReader) read() (n int, err error) {
	chunk := make([]byte, (2+r.chunkSize)+r.aead.Overhead())
	fullChunSizeBytes := packUint16LE(uint16(r.chunkSize))

	var l int

	l, err = io.ReadFull(r.r, chunk)

	if err != nil {
		return 0, err
	}

	if l > 0 {
		_, err = r.aead.Open(chunk[:0], r.seed, chunk, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		if bytes.Equal(chunk[0:2], fullChunSizeBytes) {
			n += r.chunkSize
			r.buf = append(r.buf, chunk[2:2+r.chunkSize]...)
		} else {
			ll := binary.LittleEndian.Uint16(chunk[0:2])
			n += int(ll)
			r.buf = append(r.buf, chunk[2:2+ll]...)
		}
	}

	return
}
