package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type MarketPrice struct {
	Current  float64
	Previous float64
}

// MarketDataService defines the interface for fetching stock market prices
type MarketDataService interface {
	GetLatestPrices(symbols []string) (map[string]MarketPrice, error)
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
func (s *YahooFinanceService) GetLatestPrices(symbols []string) (map[string]MarketPrice, error) {
	prices := make(map[string]MarketPrice)
	if len(symbols) == 0 {
		return prices, nil
	}

	for _, symbol := range symbols {
		url := fmt.Sprintf("https://query2.finance.yahoo.com/v8/finance/chart/%s?interval=1d&range=1d", symbol)
		
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			continue // skip on error
		}
		
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
		
		resp, err := s.client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			if resp != nil && resp.Body != nil {
				resp.Body.Close()
			}
			continue
		}

		var result struct {
			Chart struct {
				Result []struct {
					Meta struct {
						RegularMarketPrice         float64 `json:"regularMarketPrice"`
						RegularMarketPreviousClose float64 `json:"chartPreviousClose"`
					} `json:"meta"`
				} `json:"result"`
			} `json:"chart"`
		}

		err = json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()
		
		if err == nil && len(result.Chart.Result) > 0 {
			meta := result.Chart.Result[0].Meta
			prices[symbol] = MarketPrice{
				Current:  meta.RegularMarketPrice,
				Previous: meta.RegularMarketPreviousClose,
			}
		}
	}

	return prices, nil
}
