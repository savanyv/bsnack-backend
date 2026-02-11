package dtos

import "time"

type CreateProductRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Flavor string `json:"flavor"`
	Size string `json:"size"`
	Price int `json:"price"`
	Stock int `json:"stock"`
}

type ProductResponse struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Flavor string `json:"flavor"`
	Size string `json:"size"`
	Price int `json:"price"`
	Stock int `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
}
