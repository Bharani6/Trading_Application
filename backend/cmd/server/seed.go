package main

import (
	"stock-trading/internal/trade"
	"gorm.io/gorm"
	"go.uber.org/zap"
)

func SeedData(pDB *gorm.DB) error {
	lErr := pDB.Transaction(func(pTx *gorm.DB) error {
		// 1. Create or get segments
		lNse := trade.Segment{Name: "NSE"}
		if lTxErr := pTx.FirstOrCreate(&lNse, trade.Segment{Name: "NSE"}).Error; lTxErr != nil {
			return lTxErr
		}

		lUs := trade.Segment{Name: "NASDAQ"}
		if lTxErr := pTx.FirstOrCreate(&lUs, trade.Segment{Name: "NASDAQ"}).Error; lTxErr != nil {
			return lTxErr
		}

		// 2. Define top symbols
		lStocks := []struct {
			Symbol    string
			Name      string
			SegmentID uint
		}{
			{"RELIANCE.NS", "Reliance Industries", lNse.ID},
			{"TCS.NS", "Tata Consultancy Services", lNse.ID},
			{"HDFCBANK.NS", "HDFC Bank", lNse.ID},
			{"ICICIBANK.NS", "ICICI Bank", lNse.ID},
			{"INFY.NS", "Infosys", lNse.ID},
			{"SBIN.NS", "State Bank of India", lNse.ID},
			{"BHARTIARTL.NS", "Bharti Airtel", lNse.ID},
			{"ITC.NS", "ITC Ltd", lNse.ID},
			{"LT.NS", "Larsen & Toubro", lNse.ID},
			{"HINDUNILVR.NS", "Hindustan Unilever", lNse.ID},
			{"AXISBANK.NS", "Axis Bank", lNse.ID},
			{"KOTAKBANK.NS", "Kotak Mahindra Bank", lNse.ID},
			{"BAJFINANCE.NS", "Bajaj Finance", lNse.ID},
			{"BAJAJFINSV.NS", "Bajaj Finserv", lNse.ID},
			{"MARUTI.NS", "Maruti Suzuki", lNse.ID},
			{"M&M.NS", "Mahindra & Mahindra", lNse.ID},
			{"SUNPHARMA.NS", "Sun Pharmaceutical", lNse.ID},
			{"ULTRACEMCO.NS", "UltraTech Cement", lNse.ID},
			{"ASIANPAINT.NS", "Asian Paints", lNse.ID},
			{"TITAN.NS", "Titan Company", lNse.ID},
			{"NESTLEIND.NS", "Nestlé India", lNse.ID},
			{"POWERGRID.NS", "Power Grid Corporation", lNse.ID},
			{"NTPC.NS", "NTPC", lNse.ID},
			{"ONGC.NS", "Oil & Natural Gas Corporation", lNse.ID},
			{"COALINDIA.NS", "Coal India", lNse.ID},
			{"ADANIENT.NS", "Adani Enterprises", lNse.ID},
			{"ADANIPORTS.NS", "Adani Ports & SEZ", lNse.ID},
			{"BEL.NS", "Bharat Electronics", lNse.ID},
			{"HCLTECH.NS", "HCL Technologies", lNse.ID},
			{"WIPRO.NS", "Wipro", lNse.ID},
			{"TECHM.NS", "Tech Mahindra", lNse.ID},
			{"CIPLA.NS", "Cipla", lNse.ID},
			{"DRREDDY.NS", "Dr. Reddy's Laboratories", lNse.ID},
			{"INDUSINDBK.NS", "IndusInd Bank", lNse.ID},
			{"TATASTEEL.NS", "Tata Steel", lNse.ID},
			{"JSWSTEEL.NS", "JSW Steel", lNse.ID},
			{"ETERNAL.NS", "Eternal (formerly Zomato)", lNse.ID},
			{"TRENT.NS", "Trent", lNse.ID},
			{"SHRIRAMFIN.NS", "Shriram Finance", lNse.ID},
			{"GRASIM.NS", "Grasim Industries", lNse.ID},
			{"HINDALCO.NS", "Hindalco Industries", lNse.ID},
			{"BAJAJ-AUTO.NS", "Bajaj Auto", lNse.ID},
			{"HEROMOTOCO.NS", "Hero MotoCorp", lNse.ID},
			{"EICHERMOT.NS", "Eicher Motors", lNse.ID},
			{"APOLLOHOSP.NS", "Apollo Hospitals", lNse.ID},
			{"BRITANNIA.NS", "Britannia Industries", lNse.ID},
			{"TATACONSUM.NS", "Tata Consumer Products", lNse.ID},
			{"INDIGO.NS", "InterGlobe Aviation", lNse.ID},
			{"JIOFIN.NS", "Jio Financial Services", lNse.ID},
		}

		// 3. Create shares with 0 initial prices if they don't exist
		for _, lStock := range lStocks {
			lShare := trade.Share{
				Symbol:          lStock.Symbol,
				Name:            lStock.Name,
				Price:           0.0,
				PreviousPrice:   0.0,
				SegmentID:       lStock.SegmentID,
				TotalShares:     1000000,
				AvailableShares: 1000000,
			}
			if lTxErr := pTx.Where("symbol = ?", lStock.Symbol).FirstOrCreate(&lShare).Error; lTxErr != nil {
				return lTxErr
			}
		}

		return nil
	})

	if lErr != nil {
		zap.L().Error("Failed to seed database", zap.Error(lErr))
		return lErr
	}

	zap.L().Info("Database seeded successfully with initial stocks.")
	return nil
}
