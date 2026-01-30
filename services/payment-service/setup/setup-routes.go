package setup

import (
	"payment-service/infra/config"
	"payment-service/internal"
	"payment-service/internal/rest"

	"github.com/gin-gonic/gin"
)

func PrepareRoutes(server *gin.Engine) {
	envs := config.LoadEnv()

	internal.RouteDefiner(rest.CheckoutRoutes(envs), server)
}
