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
	middleware    middleware.Middleware
	log           *logrus.Logger
	validator     *validator.Validate
}

func New(resumeService resumeService.ResumeService,
	log *logrus.Logger,
	validate *validator.Validate,
	middleware middleware.Middleware) *ResumeHandler {
	return &ResumeHandler{
		resumeService: resumeService,
		log:           log,
		validator:     validate,
		middleware:    middleware,
	}
}

func (h *ResumeHandler) Start(srv fiber.Router) {
	resume := srv.Group("/resume")
	resume.Use(h.middleware.NewtokenMiddleware)
	resume.Post("", h.HandleCreateResume)
	resume.Get("", h.HandleGetUserResume)
	resume.Get("/:id", h.HandleGetResumeByID)
	resume.Put("/:id", h.HandleUpdateResume)

	resume.Get("/scoring/:id", h.HandleScoringResume)
	resume.Post("/scoring/file", h.HandleScoringResumePDF)
	resume.Get("/job-vacancy/:id", h.HandleJobVacancyFromResume)
	resume.Post("/job-vacancy/file", h.HandleJobVacancyFromPDF)
}
