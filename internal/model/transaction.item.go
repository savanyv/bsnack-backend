package model

import "github.com/google/uuid"

type TransactionItems struct {
	ID uuid.UUID `db:"id" json:"id"`
	TransactionID uuid.UUID `db:"transaction_id" json:"transaction_id"`
	ProductID uuid.UUID `db:"product_id" json:"product_id"`
	Quantity int `db:"quantity" json:"quantity"`
	Price int `db:"price" json:"price"`
}
