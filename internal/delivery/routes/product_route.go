package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/savanyv/bsnack-backend/internal/cache"
	"github.com/savanyv/bsnack-backend/internal/database"
	"github.com/savanyv/bsnack-backend/internal/delivery/handlers"
	"github.com/savanyv/bsnack-backend/internal/repository"
	"github.com/savanyv/bsnack-backend/internal/usecase"
)

func productRoute(app fiber.Router, redisClient *cache.RedisClient) {
	repo := repository.NewProductRepository(database.DB)
	usecase := usecase.NewProductUsecase(repo, redisClient)
	handler := handlers.NewProductHandler(usecase)

	app.Get("/products", handler.GetAllProduct)
	app.Get("/products/:id", handler.GetByIDProduct)
	app.Post("/products", handler.CreateProduct)
}
