package handlers

import (
	"github.com/gofiber/fiber/v2"
	dtos "github.com/savanyv/bsnack-backend/internal/dto"
	"github.com/savanyv/bsnack-backend/internal/usecase"
)

type TransactionHandler struct {
	usecase usecase.TransactionUsecase
}

func NewTransactionHandler(u usecase.TransactionUsecase) *TransactionHandler {
	return &TransactionHandler{
		usecase: u,
	}
}

func (h *TransactionHandler) CreateTransaction(c *fiber.Ctx) error {
	var req dtos.CreateTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
	}

	res, err := h.usecase.CreateTransaction(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create transaction",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": 	"Successfully create transaction",
		"data": res,
	})
}

func (h *TransactionHandler) RedeemPoint(c *fiber.Ctx) error {
	var req dtos.ReedemPointRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
	}

	if req.CustomerID == "" || req.ProductID == "" || req.PointRequired <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "customer_id, product_id, and point_required are required",
		})
	}

	if err := h.usecase.ReedemPoint(c.Context(), req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to redeem point",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": 	"Successfully redeem point",
	})
}
