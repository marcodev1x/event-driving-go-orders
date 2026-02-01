package tests

import (
	"context"
	"errors"
	"payment-service/internal"
	"payment-service/internal/domain"
	"payment-service/internal/usecases"
	mockusecases "payment-service/internal/usecases/mocks"
	events "payment-service/kafka/events/domain"
	mockkafka "payment-service/kafka/mocks"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var timeTesting = time.Now()

func TestPaymentUseCases_ValidatePayment(t *testing.T) {
	// O microservice não valida realmente se pagamentos foram válidos, já que é um projeto somente para desenho de arquitetura, k8s e testing.
	// Além de demonstrar implementação de Go, Gin, GORM, entre outros. Então os testes vão validar conteúdo mocked mesmo (preço 100 = valid, preço 50 = failed...).

	mockedEvent := events.OrderCreated{
		BaseEvent: events.BaseEvent{
			EventID:   uuid.New().String(),
			EventType: "order.created",
			Timestamp: timeTesting,
			ContentID: 1,
		},
		Checkout: domain.Checkout{
			ID:    1,
			Price: 100,
		},
	}
	orderId := 1

	testsNoError := []struct {
		testcase string
		price    float64
		key      string
	}{
		{"should return payment.confirmed if price is 100", 100, "payment.confirmed"},
		{"should return payment.failed if price is 50", 50, "payment.failed"},
	}

	for _, test := range testsNoError {
		t.Run(test.testcase, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			cacheMocked := mockusecases.NewMockCache(controller)
			producerMocked := mockkafka.NewMockEventProducer(controller)

			uses := usecases.NewPaymentUseCase(cacheMocked, producerMocked)

			producerMocked.
				EXPECT().
				PublishEvent(context.Background(), test.key, gomock.AssignableToTypeOf(events.PaymentInvoice{})).
				Return(nil).
				Times(1)

			// muda o price sem mudar a struct base
			event := mockedEvent
			event.Checkout.Price = test.price

			err := uses.ValidatePayment(event, orderId)

			assert.NoError(t, err)
		})
	}

	testsError := []struct {
		testcase string
		price    float64
	}{
		{"should throws internal API Error if backoff fails with price 100", 100},
		{"should throws internal API Error if backoff fails with price 50", 50},
	}

	for _, test := range testsError {
		t.Run(test.testcase, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			cacheMocked := mockusecases.NewMockCache(controller)
			producerMocked := mockkafka.NewMockEventProducer(controller)

			uses := usecases.NewPaymentUseCase(cacheMocked, producerMocked)

			producerMocked.
				EXPECT().
				PublishEvent(context.Background(), gomock.Any(), gomock.AssignableToTypeOf(events.PaymentInvoice{})).
				Return(errors.New("error")).
				AnyTimes() // Valida tentativas de backoff. Demora mais por conta de retry

			// muda o price sem mudar a struct base
			event := mockedEvent
			event.Checkout.Price = test.price
			expectedError := internal.NewAPIError("Erro ao enviar evento mesmo após diversas tentativas.", 500, 200)

			err := uses.ValidatePayment(event, orderId)

			assert.Equal(t, expectedError, err)
		})
	}
}
