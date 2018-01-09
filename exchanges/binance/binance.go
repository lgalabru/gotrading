package binance

import (
	"fmt"
	"gotrading/core"
)

type Binance struct {
}

func (b Binance) GetSettings() func() (core.ExchangeSettings, error) {
	return func() (core.ExchangeSettings, error) {
		var settings core.ExchangeSettings
		var err error
		fmt.Println("Loading settings from Binance")
		return settings, err

	}
}

func (b Binance) GetOrderbook() func(hit core.Hit) (core.Orderbook, error) {
	return func(hit core.Hit) (core.Orderbook, error) {
		var ob core.Orderbook
		var err error
		fmt.Println("Getting Orderbooks from Binance")
		return ob, err
	}
}

func (b Binance) GetPortfolio() func(settings core.ExchangeSettings) (core.Portfolio, error) {
	return func(settings core.ExchangeSettings) (core.Portfolio, error) {
		var p core.Portfolio
		var err error
		fmt.Println("Getting Portfolio from Binance")
		return p, err
	}
}

func (b Binance) PostOrder() func(order core.Order, settings core.ExchangeSettings) (core.Order, error) {
	return func(order core.Order, settings core.ExchangeSettings) (core.Order, error) {
		var o core.Order
		var err error
		fmt.Println("Posting Order on Binance")
		return o, err
	}
}

// func (b *Binance) Deposit(client http.Client) (bool, error) {
// 	var err error
// 	return true, err
// }

// func (b *Binance) Withdraw(client http.Client) (bool, error) {
// 	var err error
// 	return true, err
// }
