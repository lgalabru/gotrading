package binance

import (
	"fmt"
	"gotrading/core"
	"net/http"
)

type Binance struct {
}

func (b Binance) GetOrderbook() func(client http.Client, pair core.CurrencyPair) (core.Orderbook, error) {
	return func(client http.Client, pair core.CurrencyPair) (core.Orderbook, error) {
		var ob core.Orderbook
		var err error
		fmt.Println("Getting Orderbooks from Binance")
		return ob, err
	}
}

func (b Binance) GetPortfolio() func(client http.Client) (core.Portfolio, error) {
	return func(client http.Client) (core.Portfolio, error) {
		var p core.Portfolio
		var err error
		fmt.Println("Getting Portfolio from Binance")
		return p, err
	}
}

// func (b *Binance) PostOrder(client http.Client, order core.Order) (core.Order, error) {
// 	var o core.Order
// 	var err error
// 	return o, err
// }

// func (b *Binance) Deposit(client http.Client) (bool, error) {
// 	var err error
// 	return true, err
// }

// func (b *Binance) Withdraw(client http.Client) (bool, error) {
// 	var err error
// 	return true, err
// }
