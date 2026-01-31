package middlewares

import (
	"notification-service/infra/config"
	"notification-service/internal"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Middlewares struct{}

func (m *Middlewares) ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		lastErr := c.Errors.Last()

		if lastErr == nil {
			return
		}

		if structured, ok := lastErr.Err.(*internal.APIError); ok {
			config.Logger().Error("Erro na requisição.", lastErr.Err, structured.Code)

			internal.SendError(c, structured.Status, internal.HttpDTOStructured{
				Message: lastErr.Err.Error(),
				Code:    structured.Code,
			})

			return
		}

		config.Logger().Error("Erro na requisição.", lastErr.Err)

		internal.SendError(c, 500, internal.HttpDTOStructured{
			Message: "Erro interno do servidor.",
			Code:    000,
		})
	}
}

func (m *Middlewares) RateLimiter(seconds int, maxRequests int) gin.HandlerFunc {
	limits := rate.NewLimiter(rate.Limit(seconds), maxRequests)

	return func(c *gin.Context) {
		if limits.Allow() {
			c.Next()
		} else {
			internal.SendError(c, 429, internal.HttpDTOStructured{
				Message: "Muitas requisições no tempo limite.",
				Code:    429,
			})
		}

	}
}

var Interceptors = Middlewares{}
