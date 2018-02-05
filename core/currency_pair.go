package core

// CurrencyPair represents a pair of currencies.
type CurrencyPair struct {
	Base  Currency `json:"base"`
	Quote Currency `json:"quote"`
}

type CurrencyPairSettings struct {
	BasePrecision  int     `json:"basePrecision"`
	QuotePrecision int     `json:"quotePrecision"`
	MinPrice       float64 `json:"minPrice"`
	MaxPrice       float64 `json:"maxPrice"`
	MinAmount      float64 `json:"minAmount"`
	MaxAmount      float64 `json:"maxAmount"`
	MinTotal       float64 `json:"minTotal"`
	Fee            float64 `json:"fee"`
}
