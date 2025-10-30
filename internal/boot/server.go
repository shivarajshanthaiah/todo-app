package boot

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shivarajshanthaiah/todo-app/configs"
	"github.com/shivarajshanthaiah/todo-app/internal/db"
	"github.com/shivarajshanthaiah/todo-app/internal/handler"
	"github.com/shivarajshanthaiah/todo-app/internal/repo"
	"github.com/shivarajshanthaiah/todo-app/internal/routes"
	"github.com/shivarajshanthaiah/todo-app/internal/service"
)

var (
	DB   *pgxpool.Pool
	Cnfg = configs.LoadConfig()
)

// Server represents the model of the server with a Gin engine.
type Server struct {
	R *gin.Engine
}

// StartServer method starts the server on the specified port.
func (s *Server) StartServer(port string) {
	if DB != nil {
		fmt.Println("not nil")
	}
	taskRepo := repo.NewTaskRepository(DB)
	taskSvc := service.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskSvc)

	userRepo := repo.NewUserRepository(DB)
	userSvc := service.NewUserService(userRepo, Cnfg)
	userHandler := handler.NewUserHandler(userSvc)

	routes.RegisterRoutes(s.R, taskHandler, userHandler, Cnfg)
	s.R.Run(":" + port)
}

// NewServer returns a new Server instance with the default Gin engine attached.
func NewServer() *Server {
	engine := gin.Default()

	return &Server{
		R: engine,
	}
}

func Setup() {
	var err error
	Cnfg = configs.LoadConfig()
	DB, err = db.NewPsql(Cnfg)
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %v", err)
	}
}
