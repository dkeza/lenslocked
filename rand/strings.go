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

// NBytes returns number of bytes used in the base64 url encoded string
func NBytes(base64String string) (int, error) {
	b, err := base64.URLEncoding.DecodeString(base64String)
	if err != nil {
		return -1, err
	}
	return len(b), nil
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
