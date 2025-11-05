package main

import (
	"log"

	"github.com/shivarajshanthaiah/todo-app/internal/boot"
)

func main() {

	cnfg, db, redis, logger, err := boot.Setup()
	if err != nil {
		log.Fatalf("failed to initialize dependencies: %v", err)
	}
	server := boot.NewServer(db, redis, logger, cnfg)
	if err := server.StartServer(cnfg.SERVERPORT); err != nil {
		log.Fatalf("server exited with error: %v", err)
	}
}
