package core

import (
	"gotrading/exchanges/liqui"
	"gotrading/exchanges/orderbook"
	"gotrading/exchanges/ticker"

	"github.com/thrasher-/gocryptotrader/config"
	"github.com/thrasher-/gocryptotrader/currency/pair"
)

// Exchange represents an exchange and list the available pairs.
type Exchange struct {
	Name                     string             `json:"name"`
	AvailablePairs           []CurrencyPair     `json:"-"`
	Engine                   *ExchangeInterface `json:"-"`
	Liqui                    *liqui.Liqui       `json:"-"`
	IsCurrencyPairNormalized bool               `json:"-"`
}

// ExchangeInterface is an abstraction for using the engines from gocryptotrader.
type ExchangeInterface interface {
	Setup(exch config.ExchangeConfig)
	Start()
	SetDefaults()
	GetName() string
	IsEnabled() bool
	GetTickerPrice(currency pair.CurrencyPair, assetType string) (ticker.Price, error)
	UpdateTicker(currency pair.CurrencyPair, assetType string) (ticker.Price, error)
	GetOrderbookEx(currency pair.CurrencyPair, assetType string) (orderbook.Base, error)
	UpdateOrderbook(currency pair.CurrencyPair, assetType string) (orderbook.Base, error)
	GetEnabledCurrencies() []pair.CurrencyPair
	GetAuthenticatedAPISupport() bool
	GetAvailableCurrencies() []pair.CurrencyPair
	Trade(pair, orderType string, amount, price float64) (float64, error)
}
