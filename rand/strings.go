package rand

import (
	"crypto/rand"
	"encoding/base64"
)

// RememberTokenBytes is length of token
const RememberTokenBytes = 32

// Bytes generate random bytes
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)

	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// String returns random string
func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// RememberToken generate token
func RememberToken() (string, error) {
	return String(RememberTokenBytes)
}
