package exchanges

import (
	"fmt"
	"net/http"
	"strings"

	"gotrading/core"
	"gotrading/exchanges/binance"
	"gotrading/exchanges/liqui"

	"github.com/spf13/viper"
)

type Factory struct {
}

// standardizedExchange enforces standard functions for all supported exchanges
type standardizedExchange interface {
	GetOrderbook() func(client http.Client, pair core.CurrencyPair) (core.Orderbook, error)
	GetPortfolio() func(client http.Client) (core.Portfolio, error)
	// PostOrder(client http.Client, order core.Order) (core.Order, error)
	// Deposit(client http.Client) (bool, error)
	// Withdraw(client http.Client) (bool, error)
}

func (f *Factory) BuildExchange(name string) core.ExchangeBase {
	key := strings.Join([]string{"exchanges", name}, ".")
	config := viper.GetStringMapString(key)
	fmt.Println("Building", name, config)

	exchange := core.ExchangeBase{}
	switch name {
	case "Binance":
		injectStandardizedMethods(&exchange, binance.Binance{})
	case "Liqui":
		injectStandardizedMethods(&exchange, liqui.Liqui{})

	default:
	}
	return exchange
}

func injectStandardizedMethods(b *core.ExchangeBase, exch standardizedExchange) {
	b.FuncGetOrderbook = exch.GetOrderbook()
	b.FuncGetPortfolio = exch.GetPortfolio()
	// b.fnPostOrder = exch.PostOrder
	// b.fnDeposit = exch.Deposit
	// b.fnWithdraw = exch.Withdraw
}
