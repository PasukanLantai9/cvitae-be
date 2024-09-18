package config

import (
	"fmt"
	"github.com/bccfilkom/career-path-service/internal/pkg/env"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	app      *fiber.App
	db       *sqlx.DB
	handlers []handler
	log      *logrus.Logger
}

type handler interface {
	Start(srv fiber.Router)
}

func NewServer(fiberApp *fiber.App, log *logrus.Logger, validator *validator.Validate) (*Server, error) {
	bootstrap := &Server{
		app: fiberApp,
		log: log,
	}

	return bootstrap, nil
}

func (s *Server) RegisterHandler() {
	s.checkHealth()
	s.handlers = append(s.handlers)
}

func (s *Server) Run() error {
	router := s.app.Group("/api/v1")

	for _, h := range s.handlers {
		h.Start(router)
	}

	addr := env.GetString("APP_ADDR", "127.0.0.1")
	port := env.GetString("APP_PORT", "8000")

	if !env.GetBool("PRODUCTION", false) {
		s.log.Warnf("Starting server on non-production mode at http://%s:%s", addr, port)
	}

	if err := s.app.Listen(fmt.Sprintf("%s:%s", addr, port)); err != nil {
		return err
	}

	return nil
}

func (s *Server) checkHealth() {
	s.app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"message": "🚗💨Beep Beep Your Server is Healthy!",
		})
	})
}
