package dtos

type ReedemPointRequest struct {
	CustomerID string `json:"customer_id"`
	ProductID string `json:"product_id"`
	PointRequired int `json:"point_required"`
}
