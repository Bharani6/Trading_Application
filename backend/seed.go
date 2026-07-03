package main

import (
	"fmt"
	"log"

	"stock-trading/internal/config"
	"stock-trading/internal/database"
	"stock-trading/internal/trade"
	"gorm.io/gorm"
)

func main() {
	config.LoadConfig("configs")
	database.ConnectDB()
	db := database.DB

	err := db.Transaction(func(tx *gorm.DB) error {
		// 1. Create or get segments
		nse := trade.Segment{Name: "NSE"}
		if err := tx.FirstOrCreate(&nse, trade.Segment{Name: "NSE"}).Error; err != nil {
			return err
		}

		us := trade.Segment{Name: "NASDAQ"}
		if err := tx.FirstOrCreate(&us, trade.Segment{Name: "NASDAQ"}).Error; err != nil {
			return err
		}

		// 2. Define top symbols
		stocks := []struct {
			Symbol    string
			Name      string
			SegmentID uint
		}{
			{"RELIANCE.NS", "Reliance Industries", nse.ID},
			{"TCS.NS", "Tata Consultancy Services", nse.ID},
			{"HDFCBANK.NS", "HDFC Bank", nse.ID},
			{"ICICIBANK.NS", "ICICI Bank", nse.ID},
			{"INFY.NS", "Infosys", nse.ID},
			{"SBIN.NS", "State Bank of India", nse.ID},
			{"BHARTIARTL.NS", "Bharti Airtel", nse.ID},
			{"ITC.NS", "ITC Ltd", nse.ID},
			{"LT.NS", "Larsen & Toubro", nse.ID},
			{"HINDUNILVR.NS", "Hindustan Unilever", nse.ID},
			{"AXISBANK.NS", "Axis Bank", nse.ID},
			{"KOTAKBANK.NS", "Kotak Mahindra Bank", nse.ID},
			{"BAJFINANCE.NS", "Bajaj Finance", nse.ID},
			{"BAJAJFINSV.NS", "Bajaj Finserv", nse.ID},
			{"MARUTI.NS", "Maruti Suzuki", nse.ID},
			{"M&M.NS", "Mahindra & Mahindra", nse.ID},
			{"TATAMOTORS.NS", "Tata Motors", nse.ID},
			{"SUNPHARMA.NS", "Sun Pharmaceutical", nse.ID},
			{"ULTRACEMCO.NS", "UltraTech Cement", nse.ID},
			{"ASIANPAINT.NS", "Asian Paints", nse.ID},
			{"TITAN.NS", "Titan Company", nse.ID},
			{"NESTLEIND.NS", "Nestlé India", nse.ID},
			{"POWERGRID.NS", "Power Grid Corporation", nse.ID},
			{"NTPC.NS", "NTPC", nse.ID},
			{"ONGC.NS", "Oil & Natural Gas Corporation", nse.ID},
			{"COALINDIA.NS", "Coal India", nse.ID},
			{"ADANIENT.NS", "Adani Enterprises", nse.ID},
			{"ADANIPORTS.NS", "Adani Ports & SEZ", nse.ID},
			{"BEL.NS", "Bharat Electronics", nse.ID},
			{"HCLTECH.NS", "HCL Technologies", nse.ID},
			{"WIPRO.NS", "Wipro", nse.ID},
			{"TECHM.NS", "Tech Mahindra", nse.ID},
			{"CIPLA.NS", "Cipla", nse.ID},
			{"DRREDDY.NS", "Dr. Reddy's Laboratories", nse.ID},
			{"INDUSINDBK.NS", "IndusInd Bank", nse.ID},
			{"TATASTEEL.NS", "Tata Steel", nse.ID},
			{"JSWSTEEL.NS", "JSW Steel", nse.ID},
			{"ETERNAL.NS", "Eternal (formerly Zomato)", nse.ID},
			{"TRENT.NS", "Trent", nse.ID},
			{"SHRIRAMFIN.NS", "Shriram Finance", nse.ID},
			{"GRASIM.NS", "Grasim Industries", nse.ID},
			{"HINDALCO.NS", "Hindalco Industries", nse.ID},
			{"BAJAJ-AUTO.NS", "Bajaj Auto", nse.ID},
			{"HEROMOTOCO.NS", "Hero MotoCorp", nse.ID},
			{"EICHERMOT.NS", "Eicher Motors", nse.ID},
			{"APOLLOHOSP.NS", "Apollo Hospitals", nse.ID},
			{"BRITANNIA.NS", "Britannia Industries", nse.ID},
			{"TATACONSUM.NS", "Tata Consumer Products", nse.ID},
			{"INDIGO.NS", "InterGlobe Aviation", nse.ID},
			{"JIOFIN.NS", "Jio Financial Services", nse.ID},
		}

		// 3. Create shares with 0 initial prices
		for _, s := range stocks {
			share := trade.Share{
				Symbol:          s.Symbol,
				Name:            s.Name,
				Price:           0.0,
				PreviousPrice:   0.0,
				SegmentID:       s.SegmentID,
				TotalShares:     1000000,
				AvailableShares: 1000000,
			}
			if err := tx.Where("symbol = ?", s.Symbol).FirstOrCreate(&share).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	fmt.Println("Successfully seeded initial stocks into the database.")
}
