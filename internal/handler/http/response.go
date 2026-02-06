package http

import (
	"graph-task-service/internal/domain"
	"time"
)

type TaskResponse struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	Status    string  `json:"status"`
	Assignee  *string `json:"assignee,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

func FromDomain(t *domain.Task) TaskResponse {
	return TaskResponse{
		ID:        t.ID,
		Title:     t.Title,
		Status:    string(t.Status),
		Assignee:  t.Assignee,
		CreatedAt: t.CreatedAt.Format(time.RFC3339),
		UpdatedAt: t.UpdatedAt.Format(time.RFC3339),
	}
}
