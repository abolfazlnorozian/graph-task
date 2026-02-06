package http

import "graph-task-service/internal/domain"

type CreateRequest struct {
	Title    string  `json:"title" binding:"required"`
	Assignee *string `json:"assignee"`
}

type UpdateStatusRequest struct {
	Status domain.TaskStatus `json:"status" binding:"required"`
}
