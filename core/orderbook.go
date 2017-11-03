package core

type Orderbook struct {
	CurrencyPair CurrencyPair
	Order        Order
}

// type Orderbook struct {
// 	CurrencyPair CurrencyPair
// 	Bids         []Order
// 	Asks         []Order
// 	LastUpdated  time.Time
// }
