package entities

import (
	"time"
)

type Investment struct {
	ID          string
	UserID      string
	Title       string
	Description string
	Tags        []string
	IsIncrease  bool
	Status      int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
