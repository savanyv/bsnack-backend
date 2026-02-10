package model

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID uuid.UUID `db:"id" json:"id"`
	CustomerID uuid.UUID `db:"customer_id" json:"customer_id"`
	TransactionDate time.Time `db:"transaction_date" json:"transaction_date"`
	TotalPrice int `db:"total_price" json:"total_price"`
}
