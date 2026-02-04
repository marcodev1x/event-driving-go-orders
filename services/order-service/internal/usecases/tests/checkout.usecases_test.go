package tests

import (
	"errors"
	"order-service/internal"
	"order-service/internal/domain"
	repomocks "order-service/internal/repository/mocks"
	"order-service/internal/structs"
	"order-service/internal/usecases"
	"order-service/internal/usecases/mocks"
	mockkafka "order-service/kafka/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCheckoutUseCases_CreateCheckout(t *testing.T) {

	mockedPadronized := &structs.CreateCheckout{
		Price:         100,
		PaymentMethod: domain.Pix,
	}

	structErrCases := []struct {
		testcase    string
		validateMsg string
		structure   structs.CreateCheckout
	}{
		{
			"should throw validation error with price is not set in request body",
			"'Price' failed on the 'required' tag",
			structs.CreateCheckout{
				PaymentMethod: domain.Pix,
			},
		},
		{
			"should throw validation error with price is lower or equal 0 in request body",
			"'Price' failed on the 'min' tag",
			structs.CreateCheckout{
				PaymentMethod: domain.Pix,
				Price:         0.9,
			},
		},
		{
			"should throw validation error if payment method is nil",
			"'PaymentMethod' failed on the 'required' tag",
			structs.CreateCheckout{
				Price: 100,
			},
		},
		{
			"should throw validation error if payment method is not type of the expected.",
			"'PaymentMethod' failed on the 'oneof' tag",
			structs.CreateCheckout{
				PaymentMethod: "invalid",
				Price:         100,
			},
		},
	}

	for _, test := range structErrCases {
		t.Run(test.testcase, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			mockRepo := repomocks.NewMockCheckoutImplementation(controller)
			mockRedis := mocks.NewMockRedisImplementation(controller)
			mockProducer := mockkafka.NewMockEventProducer(controller)
			svc := usecases.NewCheckoutUseCase(mockRepo, mockRedis, mockProducer)

			_, err := svc.CreateCheckout(test.structure)

			expectedError := internal.NewAPIError(err.Error(), 400, 101)

			assert.Equal(t, expectedError, err)
			assert.Contains(t, err.Error(), test.validateMsg)
		})
	}

	t.Run("should throw ApiError repo registry fails", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()

		mockRepo := repomocks.NewMockCheckoutImplementation(controller)
		mockRedis := mocks.NewMockRedisImplementation(controller)
		mockProducer := mockkafka.NewMockEventProducer(controller)
		svc := usecases.NewCheckoutUseCase(mockRepo, mockRedis, mockProducer)

		mockRepo.
			EXPECT().
			CreateCheckout(gomock.Any()).
			Return(nil, errors.New("error")).
			Times(1)

		_, err := svc.CreateCheckout(*mockedPadronized)

		expectedError := internal.NewAPIError("Erro ao criar checkout.", 500, 102)

		assert.Equal(t, expectedError, err)
	})
}
