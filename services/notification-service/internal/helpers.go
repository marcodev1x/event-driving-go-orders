package internal

import (
	"notification-service/internal/domain"

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

func HandlePaymentMethodToClient(method domain.PaymentMethod) string {
	switch method {
	case domain.Pix:
		return "Pix"
	case domain.Boleto:
		return "Boleto"
	case domain.CreditCard:
		return "Cartão de Crédito"
	case domain.DebitCard:
		return "Cartão de Débito"
	case domain.CommonMoney:
		return "Dinheiro"
	}

	return "Método de pagamento não informado"
}
