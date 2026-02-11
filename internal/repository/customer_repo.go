package repository

import (
	"context"
	"database/sql"

	"github.com/savanyv/bsnack-backend/internal/model"
)

type CustomerRepository interface {
	FindByName(ctx context.Context, name string) (*model.Customer, error)
	FindByID(ctx context.Context, id string) (*model.Customer, error)
	Create(ctx context.Context, customer *model.Customer) error
	UpdatePoint(ctx context.Context, tx *sql.Tx, ID string, point int) error
	GetAll(ctx context.Context) ([]*model.Customer, error)
}

type customerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) CustomerRepository {
	return &customerRepository{
		db: db,
	}
}

func (r *customerRepository) FindByName(ctx context.Context, name string) (*model.Customer, error) {
	var customer model.Customer
	query := `
		SELECT id, name, point, created_at
		FROM customers
		WHERE name = $1
	`
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&customer.ID,
		&customer.Name,
		&customer.Point,
		&customer.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &customer, nil
}

func (r *customerRepository) FindByID(ctx context.Context, ID string) (*model.Customer, error) {
	var customer model.Customer
	query := `
		SELECT id, name, point, created_at
		FROM customers
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, ID).Scan(
		&customer.ID,
		&customer.Name,
		&customer.Point,
		&customer.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &customer, err
}

func (r *customerRepository) Create(ctx context.Context, customer *model.Customer) error {
	query := `
		INSERT INTO customers (name, point)
		VALUES ($1, $2)
		RETURNING id, created_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		customer.Name,
		customer.Point,
	).Scan(
		&customer.ID,
		&customer.CreatedAt,
	)
}

func (r *customerRepository) UpdatePoint(ctx context.Context, tx *sql.Tx, ID string, point int) error {
	query := `
		UPDATE customers
		SET point = $1
		WHERE id = $2
	`

	_, err := tx.ExecContext(ctx, query, point, ID)
	return err
}

func (r *customerRepository) GetAll(ctx context.Context) ([]*model.Customer, error) {
	query := `
		SELECT id, name, point, created_at
		FROM customers
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []*model.Customer
	for rows.Next() {
		var c model.Customer
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Point,
			&c.CreatedAt,
		); err != nil {
			return nil, err
		}
		customers = append(customers, &c)
	}

	return customers, nil
}
