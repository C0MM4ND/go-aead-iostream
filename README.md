# go-aead-iostream

IO stream for [go std AEAD](https://godoc.org/crypto/cipher#AEAD)

## Status

[![Build Status](https://travis-ci.com/maoxs2/go-aead-iostream.svg?branch=master)](https://travis-ci.com/maoxs2/go-aead-iostream)

## Conn example

https://github.com/maoxs2/go-aead-compress-conn

## Example

```go
import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"io"
	"os"
	"log"

	stream "github.com/maoxs2/go-aead-iostream"
)

func main() {
    seed := hash([]byte("Hello"))

	c1, err := aes.NewCipher(seed)
	if err != nil {
		panic(err)
	}
	c2, err := aes.NewCipher(seed)
	if err != nil {
		panic(err)
	}

	aead1, err := cipher.NewGCM(c1)
	if err != nil {
		panic(err)
	}
	aead2, err := cipher.NewGCM(c2)
	if err != nil {
		panic(err)
	}

	f1, err := os.OpenFile("test", os.O_CREATE|os.O_WRONLY, 644)
	if err != nil {
		panic(err)
	}

	f2, err := os.Open("test")
	if err != nil {
		panic(err)
	}

	var chunkSize = 64

	w := stream.NewStreamWriteCloser(seed, chunkSize, f1, aead2)

	rawMessage := []byte("Package cipher implements standard block cipher modes that can be wrapped around low-level block cipher implementations. See https://csrc.nist.gov/groups/ST/toolkit/BCM/current_modes.html and NIST Special Publication 800-38A.")

	w.Write(rawMessage)
	w.Close()

	r := stream.NewStreamReader(seed, chunkSize, f2, aead1)

	buf := make([]byte, 2048)
	dst := make([]byte, 0)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			dst = append(dst, buf[:n]...)
		}
		if err != nil && err != io.EOF {
			panic(err)
		}
		if err == io.EOF {
			break
		}

	}

	if !bytes.Equal(dst, rawMessage) {
		log.Error("dst is %s, but raw is %s", dst, rawMessage)
	} else {
		log.Println("pass")
	}

	f2.Close()
	os.Remove("test")
}


```

