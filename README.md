# go-aead-iostream

This is a fork of the [implementation of c0mm4nd](https://github.com/c0mm4nd/go-aead-iostream).

IO stream for [go std AEAD](https://godoc.org/crypto/cipher#AEAD).

## Conn example

https://github.com/c0mm4nd/go-aead-conn

## Example

```go
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
	"os"
	"log"

	stream "github.com/stweiz/go-aead-iostream"
)

func main() {
	// replace with your key
	seed := []byte("an arbitrary key")
	pathToNewFile := "/tmp/initial_file"
	pathToEncryptedFile := "/tmp/encrypted_file"
	pathToDecryptedFile := "/tmp/decrypted_file"

	err := encryptNewFile(seed, pathToNewFile, pathToEncryptedFile)
	if err != nil {
		panic(err)
	}

	err = decryptExistingFile(seed, pathToDecryptedFile, pathToEncryptedFile)
	if err != nil {
		panic(err)
	}
}

func encryptNewFile(seed []byte, pathToNewFile string, pathToEncryptedFile string) error {
	// Prepare ciphers.
	aesCipher, err := aes.NewCipher(seed)
	if err != nil {
		log.Print("Couldn't create AES cipher.")
		log.Printf("Error: %s", err)
		return err
	}
	aeadCipher, err := cipher.NewGCM(aesCipher)
	if err != nil {
		log.Print("Couldn't create AEAD cipher.")
		log.Printf("Error: %s", err)
		return err
	}

	// Create a new file and write some content into it.
	initialFile, err := os.OpenFile(pathToNewFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Couldn't create or open file. Path: %s", pathToNewFile)
		log.Printf("Error: %s", err)
		return err
	}
	defer initialFile.Close()
	rawMessage := []byte("Package cipher implements standard block cipher modes that can be wrapped around low-level block cipher implementations. See https://csrc.nist.gov/groups/ST/toolkit/BCM/current_modes.html and NIST Special Publication 800-38A.")
	initialFile.Write(rawMessage)
	initialFile.Close()

	// Reopen the new file for reading.
	initialFile, err = os.OpenFile(pathToNewFile, os.O_RDONLY, 0644)
	if err != nil {
		log.Printf("Couldn't open file for reading. Path: %s", pathToNewFile)
		log.Printf("Error: %s", err)
		return err
	}
	defer initialFile.Close()

	// Create an empty file, to which can be written later.
	encryptedFile, err := os.OpenFile(pathToEncryptedFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Couldn't create or open file. Path: %s", pathToEncryptedFile)
		log.Printf("Error: %s", err)
		return err
	}
	defer encryptedFile.Close()

	// Create the StreamWriteCloser, which can be piped into any output.
	chunkSize := 64
	encryptedWriter := stream.NewStreamWriteCloser(seed, chunkSize, encryptedFile, aeadCipher)
	defer encryptedWriter.Close()

	// Create a buffer to hold the data in the specified chunk size.
	buf := make([]byte, chunkSize)

	// Use io.CopyBuffer to read from the unencrypted file stream and write to the encrypted stream.
	if _, err := io.CopyBuffer(encryptedWriter, initialFile, buf); err != nil {
		log.Printf("Couldn't write encrypted file. Path: %s", pathToEncryptedFile)
		log.Printf("Error: %s", err)
		return err
	}

	return nil
}

func decryptExistingFile(seed []byte, pathToDecryptedFile string, pathToEncryptedFile string) error {
	// Prepare ciphers.
	block, err := aes.NewCipher(seed)
	if err != nil {
		log.Print("Couldn't create AES cipher.")
		log.Printf("Error: %s", err)
		return err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		log.Print("Couldn't create AEAD cipher.")
		log.Printf("Error: %s", err)
		return err
	}

	// Open the encrypted file for reading.
	encryptedFile, err := os.Open(pathToEncryptedFile)
	if err != nil {
		log.Printf("Couldn't create or open file. Path: %s", pathToEncryptedFile)
		log.Printf("Error: %s", err)
		return err
	}
	defer encryptedFile.Close()

	// Create a new file to write the decrypted data to.
	decryptedFile, err := os.Create(pathToDecryptedFile)
	if err != nil {
		log.Printf("Couldn't create or open file. Path: %s", pathToDecryptedFile)
		log.Printf("Error: %s", err)
		return err
	}
	defer decryptedFile.Close()

	// Create the StreamReader, which can be piped into any output.
	chunkSize := 64
	r := stream.NewStreamReader(seed, chunkSize, encryptedFile, aead)

	// Create a buffer to hold the data in the specified chunk size.
	buf := make([]byte, 64)

	// Use io.CopyBuffer to read from the encrypted file stream and write to the decrypted stream.
	if _, err := io.CopyBuffer(decryptedFile, r, buf); err != nil {
		log.Printf("Couldn't write encrypted file. Path: %s", pathToDecryptedFile)
		log.Printf("Error: %s", err)
		return err
	}

	return nil
}
```
