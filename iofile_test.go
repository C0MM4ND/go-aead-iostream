package stream_test

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"io"
	"io/ioutil"
	"os"
	"testing"

	stream "github.com/stweiz/go-aead-iostream"
)

func TestAEADFileStream(t *testing.T) {
	seed := hash([]byte("Hello"))
	rawMessage := []byte("Package cipher implements standard block cipher modes that can be wrapped around low-level block cipher implementations. See https://csrc.nist.gov/groups/ST/toolkit/BCM/current_modes.html and NIST Special Publication 800-38A.")
	chunkSize := 64

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

	// WRITE
	f1, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		panic(err)
	}

	w := stream.NewStreamWriteCloser(seed, chunkSize, f1, aead2)
	w.Write(rawMessage)
	w.Close()

	// READ
	f2, err := os.Open(f1.Name())
	if err != nil {
		panic(err)
	}

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
		t.Errorf("dst is %s, but raw is %s", dst, rawMessage)
	} else {
		t.Log("pass")
	}

	f2.Close()
}

func hash(b []byte) []byte {
	hash := sha256.Sum256(b)
	return hash[:]
}
