package boot

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shivarajshanthaiah/todo-app/configs"
	"github.com/shivarajshanthaiah/todo-app/internal/clients/psql"
	"github.com/shivarajshanthaiah/todo-app/internal/clients/redis"
	"github.com/shivarajshanthaiah/todo-app/internal/handler"
	"github.com/shivarajshanthaiah/todo-app/internal/repo"
	"github.com/shivarajshanthaiah/todo-app/internal/routes"
	"github.com/shivarajshanthaiah/todo-app/internal/service"
	"go.uber.org/zap"
)

// Server represents the model of the server with a Gin engine.
type Server struct {
	R      *gin.Engine
	DB     *pgxpool.Pool
	Cnfg   *configs.Config
	Redis  *redis.RedisService
	Logger *zap.Logger
}

// StartServer method starts the server on the specified port.
func (s *Server) StartServer(port string) error {
	if s.DB != nil {
		fmt.Println("not nil")
	}
	taskRepo := repo.NewTaskRepository(s.DB)
	taskSvc := service.NewTaskService(taskRepo, s.Logger)
	taskHandler := handler.NewTaskHandler(taskSvc)

	userRepo := repo.NewUserRepository(s.DB)
	userSvc := service.NewUserService(userRepo, s.Cnfg, s.Redis, s.Logger)
	userHandler := handler.NewUserHandler(userSvc)

	routes.RegisterRoutes(s.R, taskHandler, userHandler, s.Cnfg)
	return s.R.Run(":" + port)
}

// NewServer returns a new Server instance with dependencies injected.
func NewServer(db *pgxpool.Pool, redisClient *redis.RedisService, logger *zap.Logger, cnfg *configs.Config) *Server {
	engine := gin.Default()
	return &Server{
		R:      engine,
		DB:     db,
		Cnfg:   cnfg,
		Redis:  redisClient,
		Logger: logger,
	}
}

// Intitialise external conncetions
func Setup() (*configs.Config, *pgxpool.Pool, *redis.RedisService, *zap.Logger, error) {
	// Load configuration
	cnfg := configs.LoadConfig()

	// Connect to Postgres
	db, err := psql.NewPsql(cnfg)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to connect to postgres: %v", err)
	}
	log.Println(" Successfully connected to PostgreSQL")

	// Run migrations (if you have one)
	if err := runMigrations(db); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	// Connect to Redis
	redisClient, err := redis.SetupRedis(cnfg)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to connect to redis: %v", err)
	}
	log.Println("Successfully connected to Redis")

	// Initialize Zap logger
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to initialize zap logger: %v", err)
	}

	return cnfg, db, redisClient, logger, nil
}

// This should be read from .sql file and run. Im showing demo migration run with this func
func runMigrations(db *pgxpool.Pool) error {
	schema := `
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(63) NOT NULL PRIMARY KEY,
    username VARCHAR(63) NOT NULL,                     
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    email VARCHAR(63) NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS emailusername ON users (lower(email));

CREATE TABLE IF NOT EXISTS tasks (
  id SERIAL PRIMARY KEY,
  user_id VARCHAR(63) NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now(),
  title VARCHAR(119) NOT NULL,
  description TEXT,
  priority INT NOT NULL DEFAULT 1,
  status INT NOT NULL DEFAULT 1,
  due_at TIMESTAMP
);
`
	_, err := db.Exec(context.Background(), schema)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	log.Println("Database migrations applied successfully.")
	return nil
}
