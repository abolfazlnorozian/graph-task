package postgres_test

import (
	"context"
	"database/sql"
	"graph-task-service/internal/domain"
	"graph-task-service/internal/repository/postgres"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib" // for testing

	"github.com/stretchr/testify/require"
)

var testRepo domain.TaskRepository
var testDB *sql.DB

func TestMain(m *testing.M) {
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/task_manager?sslmode=disable"
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}

	testDB = db
	testRepo = postgres.NewTaskRepository(db)

	code := m.Run()

	_ = db.Close()
	os.Exit(code)
}

func truncateTasks(t *testing.T) {
	_, err := testDB.Exec(`TRUNCATE TABLE tasks RESTART IDENTITY`)
	require.NoError(t, err)
}

func TestTaskRepository_Create_And_Get(t *testing.T) {
	truncateTasks(t)

	ctx := context.Background()

	task := &domain.Task{
		Title:  "integration test",
		Status: "pending",
	}

	created, err := testRepo.Create(ctx, task)
	require.NoError(t, err)
	require.NotZero(t, created.ID)

	found, err := testRepo.GetByID(ctx, created.ID)
	require.NoError(t, err)

	require.Equal(t, created.Title, found.Title)
	require.Equal(t, created.Status, found.Status)
}

func TestTaskRepository_GetByID_NotFound(t *testing.T) {
	truncateTasks(t)

	ctx := context.Background()

	task, err := testRepo.GetByID(ctx, "999")

	require.Nil(t, task)
	require.Error(t, err)
}
