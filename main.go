package main

import (
	"fmt"
	"log"

	"gotrading/core"
	"gotrading/strategies"

	"github.com/thrasher-/gocryptotrader/config"
)

func main() {

	cfg := &config.Cfg
	err := cfg.LoadConfig("config.dat")
	if err != nil {
		log.Fatal(err)
	}

	exchanges := []core.Exchange{
		core.Exchange{"Alpha"},
		core.Exchange{"Beta"},
		core.Exchange{"Charlie"}} //SetupExchanges(*cfg)

	fmt.Println(exchanges)
	currencies := []core.Currency{"USD", "BTC", "ETH", "ETC"} //strings.Split(cfg.Cryptocurrencies, ",")

	mashup := core.ExchangeMashup{}
	mashup.Init(currencies, exchanges)

	mashup.Orderbooks[0][1][0] = &core.Orderbook{"1"}
	mashup.Orderbooks[1][2][0] = &core.Orderbook{"2"}
	mashup.Orderbooks[2][3][0] = &core.Orderbook{"3"}
	mashup.Orderbooks[3][0][0] = &core.Orderbook{"4"}
	mashup.Orderbooks[2][0][0] = &core.Orderbook{"5"}

	// for _, exch := range exchanges {
	// 	exch.Watch(currencyPairs, func(orderbook core.Orderbook) {
	// 		mashup.UpdateOrderbook(exchange, orderbook)
	// 		fmt.Println(results)
	// 	})
	// }

	depth := uint(3)
	threshold := 1.0
	currency := core.Currency("BTC")
	strategy := strategies.Arbitrage{}
	strategy.Run(mashup, depth, threshold, currency)
}
