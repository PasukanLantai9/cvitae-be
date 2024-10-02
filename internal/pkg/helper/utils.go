package helper

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/bccfilkom/career-path-service/internal/api/authentication"
	"github.com/gofiber/fiber/v2"
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

func GetUserFromContext(ctx *fiber.Ctx) (authentication.UserClaims, error) {
	userCtx := ctx.Locals("user")
	if userCtx == nil {
		return authentication.UserClaims{}, fmt.Errorf("user not found in context")
	}

	return userCtx.(authentication.UserClaims), nil
}
