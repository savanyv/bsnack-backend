package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/savanyv/bsnack-backend/internal/database"
	"github.com/savanyv/bsnack-backend/internal/delivery/handlers"
	"github.com/savanyv/bsnack-backend/internal/repository"
	"github.com/savanyv/bsnack-backend/internal/usecase"
)

func transactionRegister(app fiber.Router) {
	tr := repository.NewTransactionRepository(database.DB)
	tir := repository.NewTransactionItemRepository(database.DB)
	cr := repository.NewCustomerRepository(database.DB)
	pr := repository.NewProductRepository(database.DB)
	prr := repository.NewPointRedemptionRepository(database.DB)
	usecase := usecase.NewTransactionUsecase(database.DB, cr, pr, tr, tir, prr)
	handler := handlers.NewTransactionHandler(usecase)

	app.Post("/transactions", handler.CreateTransaction)
	app.Post("/transactions/redeem", handler.RedeemPoint)
	app.Get("/transactions", handler.GetTransactionByPeriod)
}
