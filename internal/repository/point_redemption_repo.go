package repository

import (
	"context"
	"database/sql"

	"github.com/savanyv/bsnack-backend/internal/model"
)

type PointRedemptionRepository interface {
	Create(ctx context.Context, tx *sql.Tx, redemption *model.PointRedemption) error
}

type pointRedemptionRepository struct {
	db *sql.DB
}

func NewPointRedemptionRepository(db *sql.DB) PointRedemptionRepository {
	return &pointRedemptionRepository{
		db: db,
	}
}

func (r *pointRedemptionRepository) Create(ctx context.Context, tx *sql.Tx, redemption *model.PointRedemption) error {
	query := `
		INSERT INTO point_redemptions (customer_id, product_id, point_used)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := tx.QueryRowContext(ctx, query, redemption.CustomerID, redemption.ProductID, redemption.PointUsed).Scan(
		&redemption.ID,
		&redemption.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}
