package models

import "github.com/shopspring/decimal"

type StockQuote struct {
	Symbol           string          `json:"symbol"`
	Name             string          `json:"name"`
	Price            decimal.Decimal `json:"price"`
	Change           decimal.Decimal `json:"change"`
	ChangePercentage decimal.Decimal `json:"change_percentage"`
	Volume           int64           `json:"volume"`
	Open             decimal.Decimal `json:"open"`
	High             decimal.Decimal `json:"high"`
	Low              decimal.Decimal `json:"low"`
	LastUpdated      int64           `json:"last_updated"`
}

type MarketStatus struct {
	IsOpen bool   `json:"is_open"`
	Status string `json:"status"` // "OPEN", "CLOSED"
}

type MarketIndex struct {
	IndexName  string          `json:"index_name"`
	Current    decimal.Decimal `json:"current"`
	Change     decimal.Decimal `json:"change"`
	Percentage decimal.Decimal `json:"percentage"`
	High       decimal.Decimal `json:"high"`
	Low        decimal.Decimal `json:"low"`
}

type MarketSummary struct {
	Status   MarketStatus  `json:"status"`
	KSE100   MarketIndex   `json:"kse100"`
	Indices  []MarketIndex `json:"indices"`
	Equities []StockQuote  `json:"equities"`
}
