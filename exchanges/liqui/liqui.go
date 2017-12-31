package liqui

import (
	"fmt"
	"gotrading/core"
	"net/http"
)

type Liqui struct {
}

func (b Liqui) GetOrderbook() func(client http.Client, pair core.CurrencyPair) (core.Orderbook, error) {
	return func(client http.Client, pair core.CurrencyPair) (core.Orderbook, error) {
		var ob core.Orderbook
		var err error
		fmt.Println("Getting Orderbooks from Liqui")
		return ob, err
	}
}

func (b Liqui) GetPortfolio() func(client http.Client) (core.Portfolio, error) {
	return func(client http.Client) (core.Portfolio, error) {
		var p core.Portfolio
		var err error
		fmt.Println("Getting Portfolio from Liqui")
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
