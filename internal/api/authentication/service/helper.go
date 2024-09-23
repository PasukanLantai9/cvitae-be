package authService

import (
	"fmt"
	"github.com/bccfilkom/career-path-service/internal/entity"
	"github.com/bccfilkom/career-path-service/internal/pkg/helper"
	"github.com/bccfilkom/career-path-service/internal/pkg/token"
	"time"
)

const (
	JWTSecret = "JWT_SECRET"
)

func generateFreshToken(user entity.User, provider entity.AuthProvider) (string, string, error) {
	userClaims := map[string]interface{}{
		"id":       user.ID,
		"email":    user.Email,
		"provider": provider,
	}

	accessToken, err := token.Sign(userClaims, JWTSecret, 30*time.Minute)
	if err != nil {
		return "", "", err
	}

	refreshToken := fmt.Sprintf("%s-%d", helper.GenerateRandomString(15), time.Now().UnixMilli())

	return accessToken, refreshToken, nil
}

func generateAccessToken(user entity.User, provider entity.AuthProvider) (string, error) {
	userClaims := map[string]interface{}{
		"id":       user.ID,
		"email":    user.Email,
		"provider": provider,
	}

	accessToken, err := token.Sign(userClaims, JWTSecret, 30*time.Minute)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
