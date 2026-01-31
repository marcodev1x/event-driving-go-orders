package setup

import (
	"notification-service/infra/config"
	"notification-service/internal"
	"notification-service/internal/rest"

	"github.com/gin-gonic/gin"
)

func PrepareRoutes(server *gin.Engine) {
	envs := config.LoadEnv()

	internal.RouteDefiner(rest.CheckoutRoutes(envs), server)
}
