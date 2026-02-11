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

type GetTransactionQuery struct {
	StartDate string `query:"start_date" json:"start_date"`
	EndDate string `query:"end_date" json:"end_date"`
}

type TransactionSummaryResponse struct {
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
	TotalCustomer int `json:"total_customer"`
	TotalIncome int `json:"total_income"`
	BestSeller string `json:"best_seller"`
	Transactions []TransactionListItem `json:"transactions"`
}

type TransactionListItem struct {
	ID string `json:"id"`
	CustomerID string `json:"customer_id"`
	TotalPrice int `json:"total_price"`
	CreatedAt string `json:"created_at"`
}
