package stream

import (
	"crypto/cipher"
	"io"
)

type StreamWriteCloser struct {
	aead cipher.AEAD

	buf []byte

	w         io.WriteCloser // encrypted conn
	seed      []byte
	chunkSize int
}

func NewStreamWriteCloser(seed []byte, chunkSize int, backend io.WriteCloser, aead cipher.AEAD) *StreamWriteCloser {
	w := &StreamWriteCloser{
		aead:      aead,
		buf:       make([]byte, 0), // save plaintext
		seed:      make([]byte, aead.NonceSize()),
		chunkSize: chunkSize,
		w:         backend,
	}

	copy(w.seed, seed)

	return w
}

func (w *StreamWriteCloser) Close() (err error) {
	return w.w.Close()
}

func (w *StreamWriteCloser) write() (n int, err error) {
	chunk := make([]byte, (2+w.chunkSize)+w.aead.Overhead())
	fullChunSizeBytes := packUint16LE(uint16(w.chunkSize))

	var nn int

	for len(w.buf) > 0 {
		if len(w.buf) > w.chunkSize {
			copy(chunk[2:], w.buf[:w.chunkSize])
			w.buf = w.buf[w.chunkSize:]
			copy(chunk[0:2], fullChunSizeBytes)

		} else {
			nn = copy(chunk[2:], w.buf)
			w.buf = w.buf[nn:]
			copy(chunk[0:2], packUint16LE(uint16(nn)))
		}

		w.aead.Seal(chunk[:0], w.seed, chunk[:(2+w.chunkSize)], nil)
		_, err = w.w.Write(chunk)
		if err != nil {
			return
		}
	}

	return
}

func (w *StreamWriteCloser) Write(src []byte) (n int, err error) {
	w.buf = append(w.buf, src...)
	n = len(src)

	_, err = w.write()
	return
}

func (w *StreamWriteCloser) WriteByte(c byte) (err error) {
	w.buf = append(w.buf, c)

	_, err = w.write()
	return
}
