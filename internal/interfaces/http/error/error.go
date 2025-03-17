package httperror

const (
	ErrInvalidRequest            = "Invalid request format"
	ErrUnauthorized              = "Unauthorized"
	ErrEmailNotFound             = "Email not found"
	ErrUserIDNotFound            = "User ID not found"
	ErrEmailMismatch             = "Email mismatch"
	ErrInternalServer            = "Internal server error"
	ErrUserNotFound              = "User not found"
	ErrInvalidRefreshToken       = "Invalid refresh token"
	ErrFailedToCreateUser        = "Failed to create user"
	ErrFailedToUpdateUser        = "Failed to update user"
	ErrFailedToGetUser           = "Failed to get user"
	ErrFailedToLogout            = "Failed to logout"
	ErrFailedToGetGoalSuggestion = "Failed to get goal suggestion"
	ErrFailedToCreateGoal        = "Failed to create goal"
	ErrFailedToUpdateGoal        = "Failed to update goal"
	ErrFailedToGetGoal           = "Failed to get goal"
	ErrFailedToGetInvestments    = "Failed to get investments"
)
