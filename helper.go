package stream

import "encoding/binary"

func packUint16LE(n uint16) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, n)
	return b
}
