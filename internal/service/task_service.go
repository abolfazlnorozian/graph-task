package service

import (
	"context"
	"errors"
	"graph-task-service/internal/domain"
)

var (
	ErrTaskNotFound  = errors.New("task not found")
	ErrInvalidStatus = errors.New("invalid task status")
	ErrEmptyTitle    = errors.New("title cannot be empty")
)

type TaskService interface {
	CreateTask(ctx context.Context, title string, assignee *string, status *domain.TaskStatus) (*domain.Task, error)
	GetTask(ctx context.Context, id string) (*domain.Task, error)
	ListTasks(ctx context.Context, filter domain.TaskFilter) ([]*domain.Task, error)
	UpdateStatus(ctx context.Context, id string, status domain.TaskStatus) (*domain.Task, error)
	DeleteTask(ctx context.Context, id string) error
}

type taskService struct {
	repo domain.TaskRepository
}

func NewTaskService(repo domain.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

func (s *taskService) CreateTask(
	ctx context.Context,
	title string,
	assignee *string,
	status *domain.TaskStatus,
) (*domain.Task, error) {

	if title == "" {
		return nil, ErrEmptyTitle
	}

	finalStatus := domain.StatusTodo
	if status != nil {
		if !domain.IsValidStatus(*status) {
			return nil, ErrInvalidStatus
		}
		finalStatus = *status
	}

	task := &domain.Task{
		Title:    title,
		Status:   finalStatus,
		Assignee: assignee,
	}

	return s.repo.Create(ctx, task)
}

func (s *taskService) UpdateStatus(
	ctx context.Context,
	id string,
	status domain.TaskStatus,
) (*domain.Task, error) {

	if !domain.IsValidStatus(status) {
		return nil, ErrInvalidStatus
	}

	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	task.Status = status

	if err := s.repo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *taskService) DeleteTask(
	ctx context.Context,
	id string,
) error {

	return s.repo.Delete(ctx, id)
}

func (s *taskService) GetTask(
	ctx context.Context,
	id string,
) (*domain.Task, error) {

	return s.repo.GetByID(ctx, id)
}

func (s *taskService) ListTasks(
	ctx context.Context,
	filter domain.TaskFilter,
) ([]*domain.Task, error) {

	// filter := domain.TaskFilter{
	// 	Status:   status,
	// 	Assignee: assignee,
	// 	Limit:    20,
	// 	Offset:   0,
	// }

	return s.repo.List(ctx, filter)
}
