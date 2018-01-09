package core

import (
	"strings"
	"time"
)

type Exchange struct {
	Name             string           `json:"name"`
	PairsEnabled     []CurrencyPair   `json:"-"`
	ExchangeSettings ExchangeSettings `json:"-"`

	FuncGetSettings  func() (ExchangeSettings, error)                            `json:"-"`
	FuncGetOrderbook func(hit Hit) (Orderbook, error)                            `json:"-"`
	FuncGetPortfolio func(settings ExchangeSettings) (Portfolio, error)          `json:"-"`
	FuncPostOrder    func(order Order, settings ExchangeSettings) (Order, error) `json:"-"`
	// fnDeposit      func(client http.Client) (bool, error)
	// fnWithdraw     func(client http.Client) (bool, error)
}

type ExchangeSettings struct {
	Name                     string                                `json:"-"`
	APIKey                   string                                `json:"-"`
	APISecret                string                                `json:"-"`
	AvailablePairs           []CurrencyPair                        `json:"-"`
	PairsSettings            map[CurrencyPair]CurrencyPairSettings `json:"-"`
	IsCurrencyPairNormalized bool                                  `json:"-"`
	Nonce                    Nonce                                 `json:"-"`
}

func (e *Exchange) LoadSettings() {
	settings, err := e.FuncGetSettings()
	if err == nil {
		nonce := Nonce{}
		nonce.Set(time.Now().Unix())
		settings.Nonce = nonce
		settings.Name = e.Name
		e.ExchangeSettings = settings
	}
}

func (e *Exchange) GetOrderbook(hit Hit) (Orderbook, error) {
	return e.FuncGetOrderbook(hit)
}

func (e *Exchange) GetPortfolio() (Portfolio, error) {
	return e.FuncGetPortfolio(e.ExchangeSettings)
}

func (e *Exchange) PostOrder(order Order) (Order, error) {
	return e.FuncPostOrder(order, e.ExchangeSettings)
}

func (e *Exchange) LoadPairsEnabled(joinedPairs string) {
	pairs := strings.Split(joinedPairs, ",")
	e.PairsEnabled = []CurrencyPair{}
	for _, pair := range pairs {
		c := strings.Split(pair, "_")
		e.PairsEnabled = append(
			e.PairsEnabled,
			CurrencyPair{Currency(c[0]), Currency(c[1])})
	}
}
