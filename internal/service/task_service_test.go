package service

import (
	"context"
	"graph-task-service/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTaskRepo struct {
	mock.Mock
}

func (m *mockTaskRepo) Create(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	args := m.Called(ctx, task)
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *mockTaskRepo) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *mockTaskRepo) List(
	ctx context.Context,
	filter domain.TaskFilter,
) ([]*domain.Task, error) {
	args := m.Called(ctx, filter)

	return args.Get(0).([]*domain.Task), args.Error(1)
}

func (m *mockTaskRepo) Update(ctx context.Context, task *domain.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *mockTaskRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateTask_Success(t *testing.T) {
	repo := new(mockTaskRepo)
	svc := NewTaskService(repo)

	repo.On(
		"Create",
		mock.Anything,
		mock.MatchedBy(func(t *domain.Task) bool {
			return t.Title == "test" &&
				t.Status == domain.StatusTodo
		}),
	).Return(&domain.Task{Title: "test"}, nil)

	task, err := svc.CreateTask(context.Background(), "test", nil)

	assert.NoError(t, err)
	assert.Equal(t, "test", task.Title)
	repo.AssertExpectations(t)
}
