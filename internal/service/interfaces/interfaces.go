package interfaces

import (
	"context"

	"github.com/shivarajshanthaiah/todo-app/internal/models"
)

type TaskServiceInterface interface {
	CreateTodoSvc(ctx context.Context, todo *models.Todo) error
	GetTodoByUserIDSvc(ctx context.Context, userID, status string, limit, offset int) (*models.PaginatedTodos, error)
	UpdateTodoByIDSvc(ctx context.Context, todo *models.Todo) error
	DeleteTodoByIDSvc(ctx context.Context, taskID int, userID string) error
}

type UserServiceInterface interface {
	UserSignUpSvc(ctx context.Context, user *models.User) error
	UserLoginSvc(ctx context.Context, login *models.Login) (string, error)
}
