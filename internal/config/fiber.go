package config

import (
	"errors"
	"github.com/bccfilkom/career-path-service/internal/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

func NewFiber(log *logrus.Logger) *fiber.App {
	app := fiber.New(
		fiber.Config{
			AppName:           "Career Path Backend",
			BodyLimit:         50 * 1024 * 1024,
			DisableKeepalive:  true,
			StrictRouting:     true,
			CaseSensitive:     true,
			EnablePrintRoutes: true,
			ErrorHandler:      newErrorHandler(log),
			JSONEncoder:       jsoniter.Marshal,
			JSONDecoder:       jsoniter.Unmarshal,
		})
	return app
}

func newErrorHandler(log *logrus.Logger) fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		var apiErr *response.Error
		if errors.As(err, &apiErr) {
			log.Errorf("client ip %s, error %s", ctx.IP(), apiErr.Error())
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
			log.Errorf("client ip %s, error %s", ctx.IP(), fiberErr.Error())
			return ctx.Status(fiberErr.Code).JSON(fiber.Map{
				"errors": fiber.Map{"message": utils.StatusMessage(fiberErr.Code), "err": err},
			})
		}

		log.Errorf("client ip %s, error %s", ctx.IP(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": fiber.Map{"message": utils.StatusMessage(fiber.StatusInternalServerError), "err": err},
		})
	}
}
