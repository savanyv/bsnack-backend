package repository

import (
	"context"
	"database/sql"

	"github.com/savanyv/bsnack-backend/internal/model"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *sql.Tx, transaction *model.Transaction) error
	UpdateTotal(ctx context.Context, tx *sql.Tx, transactionID string, totalPrice int) error
	GetByPeriod(ctx context.Context, startDate, endDate string) ([]model.Transaction, error)
	GetSummary(ctx context.Context, startDate, endDate string) (int, int, string, error)
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

func (r *transactionRepository) GetByPeriod(ctx context.Context, startDate, endDate string) ([]model.Transaction, error) {
	query := `
		SELECT id, customer_id, total_price, transaction_date
		FROM transactions
		WHERE DATE(transaction_date) BETWEEN $1 AND $2
		ORDER BY transaction_date DESC
	`

	rows, err := r.db.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var t model.Transaction
		if err := rows.Scan(
			&t.ID,
			&t.CustomerID,
			&t.TotalPrice,
			&t.TransactionDate,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (r *transactionRepository) GetSummary(ctx context.Context, startDate, endDate string) (int, int, string, error) {
	var totalCustomer int
	var totalIncome int
	var bestSeller string

	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(DISTINCT customer_id)
		FROM transactions
		WHERE transaction_date >= $1::date
		AND transaction_date <= ($2::date + INTERVAL '1 day')
	`, startDate, endDate).Scan(&totalCustomer)
	if err != nil {
		return 0, 0, "", err
	}

	err = r.db.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(total_price), 0)
		FROM transactions
		WHERE transaction_date >= $1::date
		AND transaction_date <= ($2::date + INTERVAL '1 day')
	`, startDate, endDate).Scan(&totalIncome)
	if err != nil {
		return 0, 0, "", err
	}

	err = r.db.QueryRowContext(ctx, `
		SELECT p.name
		FROM transaction_items ti
		JOIN products p ON p.id = ti.product_id
		JOIN transactions t ON t.id = ti.transaction_id
		WHERE t.transaction_date >= $1::date
		AND t.transaction_date < ($2::date + INTERVAL '1 day')
		GROUP BY p.name
		ORDER BY SUM(ti.quantity) DESC
		LIMIT 1
	`, startDate, endDate).Scan(&bestSeller)
	if err != nil {
		return 0, 0, "", err
	}

	return totalCustomer, totalIncome, bestSeller, nil
}
