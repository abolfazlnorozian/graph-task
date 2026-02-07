package service_test

import (
	"context"
	"graph-task-service/internal/domain"
	"graph-task-service/internal/service"
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
	svc := service.NewTaskService(repo)

	repo.On(
		"Create",
		mock.Anything,
		mock.MatchedBy(func(t *domain.Task) bool {
			return t.Title == "test" &&
				t.Status == domain.StatusTodo
		}),
	).Return(&domain.Task{Title: "test"}, nil)

	task, err := svc.CreateTask(context.Background(), "test", nil, nil)

	assert.NoError(t, err)
	assert.Equal(t, "test", task.Title)
	repo.AssertExpectations(t)
}

func TestGetTask_Success(t *testing.T) {
	repo := new(mockTaskRepo)
	svc := service.NewTaskService(repo)

	expected := &domain.Task{ID: "1", Title: "test"}

	repo.On(
		"GetByID",
		mock.Anything,
		"1",
	).Return(expected, nil)

	task, err := svc.GetTask(context.Background(), "1")

	assert.NoError(t, err)
	assert.Equal(t, expected, task)
	repo.AssertExpectations(t)
}

func TestGetTask_NotFound(t *testing.T) {
	repo := new(mockTaskRepo)
	svc := service.NewTaskService(repo)

	repo.On(
		"GetByID",
		mock.Anything,
		"404",
	).Return((*domain.Task)(nil), service.ErrTaskNotFound)

	task, err := svc.GetTask(context.Background(), "404")

	assert.Nil(t, task)
	assert.ErrorIs(t, err, service.ErrTaskNotFound)
	repo.AssertExpectations(t)
}

func TestListTasks_Success(t *testing.T) {
	repo := new(mockTaskRepo)
	svc := service.NewTaskService(repo)

	filter := domain.TaskFilter{}

	expected := []*domain.Task{
		{ID: "1", Title: "task 1"},
		{ID: "2", Title: "task 2"},
	}

	repo.On(
		"List",
		mock.Anything,
		filter,
	).Return(expected, nil)

	tasks, err := svc.ListTasks(context.Background(), filter)

	assert.NoError(t, err)
	assert.Len(t, tasks, 2)
	assert.Equal(t, expected, tasks)
	repo.AssertExpectations(t)
}

func TestUpdateStatus_Success(t *testing.T) {
	repo := new(mockTaskRepo)
	svc := service.NewTaskService(repo)

	task := &domain.Task{
		ID:     "1",
		Title:  "test",
		Status: domain.StatusTodo,
	}

	repo.On(
		"GetByID",
		mock.Anything,
		"1",
	).Return(task, nil)

	repo.On(
		"Update",
		mock.Anything,
		mock.MatchedBy(func(t *domain.Task) bool {
			return t.Status == domain.StatusDone
		}),
	).Return(nil)

	updated, err := svc.UpdateStatus(
		context.Background(),
		"1",
		domain.StatusDone,
	)

	assert.NoError(t, err)
	assert.Equal(t, domain.StatusDone, updated.Status)
	repo.AssertExpectations(t)
}

func TestUpdateStatus_InvalidStatus(t *testing.T) {
	repo := new(mockTaskRepo)
	svc := service.NewTaskService(repo)

	task, err := svc.UpdateStatus(
		context.Background(),
		"1",
		"invalid",
	)

	assert.Nil(t, task)
	assert.ErrorIs(t, err, service.ErrInvalidStatus)

	repo.AssertExpectations(t)
}

func TestDeleteTask_Success(t *testing.T) {
	repo := new(mockTaskRepo)
	svc := service.NewTaskService(repo)

	repo.On(
		"Delete",
		mock.Anything,
		"1",
	).Return(nil)

	err := svc.DeleteTask(context.Background(), "1")

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
