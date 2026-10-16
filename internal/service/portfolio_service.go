package service

import (
	"psx-engine-go/internal/models"
	"psx-engine-go/internal/psxscraper"
	"github.com/shopspring/decimal"
)

type PortfolioService struct {
	ingestor *psxscraper.PSXIngestor
}

func NewPortfolioService(ingestor *psxscraper.PSXIngestor) *PortfolioService {
	return &PortfolioService{
		ingestor: ingestor,
	}
}

func (ps *PortfolioService) CalculatePortfolioSummary(p models.Portfolio, txs []models.Transaction) models.PortfolioSummary {
	holdingsMap := make(map[string]*models.HoldingSummary)

	for _, tx := range txs {
		h, exists := holdingsMap[tx.Symbol]
		if !exists {
			h = &models.HoldingSummary{
				Symbol:        tx.Symbol,
				TotalQuantity: decimal.Zero,
				AvgCostBasis:  decimal.Zero,
				TotalCost:     decimal.Zero,
			}
			holdingsMap[tx.Symbol] = h
		}

		if tx.TradeType == models.TradeBuy {
			cost := tx.Quantity.Mul(tx.PricePerShare).Add(tx.Commission)
			h.TotalCost = h.TotalCost.Add(cost)
			h.TotalQuantity = h.TotalQuantity.Add(tx.Quantity)
			if !h.TotalQuantity.IsZero() {
				h.AvgCostBasis = h.TotalCost.Div(h.TotalQuantity)
			}
		} else if tx.TradeType == models.TradeSell {
			h.TotalQuantity = h.TotalQuantity.Sub(tx.Quantity)
			if h.TotalQuantity.LessThanOrEqual(decimal.Zero) {
				h.TotalQuantity = decimal.Zero
				h.TotalCost = decimal.Zero
				h.AvgCostBasis = decimal.Zero
			} else {
				h.TotalCost = h.TotalQuantity.Mul(h.AvgCostBasis)
			}
		}
	}

	var holdingsList []models.HoldingSummary
	totalInvested := decimal.Zero
	totalCurrentValue := decimal.Zero

	for symbol, h := range holdingsMap {
		if h.TotalQuantity.GreaterThan(decimal.Zero) {
			quote, found := ps.ingestor.GetCachedQuote(symbol)
			if !found {
				quote, _ = ps.ingestor.FetchQuote(symbol)
			}

			h.CurrentPrice = quote.Price
			h.CurrentValue = h.TotalQuantity.Mul(quote.Price)
			h.UnrealizedGain = h.CurrentValue.Sub(h.TotalCost)
			if !h.TotalCost.IsZero() {
				h.UnrealizedGainPct = h.UnrealizedGain.Div(h.TotalCost).Mul(decimal.NewFromInt(100))
			}

			totalInvested = totalInvested.Add(h.TotalCost)
			totalCurrentValue = totalCurrentValue.Add(h.CurrentValue)
			holdingsList = append(holdingsList, *h)
		}
	}

	totalUnrealized := totalCurrentValue.Sub(totalInvested)
	var totalUnrealizedPct decimal.Decimal
	if !totalInvested.IsZero() {
		totalUnrealizedPct = totalUnrealized.Div(totalInvested).Mul(decimal.NewFromInt(100))
	}

	return models.PortfolioSummary{
		PortfolioID:        p.ID,
		Name:               p.Name,
		TotalInvested:      totalInvested,
		CurrentValue:       totalCurrentValue,
		TotalUnrealized:    totalUnrealized,
		TotalUnrealizedPct: totalUnrealizedPct,
		Holdings:           holdingsList,
	}
}
