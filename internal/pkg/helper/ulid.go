package helper

import (
	"crypto/rand"
	"github.com/oklog/ulid/v2"
	"time"
)

func NewUlidFromTimestamp(time time.Time) (string, error) {
	ms := ulid.Timestamp(time)
	entropy := ulid.Monotonic(rand.Reader, 0)

	id, err := ulid.New(ms, entropy)
	if err != nil {
		return "", err
	}

	return id.String(), nil
}
