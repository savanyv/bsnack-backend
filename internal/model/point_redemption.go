package model

import (
	"time"

	"github.com/google/uuid"
)

type PointRedemption struct {
	ID uuid.UUID `db:"id" json:"id"`
	CustomerID uuid.UUID `db:"customer_id" json:"customer_id"`
	ProductID uuid.UUID `db:"product_id" json:"product_id"`
	PointUsed int `db:"point_used" json:"point_used"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
