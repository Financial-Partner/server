package dto

type InvestmentResponse struct {
	Title       string   `json:"title" example:"Investment in stock market"`
	Description string   `json:"description" example:"Investment in stock market is a good way to make money"`
	Tags        []string `json:"tags" example:"stock, market"`
	IsIncrease  bool     `json:"is_increase" example:"true"`
	Status      int64    `json:"status" example:"20"`
	CreatedAt   string   `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   string   `json:"updated_at" example:"2023-06-01T00:00:00Z"`
}

type GetInvestmentsResponse struct {
	Investments []InvestmentResponse `json:"investments"`
}
