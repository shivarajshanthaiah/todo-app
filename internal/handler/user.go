package handler

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shivarajshanthaiah/todo-app/internal/models"
	"github.com/shivarajshanthaiah/todo-app/internal/service/interfaces"
)

type UserHandler struct {
	service interfaces.UserServiceInterface
}

func NewUserHandler(service interfaces.UserServiceInterface) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) UserSignUpHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, time.Second*100)
	defer cancel()

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Printf("UserSignUpHandler Error in binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Status": http.StatusBadRequest,
			"Message": "error in binding data",
			"Error":   err.Error()})
		return
	}

	log.Printf("UserSignUpHandler JSON binding successful: %+v", user)

	if err := h.service.UserSignUpSvc(ctx, &user); err != nil {
		log.Printf("UserSignUpHandler Error in signup service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"Status": http.StatusInternalServerError,
			"Message": "error in signup service",
			"Error":   err.Error()})
		return
	}
	log.Printf("UserSignUpHandler Signup service completed successfully for user: %s", user.Email)

	c.JSON(http.StatusCreated, gin.H{
		"Status":  http.StatusCreated,
		"Message": "user signed up successfully",
		"Data":    "",
	})
	log.Print("UserSignUpHandler Handler completed successfully")
}

func (h *UserHandler) UserLoginHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, time.Second*100)
	defer cancel()

	var user models.Login
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Status": http.StatusBadRequest,
			"Message": "error in binding data",
			"Error":   err.Error()})
		return
	}

	token, err := h.service.UserLoginSvc(ctx, &user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Status": http.StatusBadRequest,
			"Message": "error in login service",
			"Error":   err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"Status":  http.StatusAccepted,
		"Message": "user logged in successfully",
		"Data":    token,
	})
}

func (h *UserHandler) GetUserProfileHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, time.Second*100)
	defer cancel()

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

	user, messge, err := h.service.GetUserByIDSvc(ctx, userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Status":  http.StatusInternalServerError,
			"Message": "Error fetching user profile",
			"Error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":  http.StatusOK,
		"Message": messge,
		"Data":    user,
	})
}