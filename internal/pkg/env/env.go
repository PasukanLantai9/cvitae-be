package env

import (
	"os"
	"strconv"
)

func GetString(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return value
}

func GetInt(key string, fallback int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return i
}

func GetBool(key string, fallback bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	b, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return b
}
