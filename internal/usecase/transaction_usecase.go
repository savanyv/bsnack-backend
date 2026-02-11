package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	dtos "github.com/savanyv/bsnack-backend/internal/dto"
	"github.com/savanyv/bsnack-backend/internal/model"
	"github.com/savanyv/bsnack-backend/internal/repository"
)

type TransactionUsecase interface {
	CreateTransaction(ctx context.Context, req dtos.CreateTransactionRequest) (*dtos.CreateTransactionResponse, error)
	ReedemPoint(ctx context.Context, req dtos.ReedemPointRequest) error
	GetTransactionByPeriod(ctx context.Context, startDate, endDate string) (*dtos.TransactionSummaryResponse, error)
}

type transactionUsecase struct {
	db *sql.DB
	customerRepo repository.CustomerRepository
	productRepo repository.ProductRepository
	transactionRepo repository.TransactionRepository
	transactionItemRepo repository.TransactionItemRepository
	pointRedemptionRepo repository.PointRedemptionRepository
}

func NewTransactionUsecase(db *sql.DB, cr repository.CustomerRepository, pr repository.ProductRepository, tr repository.TransactionRepository, tir repository.TransactionItemRepository, prr repository.PointRedemptionRepository) TransactionUsecase {
	return &transactionUsecase{
		db: db,
		customerRepo: cr,
		productRepo: pr,
		transactionRepo: tr,
		transactionItemRepo: tir,
		pointRedemptionRepo: prr,
	}
}

func (u *transactionUsecase) CreateTransaction(ctx context.Context, req dtos.CreateTransactionRequest) (*dtos.CreateTransactionResponse, error) {
	if req.CustomerName == "" || len(req.Items) == 0 {
		return nil, errors.New("customer_name and items are required")
	}

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	customer, err := u.customerRepo.FindByName(ctx, req.CustomerName)
	if err != nil {
		return nil, err
	}

	if customer == nil {
		customer = &model.Customer{
			Name: req.CustomerName,
			Point: 0,
		}
		if err := u.customerRepo.Create(ctx, customer); err != nil {
			return nil, err
		}
	}

	totalPrice := 0

	transaction := &model.Transaction{
		CustomerID: customer.ID,
		TotalPrice: 0,
	}

	if err := u.transactionRepo.Create(ctx, tx, transaction); err != nil {
		return nil, err
	}

	for _, item := range req.Items {
		if item.ProductID == "" || item.Quantity < 0 {
			return nil, errors.New("product_id and quantity are required")
		}

		product, err := u.productRepo.FindByID(ctx, item.ProductID)
		if err != nil {
			return nil, err
		}

		if product == nil {
			return nil, errors.New("product not found")
		}
		if product.Stock < item.Quantity {
			return nil, errors.New("stock not enough")
		}

		itemTotal := product.Price * item.Quantity
		totalPrice += itemTotal

		transactionItem := &model.TransactionItems{
			TransactionID: transaction.ID,
			ProductID: product.ID,
			Quantity: item.Quantity,
			Price: itemTotal,
		}

		if err := u.transactionItemRepo.Create(ctx, tx, transactionItem); err != nil {
			return nil, err
		}

		if err := u.productRepo.UpdateStock(ctx, tx, product.ID.String(), product.Stock - item.Quantity); err != nil {
			return nil, err
		}
	}

	transaction.TotalPrice = totalPrice
	if err := u.transactionRepo.UpdateTotal(ctx, tx, transaction.ID.String(), transaction.TotalPrice); err != nil {
		return nil, err
	}

	pointEarned := totalPrice / 1000
	if err := u.customerRepo.UpdatePoint(ctx, tx, customer.ID.String(), customer.Point + pointEarned); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	response := dtos.CreateTransactionResponse{
		TransactionID: transaction.ID.String(),
		TotalPrice: transaction.TotalPrice,
		PointEarned: pointEarned,
	}

	return &response, nil
}

func (u *transactionUsecase) ReedemPoint(ctx context.Context, req dtos.ReedemPointRequest) error {
	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		return errors.New("invalid customer id")
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		return errors.New("invalid product id")
	}

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	customer, err := u.customerRepo.FindByID(ctx, customerID.String())
	if err != nil {
		return err
	}

	if customer.Point < req.PointRequired {
		return errors.New("point not enough")
	}

	product, err := u.productRepo.FindByID(ctx, productID.String())
	if err != nil {
		return err
	}

	if product.Stock < 1 {
		return errors.New("stock empty")
	}

	redemption := &model.PointRedemption{
		CustomerID: customer.ID,
		ProductID: product.ID,
		PointUsed: req.PointRequired,
	}

	if err := u.pointRedemptionRepo.Create(ctx, tx, redemption); err != nil {
		return err
	}

	if err := u.customerRepo.UpdatePoint(ctx, tx, customer.ID.String(), customer.Point - req.PointRequired); err != nil {
		return err
	}

	if err := u.productRepo.UpdateStock(ctx, tx, product.ID.String(), product.Stock - 1); err != nil {
		return err
	}

	return tx.Commit()
}

func (u *transactionUsecase) GetTransactionByPeriod(ctx context.Context, startDate, endDate string) (*dtos.TransactionSummaryResponse, error) {
	transactions, err := u.transactionRepo.GetByPeriod(ctx, startDate, endDate)
	if err != nil {
		return nil, errors.New("failed to get transactions")
	}

	totalCustomer, totalIncome, bestSeller, err := u.transactionRepo.GetSummary(ctx, startDate, endDate)

	var list []dtos.TransactionListItem
	for _, t := range transactions {
		list = append(list, dtos.TransactionListItem{
			ID: t.ID.String(),
			CustomerID: t.CustomerID.String(),
			TotalPrice: t.TotalPrice,
			CreatedAt: t.TransactionDate.Format("2006-01-02 15:04:05"),
		})
	}

	response := dtos.TransactionSummaryResponse{
		StartDate: startDate,
		EndDate: endDate,
		TotalCustomer: totalCustomer,
		TotalIncome: totalIncome,
		BestSeller: bestSeller,
		Transactions: list,
	}

	return &response, nil
}
