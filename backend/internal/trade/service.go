package trade

import (
	"errors"
	"fmt"
	"time"

	"stock-trading/internal/wallet"
	walletpkg "stock-trading/internal/wallet"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TradeService interface {
	GetAllShares() ([]ShareResponse, error)
	BuyShare(userID string, req TradeRequest, isPending bool) error
	SellShare(userID string, req TradeRequest, isPending bool) error
	GetUserTrades(userID string) ([]Trade, error)
	ExecutePendingTrades() error
	CancelTrade(userID string, tradeID string) error
}

type tradeService struct {
	repo       TradeRepository
	walletRepo wallet.WalletRepository
}

func NewTradeService() TradeService {
	return &tradeService{
		repo:       NewTradeRepository(),
		walletRepo: wallet.NewWalletRepository(),
	}
}

func IsMarketOpen() bool {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		loc = time.FixedZone("IST", 5*3600+1800)
	}
	now := time.Now().In(loc)

	if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
		return false
	}

	hour := now.Hour()
	min := now.Minute()

	// Open between 9:00 and 15:30
	if hour < 9 || hour > 15 || (hour == 15 && min > 30) {
		return false
	}

	return true
}

func (s *tradeService) GetAllShares() ([]ShareResponse, error) {
	shares, err := s.repo.GetAllShares()
	if err != nil {
		return nil, errors.New("failed to retrieve shares")
	}

	var res []ShareResponse
	for _, share := range shares {
		res = append(res, ShareResponse{
			ID:              share.ID.String(),
			Symbol:          share.Symbol,
			Name:            share.Name,
			Price:           share.Price,
			PreviousPrice:   share.PreviousPrice,
			Segment:         share.Segment.Name,
			AvailableShares: share.AvailableShares,
		})
	}
	return res, nil
}

func (s *tradeService) BuyShare(userID string, req TradeRequest, isPending bool) error {
	uID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user id")
	}

	return s.repo.RunInTransaction(func(tx *gorm.DB) error {
		// 1. Lock Share row
		share, err := s.repo.GetShareForUpdate(tx, req.ShareID)
		if err != nil {
			return errors.New("share not found")
		}

		if share.AvailableShares < req.Quantity {
			return errors.New("insufficient shares available in market")
		}

		// 2. Lock Wallet row
		wallet, err := s.walletRepo.GetWalletForUpdate(tx, userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("insufficient wallet balance")
			}
			return errors.New("wallet not found")
		}

		totalCost := share.Price * float64(req.Quantity)
		if wallet.AvailableBalance < totalCost {
			return errors.New("insufficient wallet balance")
		}

		// 3. Deduct from wallet
		if isPending {
			wallet.AvailableBalance -= totalCost
			wallet.BlockedBalance += totalCost
		} else {
			wallet.WalletBalance -= totalCost
			wallet.AvailableBalance -= totalCost
		}
		if err := s.walletRepo.UpdateWallet(tx, wallet); err != nil {
			return err
		}

		// 4. Update share count (only if not pending)
		if !isPending {
			share.AvailableShares -= req.Quantity
			if err := s.repo.UpdateShare(tx, share); err != nil {
				return err
			}
		}

		status := "completed"
		if isPending {
			status = "pending"
		}

		// 5. Record Transaction
		transaction := &walletpkg.Transaction{
			ID:          uuid.New(),
			UserID:      uID,
			Type:        "trade_buy",
			Amount:      totalCost,
			ReferenceID: fmt.Sprintf("BUY-%d", time.Now().UnixNano()),
			Description: fmt.Sprintf("Bought %d shares of %s", req.Quantity, share.Symbol),
			Status:      status,
		}
		if err := s.walletRepo.CreateTransaction(tx, transaction); err != nil {
			return err
		}

		// 6. Record Trade
		trade := &Trade{
			ID:       uuid.New(),
			UserID:   uID,
			ShareID:  share.ID,
			Quantity: req.Quantity,
			Price:    share.Price,
			Type:     "buy",
			Status:   status,
		}
		return s.repo.CreateTrade(tx, trade)
	})
}

func (s *tradeService) SellShare(userID string, req TradeRequest, isPending bool) error {
	uID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user id")
	}

	return s.repo.RunInTransaction(func(tx *gorm.DB) error {
		// Basic check: Does user own enough of this share?
		// In a real app we'd have a UserPortfolio table holding owned shares.
		// For brevity, we assume the user owns it or we deduce from Trade history.
		// Let's implement a quick aggregation to check owned quantity.
		var ownedQty int64
		// SUM(buy) - SUM(sell completed) - SUM(sell pending)
		type Result struct {
			Total int64
		}
		var bought Result
		tx.Model(&Trade{}).Select("COALESCE(SUM(quantity), 0) as total").Where("user_id = ? AND share_id = ? AND type = 'buy' AND status = 'completed'", userID, req.ShareID).Scan(&bought)

		var sold Result
		tx.Model(&Trade{}).Select("COALESCE(SUM(quantity), 0) as total").Where("user_id = ? AND share_id = ? AND type = 'sell' AND (status = 'completed' OR status = 'pending')", userID, req.ShareID).Scan(&sold)

		ownedQty = bought.Total - sold.Total
		if ownedQty < int64(req.Quantity) {
			return errors.New("insufficient owned shares (some may be locked in pending orders)")
		}

		// 1. Lock Share row
		share, err := s.repo.GetShareForUpdate(tx, req.ShareID)
		if err != nil {
			return errors.New("share not found")
		}

		// 2. Lock Wallet row
		wallet, err := s.walletRepo.GetWalletForUpdate(tx, userID)
		if err != nil {
			return errors.New("wallet not found")
		}

		totalGain := share.Price * float64(req.Quantity)

		// 3. Add to wallet (only if completed)
		if !isPending {
			wallet.WalletBalance += totalGain
			wallet.AvailableBalance += totalGain
			if err := s.walletRepo.UpdateWallet(tx, wallet); err != nil {
				return err
			}
		}

		// 4. Update share market count (only if completed)
		if !isPending {
			share.AvailableShares += req.Quantity
			if err := s.repo.UpdateShare(tx, share); err != nil {
				return err
			}
		}

		status := "completed"
		if isPending {
			status = "pending"
		}

		// 5. Record Transaction
		transaction := &walletpkg.Transaction{
			ID:          uuid.New(),
			UserID:      uID,
			Type:        "trade_sell",
			Amount:      totalGain,
			ReferenceID: fmt.Sprintf("SELL-%d", time.Now().UnixNano()),
			Description: fmt.Sprintf("Sold %d shares of %s", req.Quantity, share.Symbol),
			Status:      status,
		}
		if err := s.walletRepo.CreateTransaction(tx, transaction); err != nil {
			return err
		}

		// 6. Record Trade
		trade := &Trade{
			ID:       uuid.New(),
			UserID:   uID,
			ShareID:  share.ID,
			Quantity: req.Quantity,
			Price:    share.Price,
			Type:     "sell",
			Status:   status,
		}
		return s.repo.CreateTrade(tx, trade)
	})
}

func (s *tradeService) GetUserTrades(userID string) ([]Trade, error) {
	return s.repo.GetTradesByUser(userID)
}

func (s *tradeService) ExecutePendingTrades() error {
	if !IsMarketOpen() {
		return nil
	}

	return s.repo.RunInTransaction(func(tx *gorm.DB) error {
		var pendingTrades []Trade
		if err := tx.Where("status = ?", "pending").Find(&pendingTrades).Error; err != nil {
			return err
		}

		for _, trade := range pendingTrades {
			// Lock Share
			share, err := s.repo.GetShareForUpdate(tx, trade.ShareID.String())
			if err != nil {
				continue
			}

			// Lock Wallet
			wallet, err := s.walletRepo.GetWalletForUpdate(tx, trade.UserID.String())
			if err != nil {
				continue
			}

			if trade.Type == "buy" {
				totalCost := trade.Price * float64(trade.Quantity)
				// Re-verify market has shares
				if share.AvailableShares < trade.Quantity {
					continue // Still not enough shares, leave as pending or fail
				}

				// Finalize wallet deduction
				wallet.BlockedBalance -= totalCost
				wallet.WalletBalance -= totalCost

				// Update market shares
				share.AvailableShares -= trade.Quantity
			} else if trade.Type == "sell" {
				totalGain := trade.Price * float64(trade.Quantity)

				// Finalize wallet addition
				wallet.WalletBalance += totalGain
				wallet.AvailableBalance += totalGain

				// Update market shares
				share.AvailableShares += trade.Quantity
			}

			// Update records
			if err := s.walletRepo.UpdateWallet(tx, wallet); err != nil {
				continue
			}
			if err := s.repo.UpdateShare(tx, share); err != nil {
				continue
			}

			// Update Trade and Transaction statuses
			trade.Status = "completed"
			trade.UpdatedAt = time.Now()
			if err := tx.Save(&trade).Error; err != nil {
				continue
			}

			// Find corresponding transaction and complete it
			var transaction walletpkg.Transaction
			if err := tx.Where("user_id = ? AND status = ? AND type = ?", trade.UserID, "pending", "trade_"+trade.Type).First(&transaction).Error; err == nil {
				transaction.Status = "completed"
				tx.Save(&transaction)
			}
		}

		return nil
	})
}

func (s *tradeService) CancelTrade(userID string, tradeID string) error {
	tID, err := uuid.Parse(tradeID)
	if err != nil {
		return errors.New("invalid trade ID")
	}
	uID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	return s.repo.RunInTransaction(func(tx *gorm.DB) error {
		var trade Trade
		if err := tx.Where("id = ?", tID).First(&trade).Error; err != nil {
			return errors.New("trade not found")
		}

		if trade.UserID != uID {
			return errors.New("unauthorized")
		}

		if trade.Status != "pending" {
			return errors.New("only pending trades can be cancelled")
		}

		if trade.Type == "buy" {
			wallet, err := s.walletRepo.GetWalletForUpdate(tx, userID)
			if err != nil {
				return errors.New("wallet not found")
			}
			refundAmount := trade.Price * float64(trade.Quantity)
			wallet.BlockedBalance -= refundAmount
			wallet.AvailableBalance += refundAmount
			if err := s.walletRepo.UpdateWallet(tx, wallet); err != nil {
				return err
			}
		}

		trade.Status = "cancelled"
		trade.UpdatedAt = time.Now()
		if err := tx.Save(&trade).Error; err != nil {
			return err
		}

		var transaction walletpkg.Transaction
		txType := "trade_buy"
		if trade.Type == "sell" {
			txType = "trade_sell"
		}

		totalAmount := trade.Price * float64(trade.Quantity)
		if err := tx.Where("user_id = ? AND type = ? AND status = 'pending' AND amount = ?", uID, txType, totalAmount).Order("created_at desc").First(&transaction).Error; err == nil {
			transaction.Status = "cancelled"
			tx.Save(&transaction)
		}

		return nil
	})
}
