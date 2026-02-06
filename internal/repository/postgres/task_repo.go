package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"graph-task-service/internal/domain"
)

type taskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) domain.TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(
	ctx context.Context,
	task *domain.Task,
) (*domain.Task, error) {

	query := `
		INSERT INTO tasks (title, status, assignee)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		task.Title,
		task.Status,
		task.Assignee,
	).Scan(
		&task.ID,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return task, nil
}

func (r *taskRepository) GetByID(
	ctx context.Context,
	id string,
) (*domain.Task, error) {

	var task domain.Task

	query := `
		SELECT id, title, status, assignee, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Status,
		&task.Assignee,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *taskRepository) List(
	ctx context.Context,
	filter domain.TaskFilter,
) ([]*domain.Task, error) {

	query := `
		SELECT id, title, status, assignee, created_at, updated_at
		FROM tasks
		WHERE 1=1
	`
	args := []any{}
	argID := 1

	if filter.Status != nil {
		query += " AND status = $" + fmt.Sprint(argID)
		args = append(args, *filter.Status)
		argID++
	}

	if filter.Assignee != nil {
		query += " AND assignee = $" + fmt.Sprint(argID)
		args = append(args, *filter.Assignee)
		argID++
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += " LIMIT $" + fmt.Sprint(argID)
		args = append(args, filter.Limit)
		argID++
	}

	if filter.Offset > 0 {
		query += " OFFSET $" + fmt.Sprint(argID)
		args = append(args, filter.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*domain.Task

	for rows.Next() {
		var t domain.Task
		if err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Status,
			&t.Assignee,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}

	return tasks, nil
}

func (r *taskRepository) Update(
	ctx context.Context,
	task *domain.Task,
) error {

	query := `
		UPDATE tasks
		SET title = $1,
		    status = $2,
		    assignee = $3,
		    updated_at = now()
		WHERE id = $4
	`

	res, err := r.db.ExecContext(
		ctx,
		query,
		task.Title,
		task.Status,
		task.Assignee,
		task.ID,
	)

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *taskRepository) Delete(
	ctx context.Context,
	id string,
) error {

	res, err := r.db.ExecContext(
		ctx,
		`DELETE FROM tasks WHERE id = $1`,
		id,
	)

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
