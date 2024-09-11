package hash

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

func WithMD5(messages []string) (string, error) {
	// Concatenate messages into a single string
	concatenated := strings.Join(messages, "")

	// Convert concatenated string to byte slice
	concatenatedBytes := []byte(concatenated)

	// Hash concatenated string using MD5
	hasher := md5.New()
	_, err := hasher.Write(concatenatedBytes)
	if err != nil {
		return "", err
	}
	hashBytes := hasher.Sum(nil)

	// Convert hash bytes to hex string
	hashString := hex.EncodeToString(hashBytes)

	return hashString, nil
}
