package main

import "github.com/shivarajshanthaiah/todo-app/internal/boot"

func main() {

	boot.Setup()
	server := boot.NewServer()
	server.StartServer(boot.Cnfg.SERVERPORT)
}
