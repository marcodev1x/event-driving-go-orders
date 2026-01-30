package rest

import (
	"order-service/internal"
	"order-service/internal/structs"
	"order-service/internal/usecases"

	"github.com/gin-gonic/gin"
)

type CheckoutRest struct {
	usecase *usecases.CheckoutUsecase
}

func NewCheckoutRest(usecases *usecases.CheckoutUsecase) *CheckoutRest {
	return &CheckoutRest{
		usecases,
	}
}

func (r *CheckoutRest) CreateCheckout(c *gin.Context) {
	var request structs.CreateCheckout

	if err := internal.BindJSON(c, &request); err != nil {
		c.Error(err)
		return
	}

	if _, err := r.usecase.CreateCheckout(request); err != nil {
		c.Error(err)
		return
	}

	internal.SendResponse(c, 201, "Pedido criado com sucesso.")
}
