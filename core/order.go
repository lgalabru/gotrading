package core

import "errors"

// OrderTransactionType describes the transaction type: Bid / Ask
type OrderTransactionType uint

const (
	// Bid - we are buying the base of a currency pair, or selling the quote
	Bid OrderTransactionType = iota
	// Ask - we are selling the base of a currency pair, or buying the quote
	Ask
)

// Order represents an order
type Order struct {
	Pair             CurrencyPair         `json:"pair"`
	Price            float64              `json:"price"`
	QuoteToBasePrice float64              `json:"quoteToBasePrice"`
	BaseVolume       float64              `json:"baseVolume"`
	QuoteVolume      float64              `json:"quoteVolume"`
	TransactionType  OrderTransactionType `json:"transactionType"`
}

// InitAsk initialize an Order, setting the transactionType to Ask
func (o *Order) InitAsk(pair CurrencyPair, price float64, baseVolume float64) {
	o.TransactionType = Ask
	o.Init(pair, price, baseVolume)
}

// InitBid initialize an Order, setting the transactionType to Bid
func (o *Order) InitBid(pair CurrencyPair, price float64, baseVolume float64) {
	o.TransactionType = Bid
	o.Init(pair, price, baseVolume)
}

// Init initialize an Order
func (o *Order) Init(pair CurrencyPair, price float64, baseVolume float64) {
	o.Pair = pair
	o.Price = price
	o.BaseVolume = baseVolume
	o.QuoteToBasePrice = 1 / price
	o.QuoteVolume = o.Price * o.BaseVolume
}

// CreateMatchingAsk returns an Ask order matching the current Bid (crossing ths spread)
func (o *Order) CreateMatchingAsk() (*Order, error) {
	if o.TransactionType != Bid {
		return nil, errors.New("order: not a bid")
	}
	m := *o
	m.TransactionType = Ask
	return &m, nil
}

// CreateMatchingBid returns a Bid order matching the current Ask (crossing ths spread)
func (o *Order) CreateMatchingBid() (*Order, error) {
	if o.TransactionType != Ask {
		return nil, errors.New("order: not a bid")
	}
	m := *o
	m.TransactionType = Bid
	return &m, nil
}
