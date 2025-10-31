package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shivarajshanthaiah/todo-app/internal/models"
	"github.com/shivarajshanthaiah/todo-app/internal/service/interfaces"
)

type TaskHandler struct {
	service interfaces.TaskServiceInterface
}

func NewTaskHandler(service interfaces.TaskServiceInterface) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) CreateTodoHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  http.StatusBadRequest,
			"Message": "Error binding request body",
			"Error":   err.Error(),
		})
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Status": http.StatusBadRequest,
			"Message": "error while user id from context",
			"Error":   ""})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Status": http.StatusBadRequest,
			"Message": "error while converting user id to string",
			"Error":   ""})
		return
	}

	todo.UserID = userIDStr

	if err := h.service.CreateTodoSvc(ctx, &todo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Status":  http.StatusInternalServerError,
			"Message": "Error creating todo",
			"Error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"Status":  http.StatusCreated,
		"Message": "Todo created successfully",
	})
}

func (h *TaskHandler) GetTodosHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	var req models.Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  http.StatusBadRequest,
			"Message": "Invalid request body",
			"Error":   err.Error(),
		})
		return
	}

	// Set defaults if not provided
	if req.Status == "" {
		req.Status = "ALL"
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Status": http.StatusBadRequest,
			"Message": "error while user id from context",
			"Error":   ""})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Status": http.StatusBadRequest,
			"Message": "error while converting user id to string",
			"Error":   ""})
		return
	}

	todos, err := h.service.GetTodoByUserIDSvc(ctx, userIDStr, req.Status, req.Limit, req.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Status":  http.StatusInternalServerError,
			"Message": "Error fetching todos",
			"Error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":  http.StatusOK,
		"Message": "Todos fetched successfully",
		"Data":    todos,
	})
}

func (h *TaskHandler) UpdateTodoHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	taskIDStr := c.Param("id")
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Status": http.StatusBadRequest,
			"Message": "error while user id from context",
			"Error":   ""})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Status": http.StatusBadRequest,
			"Message": "error while converting user id to string",
			"Error":   ""})
		return
	}

	id, err := strconv.Atoi(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "invalid task id"})
		return
	}

	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  http.StatusBadRequest,
			"Message": "Error binding request body",
			"Error":   err.Error(),
		})
		return
	}
	todo.UserID = userIDStr
	todo.ID = int64(id)

	if err := h.service.UpdateTodoByIDSvc(ctx, &todo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Status":  http.StatusInternalServerError,
			"Message": "Error updating todo",
			"Error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":  http.StatusOK,
		"Message": "Todo updated successfully",
	})
}

func (h *TaskHandler) DeleteTodoHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	taskIDStr := c.Param("id")
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Status": http.StatusBadRequest,
			"Message": "error while user id from context",
			"Error":   ""})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Status": http.StatusBadRequest,
			"Message": "error while converting user id to string",
			"Error":   ""})
		return
	}

	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  http.StatusBadRequest,
			"Message": "Invalid task ID",
			"Error":   err.Error(),
		})
		return
	}

	if err := h.service.DeleteTodoByIDSvc(ctx, taskID, userIDStr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Status":  http.StatusInternalServerError,
			"Message": "Error deleting todo",
			"Error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":  http.StatusOK,
		"Message": "Todo deleted successfully",
	})
}
