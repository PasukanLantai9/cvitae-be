package authHandler

import (
	"github.com/bccfilkom/career-path-service/internal/api/authentication"
	"github.com/gofiber/fiber/v2"
)

func (h *AuthHandler) HandleRegister(ctx *fiber.Ctx) error {
	var req authentication.CreateUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	if err := h.validator.Struct(&req); err != nil {
		return err
	}

	if err := h.authService.RegisterUser(ctx.Context(), req); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusCreated)
}
