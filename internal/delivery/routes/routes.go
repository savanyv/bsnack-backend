package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/savanyv/bsnack-backend/internal/cache"
)

func RegisterRoutes(app fiber.Router, redisClient *cache.RedisClient) {
	api := app.Group("/bsnack-api")

	customerRoute(api)
	productRoute(api, redisClient)
	transactionRegister(api)
}
