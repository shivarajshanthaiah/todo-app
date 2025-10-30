package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shivarajshanthaiah/todo-app/internal/repo/entity"
	"github.com/shivarajshanthaiah/todo-app/internal/repo/interfaces"
	"github.com/shivarajshanthaiah/todo-app/pkg/globals"
)

type TaskRepo struct {
	dao *pgxpool.Pool
}

func NewTaskRepository(dao *pgxpool.Pool) interfaces.TaskRepoInterface {
	return &TaskRepo{
		dao: dao,
	}
}

func (r *TaskRepo) CreateTodo(ctx context.Context, task *entity.Task) error {
	query := `
		INSERT INTO tasks (user_id, title, description, priority, status, due_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := r.dao.QueryRow(
		ctx,
		query,
		task.UserID,
		task.Title,
		task.Description,
		task.Priority,
		task.Status,
		task.DueAt,
	).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (r *TaskRepo) ListAllTodos(ctx context.Context, userID string, statusFilter string, limit, offset int) ([]*entity.Task, int64, error) {
	baseQuery := `
		SELECT
			id, user_id, title, description, priority, status, created_at, updated_at, due_at
		FROM
			tasks
		WHERE
			user_id = $1
	`
	countQuery := `
		SELECT COUNT(*) FROM tasks WHERE user_id = $1
	`

	args := []interface{}{userID}
	argIndex := 2

	// Add status filter if applicable
	if strings.ToUpper(statusFilter) != "ALL" && statusFilter != "" {
		statusVal, ok := globals.TaskStatus[statusFilter]
		if !ok {
			return nil, 0, fmt.Errorf("invalid status filter: %s", statusFilter)
		}
		baseQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, statusVal)
		argIndex++
	}

	// Add pagination and ordering
	baseQuery += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.dao.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tasks []*entity.Task
	for rows.Next() {
		var (
			id, priority, status        sql.NullInt32
			userID, title, description  sql.NullString
			createdAt, updatedAt, dueAt sql.NullTime
		)

		if err := rows.Scan(
			&id,
			&userID,
			&title,
			&description,
			&priority,
			&status,
			&createdAt,
			&updatedAt,
			&dueAt,
		); err != nil {
			return nil, 0, err
		}

		task := &entity.Task{
			ID:          int64(id.Int32),
			UserID:      userID.String,
			Title:       title.String,
			Description: description.String,
			Priority:    int(priority.Int32),
			Status:      int(status.Int32),
			CreatedAt:   createdAt.Time,
			UpdatedAt:   updatedAt.Time,
		}
		if dueAt.Valid {
			task.DueAt = dueAt.Time
		}
		tasks = append(tasks, task)
	}

	// Fetch total count
	var totalCount int64
	err = r.dao.QueryRow(ctx, countQuery, args[:len(args)-2]...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	return tasks, totalCount, nil
}

func (r *TaskRepo) UpdateTodoByID(ctx context.Context, id int, updatedTask *entity.Task) error {
	query := `
		UPDATE tasks
		SET
			title = $1,
			description = $2,
			priority = $3,
			status = $4,
			due_at = $5,
			updated_at = now()
		WHERE id = $6
	`

	_, err := r.dao.Exec(
		ctx,
		query,
		updatedTask.Title,
		updatedTask.Description,
		updatedTask.Priority,
		updatedTask.Status,
		updatedTask.DueAt,
		id,
	)
	return err
}

// This should be a soft delete, for now im keeping it as DELET operation
func (r *TaskRepo) DeleteTodo(ctx context.Context, id int) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.dao.Exec(ctx, query, id)
	return err
}

func (r *TaskRepo) GetTodoByID(ctx context.Context, id int) (*entity.Task, error) {
	query := `
		SELECT
			id,
			user_id,
			title,
			description,
			priority,
			status,
			due_at,
			created_at,
			updated_at
		FROM tasks
		WHERE id = $1
	`

	var (
		taskID, priority, status    sql.NullInt32
		userID, title, description  sql.NullString
		dueAt, createdAt, updatedAt sql.NullTime
	)

	err := r.dao.QueryRow(ctx, query, id).Scan(
		&taskID,
		&userID,
		&title,
		&description,
		&priority,
		&status,
		&dueAt,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	task := &entity.Task{
		ID:          int64(taskID.Int32),
		UserID:      userID.String,
		Title:       title.String,
		Description: description.String,
		Priority:    int(priority.Int32),
		Status:      int(status.Int32),
		DueAt:       dueAt.Time,
		CreatedAt:   createdAt.Time,
		UpdatedAt:   updatedAt.Time,
	}

	return task, nil
}
