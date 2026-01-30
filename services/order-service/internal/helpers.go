package internal

import (
	"github.com/gin-gonic/gin"
)

type APIError struct {
	Message string `json:"message"`
	Status  int    `json:"-"`
	Code    int    `json:"code"`
}

func NewAPIError(message string, status, code int) error {
	return &APIError{
		Message: message,
		Status:  status,
		Code:    code,
	}
}

func (e *APIError) Error() string {
	return e.Message
}

func BindJSON(ctx *gin.Context, dest any) error {
	if err := ctx.ShouldBindJSON(dest); err != nil {
		return NewAPIError(
			"Erro com dados da requisição.",
			400,
			100,
		)
	}

	return nil
}
