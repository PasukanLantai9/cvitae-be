package resumeHandler

import (
	resumeService "github.com/bccfilkom/career-path-service/internal/api/resume/service"
	"github.com/bccfilkom/career-path-service/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ResumeHandler struct {
	resumeService resumeService.ResumeService
	log           *logrus.Logger
	validator     *validator.Validate
}

func New(resumeService resumeService.ResumeService, log *logrus.Logger, validate *validator.Validate) *ResumeHandler {
	return &ResumeHandler{
		resumeService: resumeService,
		log:           log,
		validator:     validate,
	}
}

func (h *ResumeHandler) Start(srv fiber.Router) {
	resume := srv.Group("/resume")
	resume.Use(middleware.JWTAccessToken())
	resume.Post("/", h.HandleCreateResume)
	resume.Get("/", h.HandleGetUserResume)
	resume.Get("/:id", h.HandleGetResumeByID)
	resume.Put("/:id", h.HandleUpdateResume)

	resume.Get("/scoring/:id", h.HandleScoringResume)
	resume.Post("/scoring/file", h.HandleScoringResumePDF)
	resume.Get("/job-vacancy/:id")
	resume.Post("/job-vacancy/file")
}
