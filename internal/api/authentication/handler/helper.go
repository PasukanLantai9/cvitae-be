package authHandler

import (
	"github.com/bccfilkom/career-path-service/internal/api/authentication"
	"github.com/bccfilkom/career-path-service/internal/entity"
	"github.com/gofiber/fiber/v2"
)

func (h *AuthHandler) parseRefreshTokenRequest(ctx *fiber.Ctx) (entity.Session, error) {
	sessionID := ctx.Query("sessionID")
	if sessionID == "" {
		return entity.Session{}, fiber.ErrUnprocessableEntity
	}

	var req authentication.RefresTokenRequest
	if err := ctx.BodyParser(&req); err != nil {
		return entity.Session{}, err
	}

	if err := h.validator.Struct(&req); err != nil {
		return entity.Session{}, err
	}

	return entity.Session{
		ID:           sessionID,
		RefreshToken: req.RefreshToken,
	}, nil
}
