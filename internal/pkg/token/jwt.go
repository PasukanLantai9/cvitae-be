package token

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"os"
	"strings"
	"time"
)

func Sign(Data map[string]interface{}, secretEnvKey string, ExpiredAt time.Duration) (string, error) {
	expiredAt := time.Now().Add(ExpiredAt).Unix()

	JWTSecretKey := os.Getenv(secretEnvKey)

	claims := jwt.MapClaims{}
	claims["exp"] = expiredAt
	claims["authorization"] = true

	for i, v := range Data {
		claims[i] = v
	}

	to := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := to.SignedString([]byte(JWTSecretKey))

	if err != nil {
		return accessToken, err
	}

	return accessToken, nil
}

func VerifyTokenHeader(c *fiber.Ctx, secretEnvKey string) (*jwt.Token, error) {
	header := c.Get("Authorization")
	accessToken := strings.SplitAfter(header, "Bearer")[1]
	JWTSecretKey := os.Getenv(secretEnvKey)

	token, err := jwt.Parse(strings.Trim(accessToken, " "), func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
