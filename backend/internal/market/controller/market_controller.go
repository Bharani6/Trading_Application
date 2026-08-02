package controller

import (
	"net/http"
	"stock-trading/internal/market/service"

	"github.com/gin-gonic/gin"
)

type MarketController struct {
	marketSvc service.MarketDataService
}

func NewMarketController(pSvc service.MarketDataService) *MarketController {
	return &MarketController{marketSvc: pSvc}
}

func (pController *MarketController) GetIndices(pCtx *gin.Context) {
	lSymbols := []string{
		"^NSEI",       // NIFTY 50
		"^BSESN",      // SENSEX
		"^NSEBANK",    // NIFTY BANK
		"^NSEMDCP50",  // NIFTY MIDCAP 50
		"^CNXFIN",     // FINNIFTY (Nifty Financial Services)
	}

	lPrices, lErr := pController.marketSvc.GetLatestPrices(lSymbols)
	if lErr != nil {
		pCtx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch indices"})
		return
	}

	type IndexData struct {
		Name      string  `json:"name"`
		Current   float64 `json:"current"`
		Change    float64 `json:"change"`
		ChangePct float64 `json:"change_pct"`
	}

	var lResults []IndexData

	lMapping := map[string]string{
		"^NSEI":       "NIFTY",
		"^BSESN":      "SENSEX",
		"^NSEBANK":    "BANKNIFTY",
		"^NSEMDCP50":  "MIDCPNIFTY",
		"^CNXFIN":     "FINNIFTY",
	}

	for _, lSym := range lSymbols {
		if lP, lOk := lPrices[lSym]; lOk && lP.Current > 0 {
			lChange := lP.Current - lP.Previous
			lChangePct := 0.0
			if lP.Previous > 0 {
				lChangePct = (lChange / lP.Previous) * 100
			}
			lResults = append(lResults, IndexData{
				Name:      lMapping[lSym],
				Current:   lP.Current,
				Change:    lChange,
				ChangePct: lChangePct,
			})
		} else {
			lResults = append(lResults, IndexData{
				Name:      lMapping[lSym],
				Current:   0,
				Change:    0,
				ChangePct: 0,
			})
		}
	}

	pCtx.JSON(http.StatusOK, gin.H{"data": lResults})
}
