package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/savanyv/bsnack-backend/internal/database"
	"github.com/savanyv/bsnack-backend/internal/delivery/handlers"
	"github.com/savanyv/bsnack-backend/internal/repository"
	"github.com/savanyv/bsnack-backend/internal/usecase"
)

func customerRoute(app fiber.Router) {
	repo := repository.NewCustomerRepository(database.DB)
	usecase := usecase.NewCustomerUsecase(repo)
	handler := handlers.NewCustomerHandler(usecase)

	app.Get("/customers", handler.GetCustomers)
}
