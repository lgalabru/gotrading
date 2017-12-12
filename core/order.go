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
	Pair               CurrencyPair         `json:"pair"`
	Price              float64              `json:"price"`
	PriceOfQuoteToBase float64              `json:"quoteToBasePrice"`
	BaseVolume         float64              `json:"baseVolume"`
	QuoteVolume        float64              `json:"quoteVolume"`
	TransactionType    OrderTransactionType `json:"transactionType"`
	Fee                float64              `json:"fee"`
	TakerFee           float64              `json:"takerFee"`
	BaseVolumeIn       float64              `json:"baseVolumeIn"`
	BaseVolumeOut      float64              `json:"baseVolumeOut"`
	QuoteVolumeIn      float64              `json:"quoteVolumeIn"`
	QuoteVolumeOut     float64              `json:"quoteVolumeIn"`
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

// NewAsk initialize an Order, setting the transactionType to Ask
func NewAsk(pair CurrencyPair, price float64, baseVolume float64) Order {
	o := Order{}
	o.InitAsk(pair, price, baseVolume)
	return o
}

// NewBid returns an Order, setting the transactionType to Bid
func NewBid(pair CurrencyPair, price float64, baseVolume float64) Order {
	o := Order{}
	o.InitBid(pair, price, baseVolume)
	return o
}

// Init initialize an Order
func (o *Order) Init(pair CurrencyPair, price float64, baseVolume float64) {
	o.Pair = pair
	o.Price = price
	o.PriceOfQuoteToBase = 1 / price
	o.TakerFee = 0.25 / 100
	o.UpdateBaseVolume(baseVolume)
}

// UpdateBaseVolume cascade update on BaseVolume and QuoteVolume
func (o *Order) UpdateBaseVolume(baseVolume float64) {
	o.BaseVolume = baseVolume
	o.QuoteVolume = o.Price * o.BaseVolume
	o.Fee = o.BaseVolume * o.TakerFee
	o.updateVolumesInOut()
}

// UpdateQuoteVolume cascade update on BaseVolume and QuoteVolume
func (o *Order) UpdateQuoteVolume(quoteVolume float64) {
	o.QuoteVolume = quoteVolume
	o.BaseVolume = o.QuoteVolume / o.Price
	o.Fee = o.BaseVolume * o.TakerFee
	o.updateVolumesInOut()
}

func (o *Order) updateVolumesInOut() {
	if o.TransactionType == Bid {
		o.BaseVolumeIn = 0
		o.QuoteVolumeIn = o.QuoteVolume
		o.BaseVolumeOut = o.BaseVolume - o.BaseVolume*o.TakerFee
		o.QuoteVolumeOut = 0
	} else if o.TransactionType == Ask {
		o.BaseVolumeIn = o.BaseVolume
		o.QuoteVolumeIn = 0
		o.BaseVolumeOut = 0
		o.QuoteVolumeOut = o.QuoteVolume - o.QuoteVolume*o.TakerFee
	}
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
		return nil, errors.New("order: not a ask")
	}
	m := *o
	m.TransactionType = Bid
	return &m, nil
}
