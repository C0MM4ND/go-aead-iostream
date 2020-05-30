package aeadconn

import (
	"crypto/cipher"
	"net"

	stream "github.com/maoxs2/go-aead-iostream"
)

type AEADConn struct {
	net.Conn
	*stream.StreamWriteCloser
	*stream.StreamReader
}

func NewAEADConn(seed []byte, chunkSize int, conn net.Conn, aead cipher.AEAD) *AEADConn {
	return &AEADConn{
		Conn:              conn,
		StreamWriteCloser: stream.NewStreamWriteCloser(seed, chunkSize, conn, aead),
		StreamReader:      stream.NewStreamReader(seed, chunkSize, conn, aead),
	}
}

func (cc *AEADConn) Close() error {
	return cc.StreamWriteCloser.Close()
}

func (cc *AEADConn) Write(b []byte) (int, error) {
	return cc.StreamWriteCloser.Write(b)
}

func (cc *AEADConn) Read(b []byte) (int, error) {
	return cc.StreamReader.Read(b)
}
