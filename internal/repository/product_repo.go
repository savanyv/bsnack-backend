package repository

import (
	"context"
	"database/sql"

	"github.com/savanyv/bsnack-backend/internal/model"
)

type ProductRepository interface {
	FindByID(ctx context.Context, ID string) (*model.Product, error)
	Create(ctx context.Context, product *model.Product) error
	UpdateStock(ctx context.Context, tx *sql.Tx, ID string, stock int) error
	GetAll(ctx context.Context) ([]model.Product, error)
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) FindByID(ctx context.Context, ID string) (*model.Product, error) {
	var product model.Product
	query := `
		SELECT id, name, type, flavor, size, price, stock, created_at
		FROM products
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, ID).Scan(
		&product.ID,
		&product.Name,
		&product.Type,
		&product.Flavor,
		&product.Size,
		&product.Price,
		&product.Stock,
		&product.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) Create(ctx context.Context, product *model.Product) error {
	query := `
		INSERT INTO products (name, type, flavor, size, price, stock)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		product.Name,
		product.Type,
		product.Flavor,
		product.Size,
		product.Price,
		product.Stock,
	).Scan(
		&product.ID,
		&product.CreatedAt,
	)
}

func (r *productRepository) UpdateStock(ctx context.Context, tx *sql.Tx, ID string, stock int) error {
	query := `
		UPDATE products
		SET stock = $1
		WHERE id = $2
	`

	_, err := tx.ExecContext(ctx, query, stock, ID)
	return err
}

func (r *productRepository) GetAll(ctx context.Context) ([]model.Product, error) {
	query := `
		SELECT id, name, type, flavor, size, price, stock, created_at
		FROM products
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Type,
			&p.Flavor,
			&p.Size,
			&p.Price,
			&p.Stock,
			&p.CreatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}
