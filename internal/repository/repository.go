package repository

import (
	"fmt"
	"gorm.io/gorm"
	"studentgit.kata.academy/gk/exchanger/internal/models"
)

type Repository interface {
	SaveCoinData(data *models.CoinData) error
}

type PsqlRepository struct {
	db *gorm.DB
}

func NewPsqlRepository(db *gorm.DB) *PsqlRepository {
	return &PsqlRepository{
		db: db,
	}
}

func (psql *PsqlRepository) SaveCoinData(data *models.CoinData) error {
	res := psql.db.Create(data)
	if res.Error != nil {
		return fmt.Errorf("error to insert coindata : %v", res.Error)
	}
	return nil
}
