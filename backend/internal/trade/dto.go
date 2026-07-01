package trade

type TradeRequest struct {
	ShareID  string `json:"share_id" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,gt=0"`
}

type ShareResponse struct {
	ID              string  `json:"id"`
	Symbol          string  `json:"symbol"`
	Name            string  `json:"name"`
	Price           float64 `json:"price"`
	Segment         string  `json:"segment"`
	AvailableShares int     `json:"available_shares"`
}
