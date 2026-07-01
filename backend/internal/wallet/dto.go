package wallet

type FundRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

type WalletResponse struct {
	WalletBalance    float64 `json:"wallet_balance"`
	BlockedBalance   float64 `json:"blocked_balance"`
	AvailableBalance float64 `json:"available_balance"`
}
