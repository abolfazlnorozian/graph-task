package postgres

import (
	"database/sql"
	"log"
	"os"
)

func RunMigrations(db *sql.DB) error {
	log.Println("RUN_MIGRATIONS:", os.Getenv("RUN_MIGRATIONS"))
	if os.Getenv("RUN_MIGRATIONS") == "true" {
		log.Println("running database migrations...")

		schema := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	CREATE TABLE IF NOT EXISTS tasks (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		title TEXT NOT NULL,
		status TEXT NOT NULL,
		assignee TEXT,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
	CREATE INDEX IF NOT EXISTS idx_tasks_assignee ON tasks(assignee);
	`

		_, err := db.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil
}
