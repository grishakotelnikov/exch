package models

type CoinData struct {
	CoinsPair string
	AskPrice  float64
	BidPrice  float64
	Time      int64
}

type RawCoinData struct {
	Timestamp int `json:"timestamp"`
	Asks      []struct {
		Price string `json:"price"`
	} `json:"asks"`
	Bids []struct {
		Price string `json:"price"`
	} `json:"bids"`
}
