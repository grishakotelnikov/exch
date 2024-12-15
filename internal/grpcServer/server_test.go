package grpcserver

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"studentgit.kata.academy/gk/exchanger/internal/models"
	exchanger "studentgit.kata.academy/gk/exchanger/internal/proto"
)

type MockExchangerService struct {
	mock.Mock
}

func (m *MockExchangerService) GetRates(ctx context.Context, request string) (*models.CoinData, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CoinData), args.Error(1)
}

func TestGetRates(t *testing.T) {
	tests := []struct {
		name          string
		req           *exchanger.CryptoRequest
		mockSetup     func(m *MockExchangerService)
		expectedRes   *exchanger.ValueResponce
		expectedError error
	}{
		{
			name: "success case",
			req: &exchanger.CryptoRequest{
				Request: "btcusdt",
			},
			mockSetup: func(m *MockExchangerService) {
				m.On("GetRates", mock.Anything, "btcusdt").Return(&models.CoinData{
					AskPrice: 50000.0,
					BidPrice: 49900.0,
					Time:     111111,
				}, nil)
			},
			expectedRes: &exchanger.ValueResponce{
				Asks:        50000.0,
				Bids:        49900.0,
				CurrentTime: 111111,
			},
			expectedError: nil,
		},
		{
			name: "error case",
			req: &exchanger.CryptoRequest{
				Request: "btcusdt",
			},
			mockSetup: func(m *MockExchangerService) {
				m.On("GetRates", mock.Anything, "btcusdt").Return(nil, errors.New("some error"))
			},
			expectedRes:   nil,
			expectedError: errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockExchangerService)
			tt.mockSetup(mockService)

			server := NewExchangerGrpcService(mockService, nil)

			res, err := server.GetRates(context.Background(), tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRes, res)
			}

			mockService.AssertExpectations(t)
		})
	}
}
