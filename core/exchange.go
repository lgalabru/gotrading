package core

import (
	"strings"
)

type Exchange struct {
	Name             string           `json:"name"`
	PairsEnabled     []CurrencyPair   `json:"-"`
	ExchangeSettings ExchangeSettings `json:"-"`

	FuncGetSettings  func() (ExchangeSettings, error) `json:"-"`
	FuncGetOrderbook func(hit Hit) (Orderbook, error) `json:"-"`
	FuncGetPortfolio func() (Portfolio, error)        `json:"-"`
	FuncPostOrder    func(order Order) (Order, error) `json:"-"`
	// fnDeposit      func(client http.Client) (bool, error)
	// fnWithdraw     func(client http.Client) (bool, error)
}

type ExchangeSettings struct {
	AvailablePairs           []CurrencyPair                        `json:"-"`
	PairsSettings            map[CurrencyPair]CurrencyPairSettings `json:"-"`
	IsCurrencyPairNormalized bool                                  `json:"-"`
}

func (b *Exchange) LoadSettings() {
	settings, err := b.FuncGetSettings()
	if err == nil {
		b.ExchangeSettings = settings
	}
}

func (b *Exchange) GetOrderbook(hit Hit) (Orderbook, error) {
	return b.FuncGetOrderbook(hit)
}

func (b *Exchange) GetPortfolio() (Portfolio, error) {
	return b.FuncGetPortfolio()
}

func (b *Exchange) PostOrder(order Order) (Order, error) {
	return b.FuncPostOrder(order)
}

func (b *Exchange) LoadPairsEnabled(joinedPairs string) {
	pairs := strings.Split(joinedPairs, ",")
	b.PairsEnabled = []CurrencyPair{}
	for _, pair := range pairs {
		c := strings.Split(pair, "_")
		b.PairsEnabled = append(
			b.PairsEnabled,
			CurrencyPair{Currency(c[0]), Currency(c[1])})
	}
}
