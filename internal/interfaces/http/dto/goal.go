package dto

type GoalSuggestionRequest struct {
	DailyExpenses   int64 `json:"daily_expenses" example:"1000" binding:"required"`
	DailyIncome     int64 `json:"daily_income" example:"2000" binding:"required"`
	WeeklyExpenses  int64 `json:"weekly_expenses" example:"7000" binding:"required"`
	WeeklyIncome    int64 `json:"weekly_income" example:"14000" binding:"required"`
	MonthlyExpenses int64 `json:"monthly_expenses" example:"30000" binding:"required"`
	MonthlyIncome   int64 `json:"monthly_income" example:"50000" binding:"required"`
}

type GoalSuggestionResponse struct {
	SuggestedAmount int64  `json:"suggested_amount" example:"15000"`
	Period          int    `json:"period" example:"30"`
	Message         string `json:"message" example:"Based on your income and expense analysis, we recommend that you can save 15,000 yuan per month."`
}

type CreateGoalRequest struct {
	TargetAmount int64 `json:"target_amount" example:"10000" binding:"required"`
	Period       int   `json:"period" example:"30" binding:"required"`
}

type GoalResponse struct {
	TargetAmount  int64  `json:"target_amount" example:"10000"`
	CurrentAmount int64  `json:"current_amount" example:"5000"`
	Period        int    `json:"period" example:"30"`
	Status        string `json:"status" example:"Need to work harder"`
	CreatedAt     string `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt     string `json:"updated_at" example:"2023-06-01T00:00:00Z"`
}

type GetGoalResponse struct {
	Goal GoalResponse `json:"goal"`
}
