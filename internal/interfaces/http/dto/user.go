package dto

type CreateUserRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
	Name  string `json:"name" binding:"required" example:"New User Name"`
}

type CreateUserResponse struct {
	ID        string `json:"id" example:"60d6ec33f777b123e4567890"`
	Email     string `json:"email" example:"user@example.com"`
	Name      string `json:"name" example:"User Name"`
	Diamonds  int64  `json:"diamonds" example:"100"`
	Savings   int64  `json:"savings" example:"0"`
	CreatedAt string `json:"created_at" example:"2025-03-07T12:00:00Z"`
}

type UpdateUserRequest struct {
	Name string `json:"name" binding:"required" example:"New User Name"`
}

type UpdateUserCharacterRequest struct {
	CharacterID string `json:"character_id" binding:"required" example:"char_001"`
	ImageURL    string `json:"image_url" binding:"required" example:"https://example.com/characters/advisor.png"`
}

type UpdateUserResponse struct {
	ID        string `json:"id" example:"60d6ec33f777b123e4567890"`
	Email     string `json:"email" example:"user@example.com"`
	Name      string `json:"name" example:"New User Name"`
	Diamonds  int64  `json:"diamonds" example:"100"`
	Savings   int64  `json:"savings" example:"5000"`
	UpdatedAt string `json:"updated_at" example:"2025-03-07T12:00:00Z"`
}

type GetUserResponse struct {
	ID        string             `json:"id" example:"60d6ec33f777b123e4567890"`
	Email     *string            `json:"email,omitempty" example:"user@example.com"`
	Name      *string            `json:"name,omitempty" example:"User Name"`
	Wallet    *WalletResponse    `json:"wallet,omitempty"`
	Character *CharacterResponse `json:"character,omitempty"`
	CreatedAt string             `json:"created_at" example:"2025-03-07T12:00:00Z"`
	UpdatedAt string             `json:"updated_at" example:"2025-03-07T12:00:00Z"`
}

type WalletResponse struct {
	Diamonds int64 `json:"diamonds" example:"100"`
	Savings  int64 `json:"savings" example:"5000"`
}

type CharacterResponse struct {
	ID       string `json:"id" example:"char_001"`
	Name     string `json:"name" example:"Character Name"`
	ImageURL string `json:"image_url" example:"https://example.com/characters/advisor.png"`
}
