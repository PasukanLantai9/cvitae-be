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
	if id == "no-id" {
		return resume.ErrObjectIDNotProvided
	}

	res, err := h.resumeService.GetResumeByID(ctx.Context(), id, user.ID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (h *ResumeHandler) HandleUpdateResume(ctx *fiber.Ctx) error {
	var req resume.ResumeDetailDTO
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	user, err := helper.GetUserFromContext(ctx)
	if err != nil {
		return authentication.ErrUnauthorized
	}

	id := ctx.Params("id", "no-id")
	if id == "no-id" {
		return resume.ErrObjectIDNotProvided
	}

	req.ID = id
	req.UserID = user.ID

	reqEntity, _ := h.convertDTOToEntity(req)

	if err := h.resumeService.UpdateResumeByID(ctx.Context(), reqEntity); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *ResumeHandler) HandleScoringResume(ctx *fiber.Ctx) error {
	user, err := helper.GetUserFromContext(ctx)
	if err != nil {
		return authentication.ErrUnauthorized
	}

	id := ctx.Params("id", "no-id")
	if id == "no-id" {
		return resume.ErrObjectIDNotProvided
	}

	res, err := h.resumeService.ScoringResume(ctx.Context(), id, user.ID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (h *ResumeHandler) HandleScoringResumePDF(ctx *fiber.Ctx) error {
	user, err := helper.GetUserFromContext(ctx)
	if err != nil {
		return authentication.ErrUnauthorized
	}

	pdfFile, err := ctx.FormFile("resume")
	if err != nil {
		return resume.ErrPDFFileNotProvided
	}

	if pdfFile.Size == 0 {
		return resume.ErrPDFFileNotProvided
	}

	if !h.isPDF(pdfFile.Filename) {
		return resume.ErrInvalidResumeFile
	}

	res, err := h.resumeService.ScoringResumePDF(ctx.Context(), pdfFile, user.ID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (h *ResumeHandler) HandleJobVacancyFromResume(ctx *fiber.Ctx) error {
	user, err := helper.GetUserFromContext(ctx)
	if err != nil {
		return authentication.ErrUnauthorized
	}

	id := ctx.Params("id", "no-id")
	if id == "no-id" {
		return resume.ErrObjectIDNotProvided
	}

	res, err := h.resumeService.JobVacancyFromResume(ctx.Context(), id, user.ID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (h *ResumeHandler) HandleJobVacancyFromPDF(ctx *fiber.Ctx) error {
	user, err := helper.GetUserFromContext(ctx)
	if err != nil {
		return authentication.ErrUnauthorized
	}

	pdfFile, err := ctx.FormFile("resume")
	if err != nil {
		return resume.ErrPDFFileNotProvided
	}

	if pdfFile.Size == 0 {
		return resume.ErrPDFFileNotProvided
	}

	if !h.isPDF(pdfFile.Filename) {
		return resume.ErrInvalidResumeFile
	}

	res, err := h.resumeService.JobVacancyFromPDF(ctx.Context(), pdfFile, user.ID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}
