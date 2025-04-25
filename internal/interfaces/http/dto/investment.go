package dto

type OpportunityResponse struct {
	OpportunityID string   `json:"opportunity_id" example:"60d6ec33f777b123e4567890"`
	Title         string   `json:"title" example:"Real Estate"`
	Description   string   `json:"description" example:"Investment in stock market is a good way to make money"`
	Tags          []string `json:"tags" example:"high risk,long term"`
	IsIncrease    bool     `json:"is_increase" example:"true"`
	Variation     int64    `json:"variation" example:"20"`
	Duration      string   `json:"duration" example:"a month"`
	MinAmount     int64    `json:"min_amount" example:"1000"`
	CreatedAt     string   `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt     string   `json:"updated_at" example:"2023-06-01T00:00:00Z"`
}

type InvestmentResponse struct {
	OpportunityID string `json:"opportunity_id" example:"60d6ec33f777b123e4567890"`
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

type CreateOpportunityRequest struct {
	Title       string   `json:"title" example:"Real Estate" binding:"required"`
	Description string   `json:"description" example:"Investment in stock market is a good way to make money" binding:"required"`
	Tags        []string `json:"tags" example:"high risk,long term" binding:"required"`
	IsIncrease  bool     `json:"is_increase" example:"true" binding:"required"`
	Variation   int64    `json:"variation" example:"20" binding:"required"`
	Duration    string   `json:"duration" example:"a month" binding:"required"`
	MinAmount   int64    `json:"min_amount" example:"1000" binding:"required"`
}

type CreateOpportunityResponse struct {
	Opportunity OpportunityResponse `json:"opportunity"`
}
