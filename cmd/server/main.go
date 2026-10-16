package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"psx-engine-go/internal/config"
	"psx-engine-go/internal/models"
	"psx-engine-go/internal/psxscraper"
	"psx-engine-go/internal/service"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func main() {
	cfg := config.LoadConfig()

	ingestor := psxscraper.NewPSXIngestor(cfg.PSXBaseURL)
	portfolioSvc := service.NewPortfolioService(ingestor)
	dividendSvc := service.NewDividendService()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"status":      "healthy",
			"engine":      "psx-engine-go v1.0",
			"environment": cfg.Environment,
			"psx_base":    cfg.PSXBaseURL,
			"time":        time.Now().Format(time.RFC3339),
		})
	})

	r.Route("/api/v1", func(r chi.Router) {
		// Market Endpoints
		r.Get("/market/summary", func(w http.ResponseWriter, r *http.Request) {
			summary, err := ingestor.FetchMarketSummary()
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to fetch market summary: %v", err), http.StatusBadGateway)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(summary)
		})

		r.Get("/market/search", func(w http.ResponseWriter, r *http.Request) {
			query := r.URL.Query().Get("q")
			matches := ingestor.SearchCompanies(query)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(matches)
		})

		r.Get("/market/quote/{symbol}", func(w http.ResponseWriter, r *http.Request) {
			symbol := chi.URLParam(r, "symbol")
			quote, err := ingestor.FetchQuote(symbol)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to fetch quote for %s: %v", symbol, err), http.StatusBadGateway)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(quote)
		})

		// Multi-Portfolio Calculation Demo Endpoint
		r.Get("/portfolios/demo/summary", func(w http.ResponseWriter, r *http.Request) {
			demoPortfolio := models.Portfolio{
				ID:       uuid.New(),
				UserID:   uuid.New(),
				Name:     "High Dividend & Growth PSX Portfolio",
				Currency: "PKR",
			}

			demoTransactions := []models.Transaction{
				{
					Symbol:        "FFC",
					TradeType:     models.TradeBuy,
					Quantity:      decimal.NewFromInt(1000),
					PricePerShare: decimal.NewFromFloat(160.00),
					Commission:    decimal.NewFromFloat(150.00),
					TradeDate:     time.Now().AddDate(-1, 0, 0),
				},
				{
					Symbol:        "SYS",
					TradeType:     models.TradeBuy,
					Quantity:      decimal.NewFromInt(500),
					PricePerShare: decimal.NewFromFloat(420.00),
					Commission:    decimal.NewFromFloat(200.00),
					TradeDate:     time.Now().AddDate(0, -6, 0),
				},
				{
					Symbol:        "ENGRO",
					TradeType:     models.TradeBuy,
					Quantity:      decimal.NewFromInt(800),
					PricePerShare: decimal.NewFromFloat(310.00),
					Commission:    decimal.NewFromFloat(250.00),
					TradeDate:     time.Now().AddDate(0, -3, 0),
				},
			}

			summary := portfolioSvc.CalculatePortfolioSummary(demoPortfolio, demoTransactions)

			announcements := []models.DividendAnnouncement{
				{
					Symbol:          "FFC",
					PayoutPerShare:  decimal.NewFromFloat(15.50),
					BookClosureDate: time.Now().AddDate(0, 0, 10),
					FinancialPeriod: "Q3 2026",
				},
				{
					Symbol:          "ENGRO",
					PayoutPerShare:  decimal.NewFromFloat(11.00),
					BookClosureDate: time.Now().AddDate(0, 0, 25),
					FinancialPeriod: "Q3 2026",
				},
			}

			filerForecast := dividendSvc.CalculateDividendForecast(demoPortfolio.ID, models.TaxFiler, summary.Holdings, announcements)
			nonFilerForecast := dividendSvc.CalculateDividendForecast(demoPortfolio.ID, models.TaxNonFiler, summary.Holdings, announcements)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"portfolio":            summary,
				"dividend_filer_15":    filerForecast,
				"dividend_nonfiler_30": nonFilerForecast,
			})
		})
	})

	// replaces the ListenAndServe block at the end of main()
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}
	// ctx is cancelled when SIGTERM (Kubernetes) or SIGINT (Ctrl+C) arrives
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("🚀 PSX Engine Go server starting on port :%s [%s env] ...", cfg.Port, cfg.Environment)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-ctx.Done() // block until signal
	log.Println("Shutdown signal received, draining connections...")

	// give the load balancer time to stop routing to this pod
	time.Sleep(5 * time.Second)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Forced shutdown: %v", err)
	}
	log.Println("Server stopped")
}
