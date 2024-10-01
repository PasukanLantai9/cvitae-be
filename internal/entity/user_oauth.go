package entity

import "time"

type UserOauth struct {
	ID           string
	UserID       string
	Provider     AuthProvider
	OAuthUserID  string
	AccessToken  string
	RefreshToken string
	CreatedAt    time.Time
}
