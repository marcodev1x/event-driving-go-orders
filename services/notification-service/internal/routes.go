package internal

import (
	"github.com/gin-gonic/gin"
)

type HttpMethod string

const (
	GET    HttpMethod = "GET"
	POST   HttpMethod = "POST"
	PUT    HttpMethod = "PUT"
	DELETE HttpMethod = "DELETE"
	PATCH  HttpMethod = "PATCH"
)

type RouteHandler struct {
	ToAuthenticated bool
	Path            string
	Handler         gin.HandlerFunc
	Method          HttpMethod
	Middlewares     []gin.HandlerFunc
}

func RouteDefiner(routes *[]RouteHandler, server *gin.Engine) {
	for _, route := range *routes {
		server.Group("/api", route.Middlewares...).Handle(string(route.Method), route.Path, route.Handler)
	}
}
