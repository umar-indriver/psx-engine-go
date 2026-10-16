package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TradeType string

const (
	TradeBuy  TradeType = "BUY"
	TradeSell TradeType = "SELL"
)

type Portfolio struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

type Transaction struct {
	ID            uuid.UUID       `json:"id"`
	PortfolioID   uuid.UUID       `json:"portfolio_id"`
	Symbol        string          `json:"symbol"`
	TradeType     TradeType       `json:"trade_type"`
	Quantity      decimal.Decimal `json:"quantity"`
	PricePerShare decimal.Decimal `json:"price_per_share"`
	Commission    decimal.Decimal `json:"commission"`
	TradeDate     time.Time       `json:"trade_date"`
}

type HoldingSummary struct {
	Symbol           string          `json:"symbol"`
	TotalQuantity    decimal.Decimal `json:"total_quantity"`
	AvgCostBasis     decimal.Decimal `json:"avg_cost_basis"`
	TotalCost        decimal.Decimal `json:"total_cost"`
	CurrentPrice     decimal.Decimal `json:"current_price"`
	CurrentValue     decimal.Decimal `json:"current_value"`
	UnrealizedGain   decimal.Decimal `json:"unrealized_gain"`
	UnrealizedGainPct decimal.Decimal `json:"unrealized_gain_pct"`
	AnnualDividend   decimal.Decimal `json:"annual_dividend"`
	YieldOnCostPct   decimal.Decimal `json:"yield_on_cost_pct"`
}

type PortfolioSummary struct {
	PortfolioID      uuid.UUID        `json:"portfolio_id"`
	Name             string           `json:"name"`
	TotalInvested    decimal.Decimal  `json:"total_invested"`
	CurrentValue     decimal.Decimal  `json:"current_value"`
	TotalUnrealized  decimal.Decimal  `json:"total_unrealized"`
	TotalUnrealizedPct decimal.Decimal `json:"total_unrealized_pct"`
	AnnualDividend   decimal.Decimal  `json:"annual_dividend"`
	Holdings         []HoldingSummary `json:"holdings"`
}
