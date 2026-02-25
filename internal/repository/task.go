package repository

import (
	"context"
	"database/sql"
	"time"

	"ApiGolangv2/internal/entity"
)

type TaskRepository interface {
	FindAll(ctx context.Context) ([]entity.Task, error)
	FindByID(ctx context.Context, id int) (*entity.Task, error)
	FindByStatus(ctx context.Context, status string) ([]entity.Task, error)
	Create(ctx context.Context, task *entity.Task) (*entity.Task, error)
	Update(ctx context.Context, task *entity.Task) (*entity.Task, error)
	Delete(ctx context.Context, id int) error
}

type taskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) FindAll(ctx context.Context) ([]entity.Task, error) {
	query := `
		SELECT id_task, title, description, status, priority, due_date, done, created_at, updated_at
		FROM tasks
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []entity.Task
	for rows.Next() {
		var task entity.Task
		err := rows.Scan(
			&task.IDTask,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.DueDate,
			&task.Done,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *taskRepository) FindByID(ctx context.Context, id int) (*entity.Task, error) {
	query := `
		SELECT id_task, title, description, status, priority, due_date, done, created_at, updated_at
		FROM tasks
		WHERE id_task = $1
	`

	var task entity.Task
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&task.IDTask,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.DueDate,
		&task.Done,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *taskRepository) FindByStatus(ctx context.Context, status string) ([]entity.Task, error) {
	query := `
		SELECT id_task, title, description, status, priority, due_date, done, created_at, updated_at
		FROM tasks
		WHERE status = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []entity.Task
	for rows.Next() {
		var task entity.Task
		err := rows.Scan(
			&task.IDTask,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.DueDate,
			&task.Done,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *taskRepository) Create(ctx context.Context, task *entity.Task) (*entity.Task, error) {
	query := `
		INSERT INTO tasks (title, description, status, priority, due_date, done, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id_task, title, description, status, priority, due_date, done, created_at, updated_at
	`

	now := time.Now()
	var created entity.Task
	err := r.db.QueryRowContext(ctx, query,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.DueDate,
		task.Done,
		now,
		now,
	).Scan(
		&created.IDTask,
		&created.Title,
		&created.Description,
		&created.Status,
		&created.Priority,
		&created.DueDate,
		&created.Done,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *taskRepository) Update(ctx context.Context, task *entity.Task) (*entity.Task, error) {
	query := `
		UPDATE tasks
		SET title       = $1,
		    description = $2,
		    status      = $3,
		    priority    = $4,
		    due_date    = $5,
		    done        = $6,
		    updated_at  = $7
		WHERE id_task = $8
		RETURNING id_task, title, description, status, priority, due_date, done, created_at, updated_at
	`

	var updated entity.Task
	err := r.db.QueryRowContext(ctx, query,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.DueDate,
		task.Done,
		time.Now(),
		task.IDTask,
	).Scan(
		&updated.IDTask,
		&updated.Title,
		&updated.Description,
		&updated.Status,
		&updated.Priority,
		&updated.DueDate,
		&updated.Done,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *taskRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM tasks WHERE id_task = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
