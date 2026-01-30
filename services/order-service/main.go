package main

import (
	"order-service/infra"
	"order-service/setup"
)

var (
	boot infra.Bootstrap
)

func main() {
	env := boot.LoadEnv()

	boot.SetupDatabase(env)

	server := boot.RunServer()

	boot.SetupRedis(env)

	setup.PrepareRoutes(server)

	server.Run(":8080")
}
