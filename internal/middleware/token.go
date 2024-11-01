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

type tokenMiddleware struct {
}

func newTokenMiddleware() *tokenMiddleware {
	return &tokenMiddleware{}
}

func (m *middleware) NewtokenMiddleware(ctx *fiber.Ctx) error {
	clientIP := ctx.IP()

	if ctx.Get("Authorization") == "" {
		m.log.Warnf("authorization header is missing, client ip is %s", clientIP)
		return ErrUnauthorized
	}

	if !strings.Contains(ctx.Get("Authorization"), "Bearer") {
		m.log.Warnf("authorization header is missing, client ip is %s", clientIP)
		return ErrUnauthorized
	}

	userToken, err := token.VerifyTokenHeader(ctx, AccessTokenSecret)
	if err != nil {
		return ErrUnauthorized
	} else {
		claims := userToken.Claims.(jwt.MapClaims)
		user := authentication.UserClaims{
			ID:       claims["id"].(string),
			Provider: entity.AuthProvider(claims["provider"].(float64)),
			Email:    claims["email"].(string),
		}
		ctx.Locals("user", user)
		return ctx.Next()
	}
}
