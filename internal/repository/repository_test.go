package repository

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"studentgit.kata.academy/gk/exchanger/internal/models"
	"testing"
)

func TestSaveCoinData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Ошибка при открытии GORM DB: %v", err)
	}

	repo := NewPsqlRepository(gormDB)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"coin_data\"").
		WithArgs("btcusdt", 50000.0, 49900.0, 111111).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	coinData := &models.CoinData{
		CoinsPair: "btcusdt",
		AskPrice:  50000.0,
		BidPrice:  49900.0,
		Time:      111111,
	}

	err = repo.SaveCoinData(coinData)

	assert.NoError(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestSaveCoinDataError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании mock DB: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Ошибка при открытии GORM DB: %v", err)
	}

	repo := NewPsqlRepository(gormDB)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"coin_data\"").
		WithArgs("btcusdt", 50000.0, 49900.0, 111111).
		WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	coinData := &models.CoinData{
		CoinsPair: "btcusdt",
		AskPrice:  50000.0,
		BidPrice:  49900.0,
		Time:      111111,
	}

	err = repo.SaveCoinData(coinData)

	assert.Error(t, err)
	assert.Equal(t, "error to insert coindata : some error", err.Error())
	assert.Nil(t, mock.ExpectationsWereMet())
}
