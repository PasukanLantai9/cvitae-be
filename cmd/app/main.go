package main

import (
	"github.com/bccfilkom/career-path-service/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	logger := config.NewLogger()

	if err := godotenv.Load(); err != nil {
		logger.Fatal("Error loading .env file")
	}

	fiber := config.NewFiber(logger)
	validator := config.NewValidator()
	rest, err := config.NewServer(fiber, logger, validator)
	if err != nil {
		logger.Fatal(err)
	}

	if err := rest.Run(); err != nil {
		logger.Fatal(err)
	}
}
