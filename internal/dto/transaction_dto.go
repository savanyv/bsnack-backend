package dtos

type CreateTransactionRequest struct {
	CustomerName string `json:"customer_name"`
	Items []TransactionItemRequest `json:"items"`
	RedeemPoint int `json:"redeem_point"`
}

type TransactionItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity int `json:"quantity"`
}

type CreateTransactionResponse struct {
	TransactionID string `json:"transaction_id"`
	TotalPrice int `json:"total_price"`
	PointEarned int `json:"point_earned"`
}
