package handlers

import (
	"github.com/gofiber/fiber/v2"
	dtos "github.com/savanyv/bsnack-backend/internal/dto"
	"github.com/savanyv/bsnack-backend/internal/usecase"
)

type ProductHandler struct {
	productUsecase usecase.ProductUsecase
}

func NewProductHandler(pu usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{
		productUsecase: pu,
	}
}

func (h *ProductHandler) GetAllProduct(c *fiber.Ctx) error {
	products, err := h.productUsecase.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve products",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": 	"Successfully retrieve products",
		"data": products,
	})
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req dtos.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
	}

	product, err := h.productUsecase.CreateProduct(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create product",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": 	"Successfully create product",
		"data": product,
	})
}

func (h *ProductHandler) GetByIDProduct(c *fiber.Ctx) error {
	ID := c.Params("id")

	product, err := h.productUsecase.GetByID(c.Context(), ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve product",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": 	"Successfully retrieve product",
		"data": product,
	})
}
