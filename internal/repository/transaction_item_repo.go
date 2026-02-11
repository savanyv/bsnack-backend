package repository

import (
	"context"
	"database/sql"

	"github.com/savanyv/bsnack-backend/internal/model"
)

type TransactionItemRepository interface {
	Create(ctx context.Context, tx *sql.Tx, item *model.TransactionItems) error
}

type transactionItemRepository struct {
	db *sql.DB
}

func NewTransactionItemRepository(db *sql.DB) TransactionItemRepository {
	return &transactionItemRepository{
		db: db,
	}
}

func (r *transactionItemRepository) Create(ctx context.Context, tx *sql.Tx, item *model.TransactionItems) error {
	query := `
		INSERT INTO transaction_items (transaction_id, product_id, quantity, price)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := tx.QueryRowContext(ctx, query, item.TransactionID, item.ProductID, item.Quantity, item.Price).Scan(
		&item.ID,
	)
	if err == sql.ErrNoRows {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}
