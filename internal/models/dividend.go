package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TaxStatus string

const (
	TaxFiler    TaxStatus = "FILER"     // 15% Withholding Tax
	TaxNonFiler TaxStatus = "NON_FILER" // 30% Withholding Tax
)

type DividendAnnouncement struct {
	ID              uuid.UUID       `json:"id"`
	Symbol          string          `json:"symbol"`
	PayoutPerShare  decimal.Decimal `json:"payout_per_share"`
	BookClosureDate time.Time       `json:"book_closure_date"`
	PaymentDate     *time.Time      `json:"payment_date,omitempty"`
	FinancialPeriod string          `json:"financial_period"`
}

type UserDividendProjection struct {
	Symbol           string          `json:"symbol"`
	SharesOwned      decimal.Decimal `json:"shares_owned"`
	PayoutPerShare   decimal.Decimal `json:"payout_per_share"`
	GrossPayout      decimal.Decimal `json:"gross_payout"`
	TaxRatePct       decimal.Decimal `json:"tax_rate_pct"` // 15 or 30
	TaxDeduction     decimal.Decimal `json:"tax_deduction"`
	NetPayout        decimal.Decimal `json:"net_payout"`
	BookClosureDate  time.Time       `json:"book_closure_date"`
	PaymentDate      *time.Time      `json:"payment_date,omitempty"`
}

type PortfolioDividendForecast struct {
	PortfolioID   uuid.UUID                `json:"portfolio_id"`
	TaxStatus     TaxStatus                `json:"tax_status"`
	TotalGross    decimal.Decimal          `json:"total_gross"`
	TotalTax      decimal.Decimal          `json:"total_tax"`
	TotalNet      decimal.Decimal          `json:"total_net"`
	Projections   []UserDividendProjection `json:"projections"`
}
