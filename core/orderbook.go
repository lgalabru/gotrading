package core

import (
	"time"
)

type Orderbook struct {
	CurrencyPair CurrencyPair
	Bids         []Order
	Asks         []Order
	LastUpdated  time.Time
}
