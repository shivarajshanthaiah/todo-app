package service

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/shivarajshanthaiah/todo-app/configs"
	"github.com/shivarajshanthaiah/todo-app/internal/jwt"
	"github.com/shivarajshanthaiah/todo-app/internal/models"
	"github.com/shivarajshanthaiah/todo-app/internal/repo/entity"
	repo "github.com/shivarajshanthaiah/todo-app/internal/repo/interfaces"
	service "github.com/shivarajshanthaiah/todo-app/internal/service/interfaces"
)

type UserService struct {
	repo repo.UserRepoInterface
	cnfg *configs.Config
}

func NewUserService(repo repo.UserRepoInterface, cnfg *configs.Config) service.UserServiceInterface {
	return &UserService{
		repo: repo,
		cnfg: cnfg,
	}
}

func (s *UserService) UserSignUpSvc(ctx context.Context, user *models.User) error {
	genetatedID, err := uuid.NewUUID()
	if err != nil {
		log.Println("Error while generating uuid for user", err)
		return err
	}

	hashedPassword, err := jwt.HashPassword(user.Password)
	if err != nil {
		log.Println("Error while hashing the password", err)
		return err
	}

	entityUser := &entity.User{
		ID:       genetatedID.String(),
		UserName: user.Username,
		Email:    user.Email,
		Password: hashedPassword,
	}
	err = s.repo.CreateUser(ctx, entityUser)
	if err != nil {
		log.Println("Error creating user in repo:", err)
		return err
	}
	return nil
}

func (s *UserService) UserLoginSvc(ctx context.Context, login *models.Login) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, login.Email)
	if err != nil {
		log.Println("Error while fetching the user from DB", err)
		return "", err
	}

	log.Printf("User fetched from DB: ID=%s, Email=%s, HashedPassword=%s", user.ID, user.Email, user.Password)

	match := jwt.CheckPassword(login.Password, user.Password)
	log.Printf("Password match result: %v (plain=%s, hashed=%s)", match, login.Password, user.Password)

	if !match {
		log.Printf("Incorrect password for user %s", user.Email)
		return "", errors.New("password incorrect")
	}

	token, err := jwt.GenerateToken(s.cnfg.SECRETKEY, user.Email, user.ID)
	if err != nil {
		log.Printf("Error generating token for user %s: %v", user.Email, err)
		return "", err
	}

	log.Printf("Login successful for user: %s", user.Email)
	return token, nil
}
