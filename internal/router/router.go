package router

import (
	"graph-task-service/internal/handler/http"
	"os"

	"github.com/gin-gonic/gin"

	_ "graph-task-service/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	if os.Getenv("ENABLE_SWAGGER") == "true" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return r
}
