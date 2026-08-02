package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type MarketPrice struct {
	Current  float64
	Previous float64
}

type SearchResult struct {
	Symbol    string
	ShortName string
	LongName  string
	QuoteType string
}

// MarketDataService defines the interface for fetching stock market prices
type MarketDataService interface {
	GetLatestPrices(pSymbols []string) (map[string]MarketPrice, error)
	SearchSymbol(pQuery string) ([]SearchResult, error)
}

type YahooFinanceService struct {
	client *http.Client
}

func NewYahooFinanceService() MarketDataService {
	return &YahooFinanceService{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// GetLatestPrices fetches the latest prices from Yahoo Finance
func (pService *YahooFinanceService) GetLatestPrices(pSymbols []string) (map[string]MarketPrice, error) {
	lPrices := make(map[string]MarketPrice)
	if len(pSymbols) == 0 {
		return lPrices, nil
	}

	for _, lSymbol := range pSymbols {
		lURL := fmt.Sprintf("https://query2.finance.yahoo.com/v8/finance/chart/%s?interval=1d&range=1d", lSymbol)
		
		lReq, lErr := http.NewRequest("GET", lURL, nil)
		if lErr != nil {
			continue // skip on error
		}
		
		lReq.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
		
		lResp, lErr := pService.client.Do(lReq)
		if lErr != nil || lResp.StatusCode != http.StatusOK {
			if lResp != nil && lResp.Body != nil {
				lResp.Body.Close()
			}
			continue
		}

		var lResult struct {
			Chart struct {
				Result []struct {
					Meta struct {
						RegularMarketPrice         float64 `json:"regularMarketPrice"`
						RegularMarketPreviousClose float64 `json:"chartPreviousClose"`
					} `json:"meta"`
				} `json:"result"`
			} `json:"chart"`
		}

		lErr = json.NewDecoder(lResp.Body).Decode(&lResult)
		lResp.Body.Close()
		
		if lErr == nil && len(lResult.Chart.Result) > 0 {
			lMeta := lResult.Chart.Result[0].Meta
			lPrices[lSymbol] = MarketPrice{
				Current:  lMeta.RegularMarketPrice,
				Previous: lMeta.RegularMarketPreviousClose,
			}
		}
	}

	return lPrices, nil
}

func (pService *YahooFinanceService) SearchSymbol(pQuery string) ([]SearchResult, error) {
	var lResults []SearchResult
	if pQuery == "" {
		return lResults, nil
	}

	lUrlStr := fmt.Sprintf("https://query2.finance.yahoo.com/v1/finance/search?q=%s", url.QueryEscape(pQuery))
	lReq, lErr := http.NewRequest("GET", lUrlStr, nil)
	if lErr != nil {
		return nil, lErr
	}
	lReq.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")

	lResp, lErr := pService.client.Do(lReq)
	if lErr != nil {
		return nil, lErr
	}
	defer lResp.Body.Close()

	if lResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("yahoo finance search returned status %d", lResp.StatusCode)
	}

	var lSearchResp struct {
		Quotes []struct {
			Symbol    string `json:"symbol"`
			Shortname string `json:"shortname"`
			Longname  string `json:"longname"`
			QuoteType string `json:"quoteType"`
		} `json:"quotes"`
	}

	if lErr := json.NewDecoder(lResp.Body).Decode(&lSearchResp); lErr != nil {
		return nil, lErr
	}

	for _, lQ := range lSearchResp.Quotes {
		if lQ.QuoteType == "EQUITY" && (strings.HasSuffix(lQ.Symbol, ".NS") || strings.HasSuffix(lQ.Symbol, ".BO")) {
			lResults = append(lResults, SearchResult{
				Symbol:    lQ.Symbol,
				ShortName: lQ.Shortname,
				LongName:  lQ.Longname,
				QuoteType: lQ.QuoteType,
			})
		}
	}

	return lResults, nil
}
