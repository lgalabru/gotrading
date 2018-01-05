package exchanges

import (
	"fmt"
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
	GetSettings() func() (core.ExchangeSettings, error)
	GetOrderbook() func(hit core.Hit) (core.Orderbook, error)
	GetPortfolio() func() (core.Portfolio, error)
	PostOrder() func(order core.Order) (core.Order, error)
	// Deposit(client http.Client) (bool, error)
	// Withdraw(client http.Client) (bool, error)
}

func (f *Factory) BuildExchange(name string) core.Exchange {
	key := strings.Join([]string{"exchanges", name}, ".")
	config := viper.GetStringMapString(key)
	fmt.Println("Building", name, config)

	exchange := core.Exchange{}
	exchange.LoadPairsEnabled(config["pairs_enabled"])
	switch name {
	case "Binance":
		injectStandardizedMethods(&exchange, binance.Binance{})
	case "Liqui":
		injectStandardizedMethods(&exchange, liqui.Liqui{})
	default:
	}
	return exchange
}

func injectStandardizedMethods(b *core.Exchange, exch standardizedExchange) {
	b.FuncGetSettings = exch.GetSettings()
	b.FuncGetOrderbook = exch.GetOrderbook()
	b.FuncGetPortfolio = exch.GetPortfolio()
	b.FuncPostOrder = exch.PostOrder()
	// b.fnDeposit = exch.Deposit
	// b.fnWithdraw = exch.Withdraw
}
