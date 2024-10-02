package resumeHandler

import (
	"github.com/bccfilkom/career-path-service/internal/api/authentication"
	"github.com/bccfilkom/career-path-service/internal/api/resume"
	"github.com/bccfilkom/career-path-service/internal/pkg/helper"
	"github.com/gofiber/fiber/v2"
)

func (h *ResumeHandler) HandleCreateResume(ctx *fiber.Ctx) error {
	var req resume.ResumeRequest
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	if err := h.validator.Struct(&req); err != nil {
		return err
	}

	user, err := helper.GetUserFromContext(ctx)
	if err != nil {
		return authentication.ErrUnauthorized
	}

	id, err := h.resumeService.CreateResume(ctx.Context(), req, user)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id})

}

func (h *ResumeHandler) HandleGetUserResume(ctx *fiber.Ctx) error {
	user, err := helper.GetUserFromContext(ctx)
	if err != nil {
		return authentication.ErrUnauthorized
	}

	res, err := h.resumeService.GetUserResume(ctx.Context(), user.ID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (h *ResumeHandler) HandleGetResumeByID(ctx *fiber.Ctx) error {
	user, err := helper.GetUserFromContext(ctx)
	if err != nil {
		return authentication.ErrUnauthorized
	}

	id := ctx.Params("id", "no-id")

	res, err := h.resumeService.GetResumeByID(ctx.Context(), id, user.ID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}
