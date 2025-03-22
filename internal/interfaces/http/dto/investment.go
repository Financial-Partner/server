package dto

type OpportunityResponse struct {
	ID          string   `json:"id" example:"60d6ec33f777b123e4567890"`
	Title       string   `json:"title" example:"Investment in stock market"`
	Description string   `json:"description" example:"Investment in stock market is a good way to make money"`
	Tags        []string `json:"tags" example:"stock, market"`
	IsIncrease  bool     `json:"is_increase" example:"true"`
	Variation   int64    `json:"variation" example:"20"`
	Duration    string   `json:"duration" example:"a month"`
	MinAmount   int64    `json:"min_amount" example:"1000"`
	CreatedAt   string   `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   string   `json:"updated_at" example:"2023-06-01T00:00:00Z"`
}

type InvestmentResponse struct {
	ID            string `json:"id" example:"60d6ec33f777b123e4567890"`
	OpportunityID string `json:"opportunity_id" example:"60d6ec33f777b123e4567890"`
	UserID        string `json:"user_id" example:"60d6ec33f777b123e4567890"`
	Amount        int64  `json:"amount" example:"1000"`
	CreatedAt     string `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt     string `json:"updated_at" example:"2023-06-01T00:00:00Z"`
}

type GetOpportunitiesResponse struct {
	Opportunities []OpportunityResponse `json:"opportunities"`
}

type CreateUserInvestmentRequest struct {
	OpportunityID string `json:"opportunity_id" example:"60d6ec33f777b123e4567890" binding:"required"`
	Amount        int64  `json:"amount" example:"1000" binding:"required"`
}

type CreateUserInvestmentResponse struct {
	Investment InvestmentResponse `json:"investment"`
}

type GetUserInvestmentsResponse struct {
	Investments []InvestmentResponse `json:"investments"`
}
