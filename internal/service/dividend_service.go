package service

import (
	"psx-engine-go/internal/models"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type DividendService struct{}

func NewDividendService() *DividendService {
	return &DividendService{}
}

// CalculateDividendForecast applies FBR tax logic (15% Filer / 30% Non-Filer)
func (ds *DividendService) CalculateDividendForecast(
	portfolioID uuid.UUID,
	taxStatus models.TaxStatus,
	holdings []models.HoldingSummary,
	announcements []models.DividendAnnouncement,
) models.PortfolioDividendForecast {

	taxRatePct := decimal.NewFromInt(15)
	if taxStatus == models.TaxNonFiler {
		taxRatePct = decimal.NewFromInt(30)
	}

	var projections []models.UserDividendProjection
	totalGross := decimal.Zero
	totalTax := decimal.Zero
	totalNet := decimal.Zero

	annMap := make(map[string]models.DividendAnnouncement)
	for _, a := range announcements {
		annMap[a.Symbol] = a
	}

	for _, h := range holdings {
		if ann, exists := annMap[h.Symbol]; exists && h.TotalQuantity.GreaterThan(decimal.Zero) {
			gross := h.TotalQuantity.Mul(ann.PayoutPerShare)
			tax := gross.Mul(taxRatePct).Div(decimal.NewFromInt(100))
			net := gross.Sub(tax)

			proj := models.UserDividendProjection{
				Symbol:          h.Symbol,
				SharesOwned:     h.TotalQuantity,
				PayoutPerShare:  ann.PayoutPerShare,
				GrossPayout:     gross,
				TaxRatePct:      taxRatePct,
				TaxDeduction:    tax,
				NetPayout:       net,
				BookClosureDate: ann.BookClosureDate,
				PaymentDate:     ann.PaymentDate,
			}

			totalGross = totalGross.Add(gross)
			totalTax = totalTax.Add(tax)
			totalNet = totalNet.Add(net)
			projections = append(projections, proj)
		}
	}

	return models.PortfolioDividendForecast{
		PortfolioID: portfolioID,
		TaxStatus:   taxStatus,
		TotalGross:  totalGross,
		TotalTax:    totalTax,
		TotalNet:    totalNet,
		Projections: projections,
	}
}
