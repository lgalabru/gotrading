package core

import (
	"gotrading/exchange/liqui"
	"gotrading/exchange/orderbook"
	"gotrading/exchange/ticker"
	"net/http"

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

type ExchangeBase struct {
	Name                     string         `json:"name"`
	PairsEnabled             []CurrencyPair `json:"-"`
	IsCurrencyPairNormalized bool           `json:"-"`

	FuncGetOrderbook func(client http.Client, pair CurrencyPair) (Orderbook, error)
	FuncGetPortfolio func(client http.Client) (Portfolio, error)
	// fnPostOrder    func(client http.Client, order core.Order) (core.Order, error)
	// fnDeposit      func(client http.Client) (bool, error)
	// fnWithdraw     func(client http.Client) (bool, error)
}

func (b *ExchangeBase) GetOrderbook(client http.Client, pair CurrencyPair) (Orderbook, error) {
	return b.FuncGetOrderbook(client, pair)
}

func (b *ExchangeBase) GetPortfolio(client http.Client) (Portfolio, error) {
	return b.FuncGetPortfolio(client)
}
