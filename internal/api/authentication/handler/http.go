package authHandler

import (
	authService "github.com/bccfilkom/career-path-service/internal/api/authentication/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService authService.AuthService
	validator   *validator.Validate
}

func New(authService authService.AuthService, validate *validator.Validate) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validate,
	}
}

func (h *AuthHandler) Start(srv fiber.Router) {
	auth := srv.Group("/auth")

	auth.Post("/register", h.HandleRegister)
	auth.Post("/signin", h.HandleSignin)
	auth.Post("/refresh", h.HandleRefreshToken)

	oauth := srv.Group("/oauth")
	oauth.Get("/:provider", h.HandleOauthCallback)
}
