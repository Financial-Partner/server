package entities

import (
	"time"
)

type GoalSuggestion struct {
	SuggestedAmount int64
	Period          int
	Message         string
}

type Goal struct {
	ID            string
	UserID        string
	TargetAmount  int64
	CurrentAmount int64
	Period        int
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type GoalMilestone struct {
	Title         string
	TargetPercent int
	Reward        string
	IsCompleted   bool
	CompletedAt   *time.Time
}
