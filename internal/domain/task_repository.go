package domain

import (
	"context"
)

type TaskRepository interface {
	Create(ctx context.Context, task *Task) (*Task, error)
	GetByID(ctx context.Context, id string) (*Task, error)
	List(ctx context.Context, filter TaskFilter) ([]*Task, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id string) error
}
