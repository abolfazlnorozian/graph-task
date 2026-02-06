package server

import (
	"graph-task-service/internal/config"
	"graph-task-service/internal/handler/http"
	"graph-task-service/internal/repository/postgres"
	"graph-task-service/internal/router"
	"graph-task-service/internal/service"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// 1. load .env for local (if exists)
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system envs")
	}

	// 2. load config
	cfg := config.Load()

	if cfg.DBURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	// 3. init database
	db, err := postgres.NewDB(cfg.DBURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 4. init layers
	taskRepo := postgres.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := http.NewTaskHandler(taskService)

	// 5. init router
	r := router.New(taskHandler)

	// 6. run server
	log.Printf("server running on :%s\n", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
