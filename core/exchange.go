package core

import (
	"strings"
)

type Exchange struct {
	Name                     string         `json:"name"`
	AvailablePairs           []CurrencyPair `json:"-"`
	PairsEnabled             []CurrencyPair `json:"-"`
	IsCurrencyPairNormalized bool           `json:"-"`

	FuncGetOrderbook func(hit Hit) (Orderbook, error)
	FuncGetPortfolio func() (Portfolio, error)
	FuncPostOrder    func(order Order) (Order, error)
	// fnDeposit      func(client http.Client) (bool, error)
	// fnWithdraw     func(client http.Client) (bool, error)
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

func (b *Exchange) LoadAvailablePairs(joinedPairs string) {
	pairs := strings.Split(joinedPairs, ",")
	b.AvailablePairs = []CurrencyPair{}
	for _, pair := range pairs {
		c := strings.Split(pair, "_")
		b.AvailablePairs = append(
			b.AvailablePairs,
			CurrencyPair{Currency(c[0]), Currency(c[1])})
	}
}
