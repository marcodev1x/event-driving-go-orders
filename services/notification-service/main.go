package main

import (
	"notification-service/infra"
	"notification-service/setup"
)

var (
	boot infra.Bootstrap
)

func main() {
	env := boot.LoadEnv()

	// boot.SetupDatabase(env)

	server := boot.RunServer()

	boot.SetupRedis(env)

	setup.PrepareRoutes(server)

	server.Run(":8082")
}
