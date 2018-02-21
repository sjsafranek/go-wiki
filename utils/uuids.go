package utils

import (
	crand "crypto/rand"
	"fmt"
	"io"
	"math/rand"
	"time"
)

const _letters string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// NewUUID generates and returns a uuid
func NewUUID() (string, error) {
	b := make([]byte, 16)
	n, err := io.ReadFull(crand.Reader, b)
	if n != len(b) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	b[8] = b[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	b[6] = b[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

// newJobId generates and returns an job_id of desired length
func newJobId(n int) string {
	s := ""
	for i := 1; i <= n; i++ {
		s += string(_letters[rand.Intn(len(_letters))])
	}
	return s
}
