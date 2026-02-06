package main

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
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system envs")
	}

	cfg := config.Load()

	if cfg.DBURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := postgres.NewDB(cfg.DBURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := postgres.RunMigrations(db); err != nil {
		log.Fatal(err)
	}

	taskRepo := postgres.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := http.NewTaskHandler(taskService)

	r := router.New(taskHandler)

	log.Printf("server running on :%s\n", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
