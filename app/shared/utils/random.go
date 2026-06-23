package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func RandomBool() bool {
	n, err := rand.Int(rand.Reader, big.NewInt(2)) // returns 0 or 1
	if err != nil {
		return false
	}
	return n.Int64() == 1
}

func RandomInt64(min, max int64) int64 {
	if min >= max {
		return min
	}

	diff := max - min + 1
	n, err := rand.Int(rand.Reader, big.NewInt(diff))
	if err != nil {
		return min
	}

	return min + n.Int64()
}

func RandomPassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"0123456789"
	charsetLen := big.NewInt(int64(len(charset)))

	b := make([]byte, length)

	for i := range b {
		n, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		b[i] = charset[n.Int64()]
	}

	return string(b), nil
}

func RandomString(length int, numbers bool, uppercase bool, lowercase bool) (string, error) {
	if length < 1 {
		return "", nil
	}
	var charset string
	if numbers {
		charset += "0123456789"
	}
	if uppercase {
		charset += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	if lowercase {
		charset += "abcdefghijklmnopqrstuvwxyz"
	}

	if charset == "" {
		return "", fmt.Errorf("charset is empty")
	}

	charsetLen := big.NewInt(int64(len(charset)))

	b := make([]byte, length)

	for i := range b {
		n, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		b[i] = charset[n.Int64()]
	}

	return string(b), nil
}
