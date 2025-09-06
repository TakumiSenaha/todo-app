package domain

import "time"

type Todo struct {
	ID          int
	UserID      int
	Title       string
	DueDate     *time.Time
	Priority    int
	IsCompleted bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}