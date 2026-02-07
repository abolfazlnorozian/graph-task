package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"graph-task-service/internal/domain"
	handlerHttp "graph-task-service/internal/handler/http"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTaskService struct {
	mock.Mock
}

func (m *MockTaskService) CreateTask(
	ctx context.Context,
	title string,
	assignee *string,
	status *domain.TaskStatus,
) (*domain.Task, error) {

	args := m.Called(ctx, title, assignee, status)
	task, _ := args.Get(0).(*domain.Task)
	return task, args.Error(1)
}

func (m *MockTaskService) ListTasks(
	ctx context.Context,
	filter domain.TaskFilter,
) ([]*domain.Task, error) {

	args := m.Called(ctx, filter)
	tasks, _ := args.Get(0).([]*domain.Task)
	return tasks, args.Error(1)
}

func (m *MockTaskService) GetTask(
	ctx context.Context,
	id string,
) (*domain.Task, error) {

	args := m.Called(ctx, id)
	task, _ := args.Get(0).(*domain.Task)
	return task, args.Error(1)
}

func (m *MockTaskService) UpdateStatus(
	ctx context.Context,
	id string,
	status domain.TaskStatus,
) (*domain.Task, error) {

	args := m.Called(ctx, id, status)
	task, _ := args.Get(0).(*domain.Task)
	return task, args.Error(1)
}

func (m *MockTaskService) DeleteTask(
	ctx context.Context,
	id string,
) error {

	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupRouter(handler *handlerHttp.TaskHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.POST("/tasks", handler.Create)

	return r
}

func TestTaskHandler_Create_Success(t *testing.T) {
	service := new(MockTaskService)
	handler := handlerHttp.NewTaskHandler(service)
	router := setupRouter(handler)

	assignee := "abo"

	reqBody := map[string]interface{}{
		"title":    "test task",
		"assignee": assignee,
	}

	body, _ := json.Marshal(reqBody)

	expectedTask := &domain.Task{
		ID:       "c292d1f6-b03b-4490-a2cb-3bd272f05dda",
		Title:    "test task",
		Assignee: &assignee,
		Status:   domain.StatusTodo,
	}

	service.
		On(
			"CreateTask",
			mock.Anything,
			"test task",
			&assignee,
			(*domain.TaskStatus)(nil),
		).
		Return(expectedTask, nil)

	req := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), `"title":"test task"`)

	service.AssertExpectations(t)
}
