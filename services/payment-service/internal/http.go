package internal

import "github.com/gin-gonic/gin"

type HttpDTOStructured struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func SendResponse(c *gin.Context, status int, content any) {
	c.JSON(status, gin.H{"data": content})
}

func SendError(c *gin.Context, status int, content HttpDTOStructured) {
	c.AbortWithStatusJSON(status, gin.H{"error": content})
}
