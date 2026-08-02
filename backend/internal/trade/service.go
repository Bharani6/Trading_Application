package trade

import (
	"errors"
	"fmt"
	"time"

	market_service "stock-trading/internal/market/service"
	"stock-trading/internal/user"
	"stock-trading/internal/wallet"
	walletpkg "stock-trading/internal/wallet"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TradeService interface {
	GetAllShares(search string) ([]ShareResponse, error)
	BuyShare(userID string, req TradeRequest, isPending bool) error
	SellShare(userID string, req TradeRequest, isPending bool) error
	GetUserTrades(userID string) ([]Trade, error)
	ExecutePendingTrades() error
	CancelTrade(userID string, tradeID string) error
}

type tradeService struct {
	repo              TradeRepository
	walletRepo        wallet.WalletRepository
	marketDataService market_service.MarketDataService
	userRepo          user.UserRepository
}

func NewTradeService(repo TradeRepository, walletRepo wallet.WalletRepository, marketDataService market_service.MarketDataService) TradeService {
	return &tradeService{
		repo:              repo,
		walletRepo:        walletRepo,
		marketDataService: marketDataService,
		userRepo:          user.NewUserRepository(),
	}
}

func IsMarketOpen() bool {
	lLoc, lErr := time.LoadLocation("Asia/Kolkata")
	if lErr != nil {
		lLoc = time.FixedZone("IST", 5*3600+1800)
	}
	lNow := time.Now().In(lLoc)

	if lNow.Weekday() == time.Saturday || lNow.Weekday() == time.Sunday {
		return false
	}

	lHour := lNow.Hour()
	lMin := lNow.Minute()

	// Open between 9:00 and 15:30
	if (lHour > 9 || (lHour == 9 && lMin >= 0)) && (lHour < 15 || (lHour == 15 && lMin <= 30)) {
		return true
	}
	return false
}

func (pService *tradeService) GetAllShares(pSearch string) ([]ShareResponse, error) {
	lShares, lErr := pService.repo.GetAllShares(pSearch)
	if lErr != nil {
		return nil, errors.New("failed to retrieve shares")
	}

	// Trigger dynamic market search if a search query is provided
	if pSearch != "" {
		if lResults, lErr := pService.marketDataService.SearchSymbol(pSearch); lErr == nil && len(lResults) > 0 {
			var lNewSymbols []string
			for _, lRes := range lResults {
				lNewSymbols = append(lNewSymbols, lRes.Symbol)
			}
			// Fetch prices for the newly discovered symbols
			lPrices, lErr := pService.marketDataService.GetLatestPrices(lNewSymbols)
			if lErr == nil {
				// Initialize the segments
				lNseSeg := &Segment{Name: "NSE"}
				pService.repo.FirstOrCreateSegment(lNseSeg)
				
				lNasdaqSeg := &Segment{Name: "NASDAQ"}
				pService.repo.FirstOrCreateSegment(lNasdaqSeg)
				
				lBseSeg := &Segment{Name: "BSE"}
				pService.repo.FirstOrCreateSegment(lBseSeg)
				
				for _, lRes := range lResults {
					lPriceData := lPrices[lRes.Symbol]
					
					var lSegID uint
					var lSeg Segment
					// Simple heuristic for segment
					if len(lRes.Symbol) > 3 && lRes.Symbol[len(lRes.Symbol)-3:] == ".NS" {
						lSegID = lNseSeg.ID
						lSeg = *lNseSeg
					} else if len(lRes.Symbol) > 3 && lRes.Symbol[len(lRes.Symbol)-3:] == ".BO" {
						lSegID = lBseSeg.ID
						lSeg = *lBseSeg
					} else {
						lSegID = lNasdaqSeg.ID
						lSeg = *lNasdaqSeg
					}
					
					lNewShare := &Share{
						ID:              uuid.New(),
						Symbol:          lRes.Symbol,
						Name:            lRes.LongName,
						Price:           lPriceData.Current,
						PreviousPrice:   lPriceData.Previous,
						SegmentID:       lSegID,
						TotalShares:     1000000,
						AvailableShares: 1000000,
					}
					// If LongName is empty, fallback to ShortName
					if lNewShare.Name == "" {
						lNewShare.Name = lRes.ShortName
					}
					
					// FirstOrCreateShare returns the DB record if it exists
					if lErr := pService.repo.FirstOrCreateShare(lNewShare); lErr == nil {
						// Only append if it wasn't already in our local search results
						lAlreadyExists := false
						for _, lExisting := range lShares {
							if lExisting.Symbol == lNewShare.Symbol {
								lAlreadyExists = true
								break
							}
						}
						
						if !lAlreadyExists {
							lNewShare.Segment = lSeg
							lShares = append(lShares, *lNewShare)
						}
					}
				}
			}
		}
	}

	var lRes []ShareResponse
	for _, lShare := range lShares {
		lRes = append(lRes, ShareResponse{
			ID:              lShare.ID.String(),
			Symbol:          lShare.Symbol,
			Name:            lShare.Name,
			Price:           lShare.Price,
			PreviousPrice:   lShare.PreviousPrice,
			Segment:         lShare.Segment.Name,
			AvailableShares: lShare.AvailableShares,
		})
	}
	return lRes, nil
}

func (pService *tradeService) BuyShare(pUserID string, pReq TradeRequest, pIsPending bool) error {
	lUID, lErr := uuid.Parse(pUserID)
	if lErr != nil {
		return errors.New("invalid user id")
	}

	lU, lErr := pService.userRepo.GetUserByID(pUserID)
	if lErr != nil {
		return errors.New("invalid user")
	}
	if lU.Status == "closure_requested" {
		return errors.New("account closure requested, trading not permitted")
	}

	return pService.repo.RunInTransaction(func(pTx *gorm.DB) error {
		// 1. Lock Share row
		lShare, lTxErr := pService.repo.GetShareForUpdate(pTx, pReq.ShareID)
		if lTxErr != nil {
			return errors.New("share not found")
		}

		if lShare.AvailableShares < pReq.Quantity {
			return errors.New("insufficient shares available in market")
		}

		// 2. Lock Wallet row
		lWallet, lTxErr := pService.walletRepo.GetWalletForUpdate(pTx, pUserID)
		if lTxErr != nil {
			if errors.Is(lTxErr, gorm.ErrRecordNotFound) {
				return errors.New("insufficient wallet balance")
			}
			return errors.New("wallet not found")
		}

		lTotalCost := lShare.Price * float64(pReq.Quantity)
		if lWallet.AvailableBalance < lTotalCost {
			return errors.New("insufficient wallet balance")
		}

		// 3. Deduct from wallet
		if pIsPending {
			lWallet.AvailableBalance -= lTotalCost
			lWallet.BlockedBalance += lTotalCost
		} else {
			lWallet.WalletBalance -= lTotalCost
			lWallet.AvailableBalance -= lTotalCost
		}
		if lTxErr := pService.walletRepo.UpdateWalletWithVersion(pTx, lWallet); lTxErr != nil {
			return lTxErr
		}

		// 4. Update share count (only if not pending)
		if !pIsPending {
			lShare.AvailableShares -= pReq.Quantity
			if lTxErr := pService.repo.UpdateShareWithVersion(pTx, lShare); lTxErr != nil {
				return lTxErr
			}
		}

		lStatus := "completed"
		if pIsPending {
			lStatus = "pending"
		}

		// 5. Record Transaction
		lTransaction := &walletpkg.Transaction{
			ID:          uuid.New(),
			UserID:      lUID,
			Type:        "trade_buy",
			Amount:      lTotalCost,
			ReferenceID: fmt.Sprintf("BUY-%d", time.Now().UnixNano()),
			Description: fmt.Sprintf("Bought %d shares of %s", pReq.Quantity, lShare.Symbol),
			Status:      lStatus,
		}
		if lTxErr := pService.walletRepo.CreateTransaction(pTx, lTransaction); lTxErr != nil {
			return lTxErr
		}

		// 6. Record Trade
		lTrade := &Trade{
			ID:       uuid.New(),
			UserID:   lUID,
			ShareID:  lShare.ID,
			Quantity: pReq.Quantity,
			Price:    lShare.Price,
			Type:     "buy",
			Status:   lStatus,
		}
		return pService.repo.CreateTrade(pTx, lTrade)
	})
}

func (pService *tradeService) SellShare(pUserID string, pReq TradeRequest, pIsPending bool) error {
	lUID, lErr := uuid.Parse(pUserID)
	if lErr != nil {
		return errors.New("invalid user id")
	}

	lU, lErr := pService.userRepo.GetUserByID(pUserID)
	if lErr != nil {
		return errors.New("invalid user")
	}
	if lU.Status == "closure_requested" {
		return errors.New("account closure requested, trading not permitted")
	}

	return pService.repo.RunInTransaction(func(pTx *gorm.DB) error {
		// Basic check: Does user own enough of this share?
		// In a real app we'd have a UserPortfolio table holding owned shares.
		// For brevity, we assume the user owns it or we deduce from Trade history.
		// Let's implement a quick aggregation to check owned quantity.
		var lOwnedQty int64
		// SUM(buy) - SUM(sell completed) - SUM(sell pending)
		type Result struct {
			Total int64
		}
		var lBought Result
		pTx.Model(&Trade{}).Select("COALESCE(SUM(quantity), 0) as total").Where("user_id = ? AND share_id = ? AND type = 'buy' AND status = 'completed'", pUserID, pReq.ShareID).Scan(&lBought)

		var lSold Result
		pTx.Model(&Trade{}).Select("COALESCE(SUM(quantity), 0) as total").Where("user_id = ? AND share_id = ? AND type = 'sell' AND (status = 'completed' OR status = 'pending')", pUserID, pReq.ShareID).Scan(&lSold)

		lOwnedQty = lBought.Total - lSold.Total
		if lOwnedQty < int64(pReq.Quantity) {
			return errors.New("insufficient owned shares (some may be locked in pending orders)")
		}

		// 1. Lock Share row
		lShare, lTxErr := pService.repo.GetShareForUpdate(pTx, pReq.ShareID)
		if lTxErr != nil {
			return errors.New("share not found")
		}

		// 2. Lock Wallet row
		lWallet, lTxErr := pService.walletRepo.GetWalletForUpdate(pTx, pUserID)
		if lTxErr != nil {
			return errors.New("wallet not found")
		}

		lTotalGain := lShare.Price * float64(pReq.Quantity)

		// 3. Add to wallet (only if completed)
		if !pIsPending {
			lWallet.WalletBalance += lTotalGain
			lWallet.AvailableBalance += lTotalGain
			if lTxErr := pService.walletRepo.UpdateWalletWithVersion(pTx, lWallet); lTxErr != nil {
				return lTxErr
			}
		}

		// 4. Update share market count (only if completed)
		if !pIsPending {
			lShare.AvailableShares += pReq.Quantity
			if lTxErr := pService.repo.UpdateShareWithVersion(pTx, lShare); lTxErr != nil {
				return lTxErr
			}
		}

		lStatus := "completed"
		if pIsPending {
			lStatus = "pending"
		}

		// 5. Record Transaction
		lTransaction := &walletpkg.Transaction{
			ID:          uuid.New(),
			UserID:      lUID,
			Type:        "trade_sell",
			Amount:      lTotalGain,
			ReferenceID: fmt.Sprintf("SELL-%d", time.Now().UnixNano()),
			Description: fmt.Sprintf("Sold %d shares of %s", pReq.Quantity, lShare.Symbol),
			Status:      lStatus,
		}
		if lTxErr := pService.walletRepo.CreateTransaction(pTx, lTransaction); lTxErr != nil {
			return lTxErr
		}

		// 6. Record Trade
		lTrade := &Trade{
			ID:       uuid.New(),
			UserID:   lUID,
			ShareID:  lShare.ID,
			Quantity: pReq.Quantity,
			Price:    lShare.Price,
			Type:     "sell",
			Status:   lStatus,
		}
		return pService.repo.CreateTrade(pTx, lTrade)
	})
}

func (pService *tradeService) GetUserTrades(pUserID string) ([]Trade, error) {
	return pService.repo.GetTradesByUser(pUserID)
}

func (pService *tradeService) ExecutePendingTrades() error {
	if !IsMarketOpen() {
		return nil
	}

	return pService.repo.RunInTransaction(func(pTx *gorm.DB) error {
		var lPendingTrades []Trade
		if lTxErr := pTx.Where("status = ?", "pending").Find(&lPendingTrades).Error; lTxErr != nil {
			return lTxErr
		}

		for _, lTrade := range lPendingTrades {
			// Lock Share
			lShare, lErr := pService.repo.GetShareForUpdate(pTx, lTrade.ShareID.String())
			if lErr != nil {
				continue
			}

			// Lock Wallet
			lWallet, lErr := pService.walletRepo.GetWalletForUpdate(pTx, lTrade.UserID.String())
			if lErr != nil {
				continue
			}

			if lTrade.Type == "buy" {
				lTotalCost := lTrade.Price * float64(lTrade.Quantity)
				// Re-verify market has shares
				if lShare.AvailableShares < lTrade.Quantity {
					continue // Still not enough shares, leave as pending or fail
				}

				// Finalize wallet deduction
				lWallet.BlockedBalance -= lTotalCost
				lWallet.WalletBalance -= lTotalCost

				// Update market shares
				lShare.AvailableShares -= lTrade.Quantity
			} else if lTrade.Type == "sell" {
				lTotalGain := lTrade.Price * float64(lTrade.Quantity)

				// Finalize wallet addition
				lWallet.WalletBalance += lTotalGain
				lWallet.AvailableBalance += lTotalGain

				// Update market shares
				lShare.AvailableShares += lTrade.Quantity
			}

			// Update records
			if lErr := pService.walletRepo.UpdateWalletWithVersion(pTx, lWallet); lErr != nil {
				continue
			}
			if lErr := pService.repo.UpdateShareWithVersion(pTx, lShare); lErr != nil {
				continue
			}

			// Update Trade and Transaction statuses
			lTrade.Status = "completed"
			lTrade.UpdatedAt = time.Now()
			if lErr := pTx.Save(&lTrade).Error; lErr != nil {
				continue
			}

			// Find corresponding transaction and complete it
			var lTransaction walletpkg.Transaction
			if lErr := pTx.Where("user_id = ? AND status = ? AND type = ?", lTrade.UserID, "pending", "trade_"+lTrade.Type).First(&lTransaction).Error; lErr == nil {
				lTransaction.Status = "completed"
				pTx.Save(&lTransaction)
			}
		}

		return nil
	})
}

func (pService *tradeService) CancelTrade(pUserID string, pTradeID string) error {
	lTID, lErr := uuid.Parse(pTradeID)
	if lErr != nil {
		return errors.New("invalid trade ID")
	}
	lUID, lErr := uuid.Parse(pUserID)
	if lErr != nil {
		return errors.New("invalid user ID")
	}

	return pService.repo.RunInTransaction(func(pTx *gorm.DB) error {
		var lTrade Trade
		if lTxErr := pTx.Where("id = ?", lTID).First(&lTrade).Error; lTxErr != nil {
			return errors.New("trade not found")
		}

		if lTrade.UserID != lUID {
			return errors.New("unauthorized")
		}

		if lTrade.Status != "pending" {
			return errors.New("only pending trades can be cancelled")
		}

		if lTrade.Type == "buy" {
			lWallet, lTxErr := pService.walletRepo.GetWalletForUpdate(pTx, pUserID)
			if lTxErr != nil {
				return errors.New("wallet not found")
			}
			lRefundAmount := lTrade.Price * float64(lTrade.Quantity)
			lWallet.BlockedBalance -= lRefundAmount
			lWallet.AvailableBalance += lRefundAmount
			if lTxErr := pService.walletRepo.UpdateWalletWithVersion(pTx, lWallet); lTxErr != nil {
				return lTxErr
			}
		}

		lTrade.Status = "cancelled"
		lTrade.UpdatedAt = time.Now()
		if lTxErr := pTx.Save(&lTrade).Error; lTxErr != nil {
			return lTxErr
		}

		var lTransaction walletpkg.Transaction
		lTxType := "trade_buy"
		if lTrade.Type == "sell" {
			lTxType = "trade_sell"
		}

		lTotalAmount := lTrade.Price * float64(lTrade.Quantity)
		if lTxErr := pTx.Where("user_id = ? AND type = ? AND status = 'pending' AND amount = ?", lUID, lTxType, lTotalAmount).Order("created_at desc").First(&lTransaction).Error; lTxErr == nil {
			lTransaction.Status = "cancelled"
			pTx.Save(&lTransaction)
		}

		return nil
	})
}
