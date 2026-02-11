package repository

import (
	"context"
	"database/sql"

	"github.com/savanyv/bsnack-backend/internal/model"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *sql.Tx, transaction *model.Transaction) error
	UpdateTotal(ctx context.Context, tx *sql.Tx, transactionID string, totalPrice int) error
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) Create(ctx context.Context, tx *sql.Tx, transaction *model.Transaction) error {
	query := `
		INSERT INTO transactions (customer_id, total_price)
		VALUES ($1, $2)
		RETURNING id, transaction_date
	`

	err := tx.QueryRowContext(ctx, query, transaction.CustomerID, transaction.TotalPrice).Scan(
		&transaction.ID,
		&transaction.TransactionDate,
	)
	if err == sql.ErrNoRows {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}

func (r transactionRepository) UpdateTotal(ctx context.Context, tx *sql.Tx, transactionID string, totalPrice int) error {
	query := `
		UPDATE transactions
		SET total_price = $1
		WHERE id = $2
	`

	_, err := tx.ExecContext(ctx, query, totalPrice, transactionID)
	return err
}
