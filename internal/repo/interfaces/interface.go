package interfaces

import (
	"context"

	"github.com/shivarajshanthaiah/todo-app/internal/repo/entity"
)

type UserRepoInterface interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByID(ctx context.Context, ID string) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
}

type TaskRepoInterface interface {
	CreateTodo(ctx context.Context, task *entity.Task) error
	ListAllTodos(ctx context.Context, userID string, statusFilter string, limit, offset int) ([]*entity.Task, int64, error)
	UpdateTodoByID(ctx context.Context, id int, updatedTask *entity.Task) error
	DeleteTodo(ctx context.Context, id int) error
	GetTodoByID(ctx context.Context, id int) (*entity.Task, error)
}
