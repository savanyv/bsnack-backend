package model

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID uuid.UUID `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Point int `db:"point" json:"point"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
