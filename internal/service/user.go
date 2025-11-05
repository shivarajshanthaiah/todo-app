package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/shivarajshanthaiah/todo-app/configs"
	redisCl "github.com/shivarajshanthaiah/todo-app/internal/clients/redis"
	"github.com/shivarajshanthaiah/todo-app/internal/jwt"
	"github.com/shivarajshanthaiah/todo-app/internal/models"
	"github.com/shivarajshanthaiah/todo-app/internal/repo/entity"
	repo "github.com/shivarajshanthaiah/todo-app/internal/repo/interfaces"
	service "github.com/shivarajshanthaiah/todo-app/internal/service/interfaces"
	"go.uber.org/zap"
)

type UserService struct {
	repo   repo.UserRepoInterface
	cnfg   *configs.Config
	redis  *redisCl.RedisService
	logger *zap.Logger
}

func NewUserService(repo repo.UserRepoInterface, cnfg *configs.Config, redis *redisCl.RedisService, logger *zap.Logger) service.UserServiceInterface {
	return &UserService{
		repo:   repo,
		cnfg:   cnfg,
		redis:  redis,
		logger: logger,
	}
}

func (s *UserService) UserSignUpSvc(ctx context.Context, user *models.User) error {
	genetatedID, err := uuid.NewUUID()
	if err != nil {
		s.logger.Error("Error while generating uuid for user", zap.Error(err)) // imitate this for all the implementation of log
		return err
	}

	hashedPassword, err := jwt.HashPassword(user.Password)
	if err != nil {
		log.Println("Error while hashing the password", err)
		return err
	}

	entityUser := &entity.User{
		ID:       genetatedID.String(),
		UserName: user.UserName,
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

func (s *UserService) GetUserByIDSvc(ctx context.Context, userID string) (*models.User, string, error) {

	cacheKey := "user_" + userID
	cachedUser, err := s.redis.GetFromRedis(cacheKey)

	if err == nil {
		// Unmarshal cached data into a User struct if cache hit
		var user models.User
		if err := json.Unmarshal([]byte(cachedUser), &user); err != nil {
			log.Printf("Error unmarshalling cached user data: %v", err)
		} else {
			return &models.User{
				ID:       user.ID,
				UserName: user.UserName,
				Email:    user.Email,
			}, "fethed from cache", nil
		}
	} else if err != redis.Nil {
		// Handle Redis error (other than missing key)
		log.Printf("Error accessing Redis: %v", err)
		return nil, "", err
	}

	// Cache miss: fetch user from the database
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		log.Println("Error while fetching the user from DB", err)
		return nil, "", err
	}

	userModel := &models.User{
		ID:       user.ID,
		UserName: user.UserName,
		Email:    user.Email,
		Password: "", // don’t expose password
	}

	// Cache the retrieved user data for future requests
	userData, err := json.Marshal(user)
	if err == nil {
		// seting 2 mins cachefor testing
		_ = s.redis.SetDataInRedis(cacheKey, userData, time.Minute*2)
	} else {
		log.Printf("Error marshalling user data for caching: %v", err)
	}

	return userModel, "fetched from DB", nil
}
