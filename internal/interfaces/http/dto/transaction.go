package dto

type TransactionResponse struct {
	Amount      int    `json:"amount" example:"1000" binding:"required"`
	Category    string `json:"category" example:"Food" binding:"required"`
	Type        string `json:"transaction_type" example:"Expense" binding:"required"`
	Date        string `json:"date" example:"2023-01-01" binding:"required"`
	Description string `json:"description" example:"Lunch" binding:"required"`
	CreatedAt   string `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   string `json:"updated_at" example:"2023-06-01T00:00:00Z"`
}

type GetTransactionsResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
}

type CreateTransactionRequest struct {
	Amount      int    `json:"amount" example:"1000" binding:"required"`
	Category    string `json:"category" example:"Food" binding:"required"`
	Type        string `json:"transaction_type" example:"Expense" binding:"required"`
	Date        string `json:"date" example:"2023-01-01" binding:"required"`
	Description string `json:"description" example:"Lunch" binding:"required"`
}
