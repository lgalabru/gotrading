package core

// Orderbook represents an orderbook
type Orderbook struct {
	CurrencyPair CurrencyPair
	Bids         []Order
	Asks         []Order
	// 	LastUpdated  time.Time
}
