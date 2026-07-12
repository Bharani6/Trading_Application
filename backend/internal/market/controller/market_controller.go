package controller

import (
	"net/http"
	"stock-trading/internal/market/service"

	"github.com/gin-gonic/gin"
)

type MarketController struct {
	marketSvc service.MarketDataService
}

func NewMarketController(svc service.MarketDataService) *MarketController {
	return &MarketController{marketSvc: svc}
}

func (c *MarketController) GetIndices(ctx *gin.Context) {
	symbols := []string{
		"^NSEI",       // NIFTY 50
		"^BSESN",      // SENSEX
		"^NSEBANK",    // NIFTY BANK
		"^NSEMDCP50",  // NIFTY MIDCAP 50
		"^CNXFIN",     // FINNIFTY (Nifty Financial Services)
	}

	prices, err := c.marketSvc.GetLatestPrices(symbols)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch indices"})
		return
	}

	type IndexData struct {
		Name      string  `json:"name"`
		Current   float64 `json:"current"`
		Change    float64 `json:"change"`
		ChangePct float64 `json:"change_pct"`
	}

	var results []IndexData

	mapping := map[string]string{
		"^NSEI":       "NIFTY",
		"^BSESN":      "SENSEX",
		"^NSEBANK":    "BANKNIFTY",
		"^NSEMDCP50":  "MIDCPNIFTY",
		"^CNXFIN":     "FINNIFTY",
	}

	for _, sym := range symbols {
		if p, ok := prices[sym]; ok && p.Current > 0 {
			change := p.Current - p.Previous
			changePct := 0.0
			if p.Previous > 0 {
				changePct = (change / p.Previous) * 100
			}
			results = append(results, IndexData{
				Name:      mapping[sym],
				Current:   p.Current,
				Change:    change,
				ChangePct: changePct,
			})
		} else {
			results = append(results, IndexData{
				Name:      mapping[sym],
				Current:   0,
				Change:    0,
				ChangePct: 0,
			})
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"data": results})
}
