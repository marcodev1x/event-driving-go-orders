package main

import (
	"payment-service/infra"
	"payment-service/setup"
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
