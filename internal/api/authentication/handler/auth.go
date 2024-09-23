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

func (h *AuthHandler) HandleSignin(ctx *fiber.Ctx) error {
	var req authentication.SigninUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	if err := h.validator.Struct(&req); err != nil {
		return err
	}

	res, err := h.authService.SinginUser(ctx.Context(), req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (h *AuthHandler) HandleRefreshToken(ctx *fiber.Ctx) error {
	req, err := h.parseRefreshTokenRequest(ctx)
	if err != nil {
		return err
	}

	res, err := h.authService.GenerateNewAccessToken(ctx.Context(), req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}
