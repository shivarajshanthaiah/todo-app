package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shivarajshanthaiah/todo-app/configs"
	"github.com/shivarajshanthaiah/todo-app/internal/handler"
	"github.com/shivarajshanthaiah/todo-app/internal/middleware"
)

func RegisterRoutes(router *gin.Engine, todoHndlr *handler.TaskHandler, userHndlr *handler.UserHandler, cnfg *configs.Config) {

	v1 := router.Group("/api/v1")
	{
		v1.POST("/signup", userHndlr.UserSignUpHandler)
		v1.POST("/login", userHndlr.UserLoginHandler)
	}

	user := v1.Group("user")
	user.Use(middleware.Authorization(cnfg.SECRETKEY))
	{
		user.POST("/todos", todoHndlr.CreateTodoHandler)
		user.POST("/todos/list", todoHndlr.GetTodosHandler)
		user.PATCH("/todos/:id", todoHndlr.UpdateTodoHandler)
		// user.PUT("/todos", todoHndlr.UpdateTodoHandler)
		user.DELETE("/todos/:id", todoHndlr.DeleteTodoHandler)
		user.GET("/get/profile", userHndlr.GetUserProfileHandler)
	}
}
