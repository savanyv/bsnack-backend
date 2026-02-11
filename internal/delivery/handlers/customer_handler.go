package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	dtos "github.com/savanyv/bsnack-backend/internal/dto"
	"github.com/savanyv/bsnack-backend/internal/usecase"
)

type CustomerHandler struct {
	customerUsecase usecase.CustomerUsecase
}

func NewCustomerHandler(cu usecase.CustomerUsecase) *CustomerHandler {
	return &CustomerHandler{
		customerUsecase: cu,
	}
}

func (h *CustomerHandler) GetCustomers(c *fiber.Ctx) error {
	monthParam := c.Query("month")

	var month time.Time
	if monthParam != "" {
		parsed, err := time.Parse("2006-01", monthParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Invalid month format",
				"message": "Month must be in YYYY-MM format.",
			})
		}
		month = parsed
	}

	result, err := h.customerUsecase.GetCustomers(c.Context(), dtos.CustomerQuery{Month: month})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve customers",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": result,
	})
}
