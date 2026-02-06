package http

import (
	"errors"
	"graph-task-service/internal/domain"
	"graph-task-service/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	service service.TaskService
}

func NewTaskHandler(s service.TaskService) *TaskHandler {
	return &TaskHandler{service: s}
}

// CreateTask godoc
// @Summary      Create a task
// @Description  Create a new task in the system
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        request body http.CreateRequest true "Create task payload"
// @Success      201 {object} http.TaskResponse "Task created successfully"
// @Failure      400 {object} http.ErrorResponse "Invalid request body"
// @Failure      500 {object} http.ErrorResponse "Internal server error"
// @Router       /tasks [post]
func (h *TaskHandler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.service.CreateTask(
		c.Request.Context(),
		req.Title,
		req.Assignee,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, FromDomain(task))
}

// List godoc
// @Summary      List tasks
// @Description  List tasks with optional filtering and pagination
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        status    query     string  false  "Task status" Enums(todo,in_progress,done)
// @Param        assignee  query     string  false  "Assignee username"
// @Param        limit     query     int     false  "Limit"   default(20)
// @Param        offset    query     int     false  "Offset"  default(0)
// @Success      200  {array}   http.TaskResponse
// @Failure      400  {object}  http.ErrorResponse
// @Failure      500  {object}  http.ErrorResponse
// @Router       /tasks [get]
func (h *TaskHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	var (
		status   *domain.TaskStatus
		assignee *string
	)

	if s := c.Query("status"); s != "" {
		st := domain.TaskStatus(s)

		if !domain.IsValidStatus(st) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
			return
		}
		status = &st
	}

	if a := c.Query("assignee"); a != "" {
		assignee = &a
	}

	filter := domain.TaskFilter{
		Status:   status,
		Assignee: assignee,
		Limit:    limit,
		Offset:   offset,
	}

	tasks, err := h.service.ListTasks(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]TaskResponse, 0, len(tasks))
	for _, t := range tasks {
		resp = append(resp, FromDomain(t))
	}

	c.JSON(http.StatusOK, resp)
}

// GetByID godoc
// @Summary      Get task by ID
// @Description  Retrieve a task by its ID
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Task ID"
// @Success      200  {object}  http.TaskResponse
// @Failure      400  {object}  http.ErrorResponse
// @Failure      404  {object}  http.ErrorResponse
// @Failure      500  {object}  http.ErrorResponse
// @Router       /tasks/{id} [get]
func (h *TaskHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	task, err := h.service.GetTask(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, FromDomain(task))
}

// UpdateStatus godoc
// @Summary      Update task status
// @Description  Update status of a task
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id       path      string                     true  "Task ID"
// @Param        request  body      http.UpdateStatusRequest   true  "Update status payload"
// @Success      200      {object}  http.TaskResponse
// @Failure      400      {object}  http.ErrorResponse
// @Failure      404      {object}  http.ErrorResponse
// @Failure      500      {object}  http.ErrorResponse
// @Router       /tasks/{id}/status [patch]
func (h *TaskHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	task, err := h.service.UpdateStatus(
		c.Request.Context(),
		id,
		req.Status,
	)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, FromDomain(task))
}

// Delete godoc
// @Summary      Delete task
// @Description  Delete a task by ID
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id   path  string  true  "Task ID"
// @Success      204  "No Content"
// @Failure      400  {object}  http.ErrorResponse
// @Failure      404  {object}  http.ErrorResponse
// @Failure      500  {object}  http.ErrorResponse
// @Router       /tasks/{id} [delete]
func (h *TaskHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	err := h.service.DeleteTask(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
