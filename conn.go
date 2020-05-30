package stream

import (
	"crypto/cipher"
	"net"
)

type CryptoConn struct {
	net.Conn
	*StreamWriteCloser
	*StreamReader
}

func NewCryptoConn(seed []byte, chunkSize int, conn net.Conn, aead cipher.AEAD) *CryptoConn {
	return &CryptoConn{
		Conn:              conn,
		StreamWriteCloser: NewStreamWriteCloser(seed, chunkSize, conn, aead),
		StreamReader:      NewStreamReader(seed, chunkSize, conn, aead),
	}
}

func (cc *CryptoConn) Close() {
	cc.StreamWriteCloser.Close()
}

func (cc *CryptoConn) Write(b []byte) (int, error) {
	return cc.StreamWriteCloser.Write(b)
}

func (cc *CryptoConn) Read(b []byte) (int, error) {
	return cc.StreamReader.Read(b)
}
