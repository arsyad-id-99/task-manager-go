package repository

import (
	"context"

	"github.com/arsyad-id-99/task-manager-go/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) FindAllByUser(ctx context.Context, userID string) ([]model.Task, error) {
	query := `
        SELECT id, user_id, title, description, status, created_at, updated_at
        FROM tasks WHERE user_id = $1
        ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Description,
			&t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *TaskRepository) FindByID(ctx context.Context, id, userID string) (*model.Task, error) {
	task := &model.Task{}
	query := `
        SELECT id, user_id, title, description, status, created_at, updated_at
        FROM tasks WHERE id = $1 AND user_id = $2`
	err := r.db.QueryRow(ctx, query, id, userID).
		Scan(&task.ID, &task.UserID, &task.Title, &task.Description,
			&task.Status, &task.CreatedAt, &task.UpdatedAt)
	return task, err
}

func (r *TaskRepository) Create(ctx context.Context, task *model.Task) error {
	query := `
        INSERT INTO tasks (user_id, title, description)
        VALUES ($1, $2, $3)
        RETURNING id, status, created_at, updated_at`
	return r.db.QueryRow(ctx, query, task.UserID, task.Title, task.Description).
		Scan(&task.ID, &task.Status, &task.CreatedAt, &task.UpdatedAt)
}

func (r *TaskRepository) UpdateStatus(ctx context.Context, id, userID string, status model.TaskStatus) error {
	query := `
        UPDATE tasks SET status = $1, updated_at = NOW()
        WHERE id = $2 AND user_id = $3`
	result, err := r.db.Exec(ctx, query, status, id, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
