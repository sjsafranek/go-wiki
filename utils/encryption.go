package utils

import (
	"crypto/md5"
	"crypto/sha512"
	"fmt"
	"hash"
	"io"
	"os"
)

// Md5FromBytes generates mb5 from bytea
func Md5FromBytes(data []byte) string {
	h := md5.New()
	return hashByteA(data, h)
}

// Md5FromFile generates mb5 hash from file contents
func Md5FromFile(filename string) (string, error) {
	h := md5.New()
	return hashFile(filename, h)
}

// SHA512FromBytes generates sha512 hash from bytea
func SHA512FromBytes(data []byte) string {
	h := sha512.New()
	return hashByteA(data, h)
}

// SHA512FromFile generates sha512 file contents
func SHA512FromFile(filename string) (string, error) {
	h := sha512.New()
	return hashFile(filename, h)
}

// hashByteA creates hash string from byte array
func hashByteA(bytea []byte, h hash.Hash) string {
	io.WriteString(h, string(bytea))
	return hash2string(h)
}

// hashFile creates hash string from file contents
func hashFile(filename string, h hash.Hash) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hash2string(h), nil
}

// hash2string converts hash to string
func hash2string(h hash.Hash) string {
	return fmt.Sprintf("%x", h.Sum(nil))
}
