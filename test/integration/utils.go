package integration

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bccfilkom/career-path-service/internal/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"io"
	"net/http"
	"os"
)

var app = fiber.New(fiber.Config{
	ErrorHandler: newErrorHandler(),
})

func GenerateDSN() string {
	DBName := os.Getenv("DB_NAME")
	DBUser := os.Getenv("DB_USER")
	DBPassword := os.Getenv("DB_PASSWORD")
	DBHost := os.Getenv("DB_HOST")
	DBPort := os.Getenv("DB_PORT")

	DBDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		DBHost, DBUser, DBPassword, DBName, DBPort,
	)

	return DBDSN
}

func newErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		var apiErr *response.Error
		if errors.As(err, &apiErr) {
			return ctx.Status(apiErr.Code).JSON(fiber.Map{
				"errors": fiber.Map{"message": apiErr.Error()},
			})
		}

		var validationErr validator.ValidationErrors
		if errors.As(err, &validationErr) {
			fieldErr := fiber.Map{}
			for _, e := range validationErr {
				fieldErr[e.Field()] = e.Error()
			}
			fieldErr["message"] = utils.StatusMessage(fiber.StatusUnprocessableEntity)
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"errors": fieldErr,
			})
		}

		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			return ctx.Status(fiberErr.Code).JSON(fiber.Map{
				"errors": fiber.Map{"message": utils.StatusMessage(fiberErr.Code), "err": err},
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": fiber.Map{"message": utils.StatusMessage(fiber.StatusInternalServerError), "err": err.Error()},
		})
	}
}

func printResponse(response *http.Response) {
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("failed to read response body: %s\n", err.Error())
	}

	var responseBody map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &responseBody); err != nil {
		fmt.Printf("Failed to unmarshal response body: %s\n", err.Error())
	} else {
		fmt.Printf("Parsed response body: %#v\n", responseBody)
	}
}
