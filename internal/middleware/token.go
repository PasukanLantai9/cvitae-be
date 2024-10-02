package middleware

import (
	"github.com/bccfilkom/career-path-service/internal/api/authentication"
	"github.com/bccfilkom/career-path-service/internal/entity"
	"github.com/bccfilkom/career-path-service/internal/pkg/response"
	"github.com/bccfilkom/career-path-service/internal/pkg/token"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
)

var (
	ErrUnauthorized = response.NewError(http.StatusUnauthorized, "unauthorized, access token invalid or expired")
)

const (
	AccessTokenSecret = "JWT_ACCESS_TOKEN_SECRET"
)

func JWTAccessToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Get("Authorization") == "" {
			return ErrUnauthorized
		}

		if !strings.Contains(c.Get("Authorization"), "Bearer") {
			return ErrUnauthorized
		}

		userToken, err := token.VerifyTokenHeader(c, AccessTokenSecret)
		if err != nil {
			return ErrUnauthorized
		} else {
			claims := userToken.Claims.(jwt.MapClaims)
			user := authentication.UserClaims{
				ID:       claims["id"].(string),
				Provider: entity.AuthProvider(claims["provider"].(float64)),
				Email:    claims["email"].(string),
			}
			c.Locals("user", user)
			return c.Next()
		}
	}
}
