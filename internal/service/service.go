package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"studentgit.kata.academy/gk/exchanger/internal/models"
	"studentgit.kata.academy/gk/exchanger/internal/repository"
)

type Exchangerer interface {
	GetRates(ctx context.Context, coinPair string) (*models.CoinData, error)
}

type GarantexExchange struct {
	rep repository.Repository
}

func NewGarantexExchange(rep repository.Repository) *GarantexExchange {
	return &GarantexExchange{
		rep: rep,
	}
}

func (g *GarantexExchange) GetRates(ctx context.Context, coinPair string) (*models.CoinData, error) {
	if len(coinPair) == 0 || len(coinPair) > 8 {
		return nil, fmt.Errorf("incorrect pair - %s", coinPair)
	}

	url := fmt.Sprintf("http://garantex.org/api/v2/depth?market=%s", coinPair)
	responce, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error to get coin data %s, from url - %s : %v", coinPair, url, err)
	}
	defer responce.Body.Close()

	if responce.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error to get requsest from url - %s : unaccepted code - %d", url, responce.StatusCode)
	}
	var rawData models.RawCoinData
	err = json.NewDecoder(responce.Body).Decode(&rawData)
	if err != nil {
		return nil, fmt.Errorf("error to decode coin data  : %v", err)
	}

	resData, err := g.rawToRes(rawData, coinPair)
	if err != nil {
		return nil, fmt.Errorf("error to parse raw  data  : %v", err)
	}

	err = g.rep.SaveCoinData(resData)
	if err != nil {
		return nil, fmt.Errorf("error to save   data into db : %v", err)
	}

	return resData, nil
}

func (g *GarantexExchange) rawToRes(data models.RawCoinData, pair string) (*models.CoinData, error) {

	var resData models.CoinData

	resData.CoinsPair = pair
	askFloat64, err := strconv.ParseFloat(data.Asks[0].Price, 64)
	if err != nil {
		return nil, fmt.Errorf("error to parse ask price  : %v", err)
	}

	bidFloat64, err := strconv.ParseFloat(data.Bids[0].Price, 64)
	if err != nil {
		return nil, fmt.Errorf("error to parse bid price  : %v", err)
	}

	resData.AskPrice = askFloat64
	resData.BidPrice = bidFloat64
	resData.Time = int64(data.Timestamp)

	return &resData, nil
}
