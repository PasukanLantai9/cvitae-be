package middleware

import "github.com/gofiber/fiber/v2"

type loggingMiddleware struct {
}

func newloggingMiddleware() *loggingMiddleware {
	return &loggingMiddleware{}
}

func (m *middleware) NewLoggingMiddleware(ctx *fiber.Ctx) error {
	clientIP := ctx.IP()
	bodyRequest := ctx.Request().Body()

	m.log.Infof("clientIP: %s, bodyRequest: %s", clientIP, bodyRequest)

	return ctx.Next()
}
