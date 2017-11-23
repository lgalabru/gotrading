package core

// CurrencyPair represents a pair of currencies.
type CurrencyPair struct {
	Base  Currency `json:"base"`
	Quote Currency `json:"quote"`
}
