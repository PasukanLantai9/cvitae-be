package entity

import "time"

type Session struct {
	ID           string
	UserID       string
	RefreshToken string
	CreatedAt    string
	ExpiresAt    time.Time
	AuthProvider AuthProvider
}

type AuthProvider uint8

const (
	AuthProviderUnknown  AuthProvider = 0
	AuthProviderNative   AuthProvider = 1
	AuthProviderGoogle   AuthProvider = 2
	AuthProviderLinkedIn AuthProvider = 3
)

var AuthProviderMap = map[AuthProvider]string{
	AuthProviderNative:   "Native",
	AuthProviderGoogle:   "Google",
	AuthProviderLinkedIn: "LinkedIn",
}

func (a AuthProvider) String() string {
	return AuthProviderMap[a]
}

func (s AuthProvider) Value() uint8 {
	return uint8(s)
}
