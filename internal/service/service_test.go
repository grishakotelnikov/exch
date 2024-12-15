package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"studentgit.kata.academy/gk/exchanger/internal/models"
	"testing"
)

type MockRepository struct {
}

func (m *MockRepository) SaveCoinData(data *models.CoinData) error {
	return nil
}

func TestGarantexExchange_GetRates(t *testing.T) {
	mockRepo := &MockRepository{}
	service := NewGarantexExchange(mockRepo)

	t.Run("success case", func(t *testing.T) {
		coinPair := "btcusdt"

		coinData, err := service.GetRates(context.Background(), coinPair)

		assert.NoError(t, err)
		assert.NotNil(t, coinData)
		assert.Equal(t, "btcusdt", coinData.CoinsPair)
		assert.Greater(t, coinData.AskPrice, 0.0)
		assert.Greater(t, coinData.BidPrice, 0.0)
		assert.Greater(t, coinData.Time, int64(0))
	})

	t.Run("error case - incorrect pair", func(t *testing.T) {
		coinPair := "btc"

		coinData, err := service.GetRates(context.Background(), coinPair)

		assert.Error(t, err)
		assert.Nil(t, coinData)
		assert.Contains(t, err.Error(), "error to get requsest")
	})

	t.Run("error case - nil req", func(t *testing.T) {
		coinPair := ""

		// Вызываем метод
		coinData, err := service.GetRates(context.Background(), coinPair)

		assert.Error(t, err)
		assert.Nil(t, coinData)
		assert.Contains(t, err.Error(), "incorrect pair")
	})
}
