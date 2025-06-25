package dto

import (
	"go-io-bound-api/internal/task"
	"time"
)

type TaskResponse struct {
	ID          string      `json:"id"`
	Status      task.Status `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	StartedAt   time.Time   `json:"started_at"`
	CompletedAt time.Time   `json:"completed_at"`
	Duration    string      `json:"duration"`
}
