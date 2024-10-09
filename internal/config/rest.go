package config

import (
	"fmt"
	authHandler "github.com/bccfilkom/career-path-service/internal/api/authentication/handler"
	authRepository "github.com/bccfilkom/career-path-service/internal/api/authentication/repository"
	authService "github.com/bccfilkom/career-path-service/internal/api/authentication/service"
	resumeHandler "github.com/bccfilkom/career-path-service/internal/api/resume/handler"
	resumeRepository "github.com/bccfilkom/career-path-service/internal/api/resume/repository"
	resumeService "github.com/bccfilkom/career-path-service/internal/api/resume/service"
	"github.com/bccfilkom/career-path-service/internal/pkg/cronjob"
	"github.com/bccfilkom/career-path-service/internal/pkg/env"
	"github.com/bccfilkom/career-path-service/pkg/google"
	"github.com/bccfilkom/career-path-service/pkg/mongo"
	"github.com/bccfilkom/career-path-service/pkg/postgres"
	redisdb "github.com/bccfilkom/career-path-service/pkg/redis"
	jobMatching "github.com/bccfilkom/career-path-service/pkg/rpc/job_matching"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Server struct {
	app       *fiber.App
	db        *sqlx.DB
	handlers  []handler
	log       *logrus.Logger
	validator *validator.Validate
	google    google.Google
}

type handler interface {
	Start(srv fiber.Router)
}

func NewServer(fiberApp *fiber.App, log *logrus.Logger, validator *validator.Validate, googles google.Google) (*Server, error) {
	bootstrap := &Server{
		app:       fiberApp,
		log:       log,
		validator: validator,
		google:    googles,
	}

	return bootstrap, nil
}

func (s *Server) registerHandler() {
	// third party dependecies
	db, err := postgres.NewInstance()
	if err != nil {
		s.log.Fatalf("could not connect to database: %v", err)
	}

	mongodb, err := mongo.NewMongoInstance()
	if err != nil {
		s.log.Fatalf("could not connect to mongodb: %v", err)
	}

	redisClient, err := redisdb.NewInstance()
	if err != nil {
		s.log.Fatalf("could not connect to redisdb: %v", err)
	}

	cronjobClient := cronjob.NewInstance()
	mlRPCClient := jobMatching.NewRpcClient(s.log)

	// repository
	authRepos := authRepository.New(db)
	resumeRepos := resumeRepository.New(mongodb, db, redisClient)

	// service
	authServices := authService.New(authRepos, s.google)
	resumeServices := resumeService.New(resumeRepos, mlRPCClient)

	// handler
	authHandlers := authHandler.New(authServices, s.validator)
	resumeHandlers := resumeHandler.New(resumeServices, s.log, s.validator)

	cronjobClient.SetupCronJob(time.Minute*10, func() error {
		return resumeServices.SyncResumesFromRedisToMongo(redisClient, mongodb.Collection("resume"))
	})

	s.checkHealth()
	s.handlers = append(s.handlers, authHandlers, resumeHandlers)
}

func (s *Server) Run() error {
	router := s.app.Group("/api/v1")

	s.registerHandler()

	s.app.Use(cors.New())
	s.app.Use(logger.New())
	for _, h := range s.handlers {
		h.Start(router)
	}

	addr := env.GetString("APP_ADDR", "0.0.0.0")
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
