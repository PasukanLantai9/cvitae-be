package helper

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateRandomString(length int) string {
	randomBytes := make([]byte, length)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return ""
	}

	str := base64.RawURLEncoding.EncodeToString(randomBytes)
	return str[:length]
}
