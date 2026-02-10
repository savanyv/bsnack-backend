package model

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID uuid.UUID `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Type string `db:"type" json:"type"`
	Flavor string `db:"flavor" json:"flavor"`
	Size string `db:"size" json:"size"`
	Price int `db:"price" json:"price"`
	Stock int `db:"stock" json:"stock"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
