package middlewares

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func MethodValidationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		allowedMethods := map[string]bool{
			http.MethodGet: true,
			http.MethodPost: true,
			http.MethodPut: true,
			http.MethodDelete: true,
			http.MethodOptions: true,
		}

		if !allowedMethods[c.Method()] {
			return c.Status(http.StatusMethodNotAllowed).JSON(fiber.Map{
				"error": "Method Not Allowed",
				"message": "The requested HTTP method is not allowed for this endpoint.",
			})
		}

		return c.Next()
	}
}
