package service

import (
	"context"
	"errors"
	"log"

	"github.com/shivarajshanthaiah/todo-app/internal/models"
	"github.com/shivarajshanthaiah/todo-app/internal/repo/entity"
	repo "github.com/shivarajshanthaiah/todo-app/internal/repo/interfaces"
	service "github.com/shivarajshanthaiah/todo-app/internal/service/interfaces"
	"github.com/shivarajshanthaiah/todo-app/pkg/globals"
)

type TaskService struct {
	repo repo.TaskRepoInterface
}

func NewTaskService(repo repo.TaskRepoInterface) service.TaskServiceInterface {
	return &TaskService{
		repo: repo,
	}
}

func (s *TaskService) CreateTodoSvc(ctx context.Context, todo *models.Todo) error {
	// Convert Priority
	priorityVal, ok := globals.TaskPriority[todo.Priority]
	if !ok {
		log.Println("Invalid priority value:", todo.Priority)
		return errors.New("invalid task priority")
	}

	// Convert Status
	statusVal, ok := globals.TaskStatus[todo.Status]
	if !ok {
		log.Println("Invalid status value:", todo.Status)
		return errors.New("invalid task status")
	}

	// Map model to entity
	entityTask := entity.Task{
		UserID:      todo.UserID,
		Title:       todo.Title,
		Description: todo.Description,
		Priority:    priorityVal,
		Status:      statusVal,
		DueAt:       todo.DueAt,
	}

	err := s.repo.CreateTodo(ctx, &entityTask)
	if err != nil {
		log.Println("Error creating todo in repo:", err)
		return err
	}
	return nil
}

func (s *TaskService) GetTodoByUserIDSvc(ctx context.Context, userID, status string, limit, offset int) (*models.PaginatedTodos, error) {

	tasks, total, err := s.repo.ListAllTodos(ctx, userID, status, limit, offset)
	if err != nil {
		log.Println("Error fetching todos from repo:", err)
		return nil, err
	}

	var todos []*models.Todo
	for _, task := range tasks {
		todos = append(todos, &models.Todo{
			ID:          task.ID,
			UserID:      task.UserID,
			Title:       task.Title,
			Description: task.Description,
			Priority:    globals.TaskPriorityReverse[task.Priority],
			Status:      globals.TaskStatusReverse[task.Status],
			DueAt:       task.DueAt,
			Created:     task.CreatedAt,
			Updated:     task.UpdatedAt,
		})
	}

	return &models.PaginatedTodos{
		TotalCount: total,
		Todos:      todos,
	}, nil
}

func (s *TaskService) UpdateTodoByIDSvc(ctx context.Context, todo *models.Todo) error {
	// Convert Priority
	priorityVal, ok := globals.TaskPriority[todo.Priority]
	if !ok {
		log.Println("Invalid priority value:", todo.Priority)
		return errors.New("invalid task priority")
	}

	// Convert Status
	statusVal, ok := globals.TaskStatus[todo.Status]
	if !ok {
		log.Println("Invalid status value:", todo.Status)
		return errors.New("invalid task status")
	}

	// Map model to entity
	entityTask := &entity.Task{
		ID:          todo.ID,
		UserID:      todo.UserID,
		Title:       todo.Title,
		Description: todo.Description,
		Priority:    priorityVal,
		Status:      statusVal,
		DueAt:       todo.DueAt,
	}

	log.Println("modified task", entityTask)

	err := s.repo.UpdateTodoByID(ctx, int(todo.ID), entityTask)
	if err != nil {
		log.Println("Error updating todo in repo:", err)
		return err
	}

	return nil
}

// We can do soft delete by adding boolean column is_deleted instead of directly deleting
func (s *TaskService) DeleteTodoByIDSvc(ctx context.Context, taskID int, userID string) error {
	task, err := s.repo.GetTodoByID(ctx, taskID)
	if err != nil {
		log.Println("Error fetching todo from repo:", err)
		return err
	}

	log.Printf("Attempting to delete task %d for user ID (from param): %s", taskID, userID)

	log.Println(task)

	// Check ownership
	if task == nil || task.UserID != userID {
		return errors.New("task not found or unauthorized")
	}

	// Perform delete
	err = s.repo.DeleteTodo(ctx, taskID)
	if err != nil {
		log.Println("Error deleting todo in repo:", err)
		return err
	}
	return nil
}
