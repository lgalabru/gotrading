package core

// CurrencyPair represents a pair of currencies.
type CurrencyPair struct {
	Base  Currency `json:"base"`
	Quote Currency `json:"quote"`
}

type CurrencyPairSettings struct {
	DecimalPlaces int     `json:"decimal_places"`
	MinPrice      float64 `json:"min_price"`
	MaxPrice      float64 `json:"max_price"`
	MinAmount     float64 `json:"min_amount"`
	MaxAmount     float64 `json:"max_amount"`
	MinTotal      float64 `json:"min_total"`
	Fee           float64 `json:"fee"`
}
