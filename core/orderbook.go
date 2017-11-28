package core

import "time"

// Orderbook represents an orderbook
type Orderbook struct {
	CurrencyPair CurrencyPair `json:"pair"`
	Bids         []Order      `json:"bids"`
	Asks         []Order      `json:"asks"`
	LastUpdate   time.Time    `json:"lastUpdate"`
}
