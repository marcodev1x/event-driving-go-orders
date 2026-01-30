package setup

import (
	"order-service/infra/config"
	"order-service/internal"
	"order-service/internal/rest"

	"github.com/gin-gonic/gin"
)

func PrepareRoutes(server *gin.Engine) {
	envs := config.LoadEnv()

	internal.RouteDefiner(rest.CheckoutRoutes(envs), server)
}
