package router

import (
	"graph-task-service/internal/handler/http"

	"github.com/gin-gonic/gin"
)

func New(
	taskHandler *http.TaskHandler,
) *gin.Engine {

	r := gin.New()

	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	tasks := r.Group("/tasks")
	{
		tasks.POST("", taskHandler.Create)
		tasks.GET("", taskHandler.List)

		tasks.GET("/:id", taskHandler.GetByID)
		tasks.PATCH("/:id/status", taskHandler.UpdateStatus)
		tasks.DELETE("/:id", taskHandler.Delete)
	}

	return r
}
