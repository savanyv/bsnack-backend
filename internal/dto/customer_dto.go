package dtos

import "time"

type CustomerQuery struct {
	Month time.Time
}

type CustomerResponse struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Point int `json:"point"`
	IsNew bool `json:"is_new"`
	CreatedAt string `json:"created_at"`
}
