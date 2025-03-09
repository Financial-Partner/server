package dto

type LoginRequest struct {
	FirebaseToken string `json:"firebase_token" binding:"required" example:"eyJhbGciOiJSUzI1NiIsImtpZCI6I..."`
}

type LoginResponse struct {
	AccessToken  string       `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6..."`
	RefreshToken string       `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6..."`
	ExpiresIn    int          `json:"expires_in" example:"3600"`
	TokenType    string       `json:"token_type" example:"Bearer"`
	User         UserResponse `json:"user"`
}

type UserResponse struct {
	ID        string `json:"id" example:"60d6ec33f777b123e4567890"`
	Email     string `json:"email" example:"user@example.com"`
	Name      string `json:"name" example:"User Name"`
	Diamonds  int64  `json:"diamonds" example:"100"`
	CreatedAt string `json:"created_at" example:"2025-03-07T12:00:00Z"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6..."`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6..."`
	ExpiresIn    int    `json:"expires_in" example:"3600"`
	TokenType    string `json:"token_type" example:"Bearer"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6..."`
}

type LogoutResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Successfully logged out"`
}
