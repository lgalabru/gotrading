package services

import (
	"fmt"
	"time"

	"gotrading/core"
	"gotrading/graph"

	"github.com/thrasher-/gocryptotrader/config"
	"github.com/thrasher-/gocryptotrader/currency/pair"
)

type orderbookReceived func(endpoint graph.Endpoint)

func LoadExchange(cfg *config.Config, name string, exch core.ExchangeInterface) core.Exchange {
	config, _ := cfg.GetExchangeConfig(name)
	exch.SetDefaults()
	exch.Setup(config)
	var rawPairs = exch.GetAvailableCurrencies()
	pairs := make([]core.CurrencyPair, len(rawPairs))
	for i, c := range rawPairs {
		pairs[i] = core.CurrencyPair{
			core.Currency(c.GetFirstCurrency()),
			core.Currency(c.GetSecondCurrency())}
	}
	return core.Exchange{name, pairs, &exch}
}

func StartPollingOrderbooks(exch core.Exchange, nodes []graph.EndpointLookup, delayBetweenReqs time.Duration, fn orderbookReceived) {
	var i = int(0)
	for {
		n := nodes[i%len(nodes)]
		i += 1
		time.Sleep(delayBetweenReqs * time.Millisecond)

		cp := pair.NewCurrencyPair(string(n.Endpoint.From), string(n.Endpoint.To))
		src, err := (*exch.Engine).UpdateOrderbook(cp, "SPOT")
		if err == nil {
			dst := &core.Orderbook{}
			dst.CurrencyPair = core.CurrencyPair{n.Endpoint.From, n.Endpoint.To}
			dst.Bids = make([]core.Order, 0)
			dst.Asks = make([]core.Order, 0)
			dst.LastUpdate = time.Now()
			// fmt.Println("1 ------------------")
			// fmt.Println(src.Asks)
			for _, ask := range src.Asks {
				// if exch.Name == "Poloniex" {
				// 	dst.Bids = append(dst.Bids, core.Order{1 / ask.Price, ask.Amount, core.Buy})
				// } else {
				dst.Asks = append(dst.Asks, core.NewAsk(dst.CurrencyPair, ask.Price, ask.Amount))
				// }
			}
			// fmt.Println("2 ------------------")
			// fmt.Println(src.Bids)
			for _, bid := range src.Bids {
				// if exch.Name == "Poloniex" {
				// 	dst.Asks = append(dst.Asks, core.Order{1 / bid.Price, bid.Amount, core.Sell})
				// } else {
				dst.Bids = append(dst.Bids, core.NewBid(dst.CurrencyPair, bid.Price, bid.Amount))
				// }
			}
			n.Endpoint.Orderbook = dst

			// fmt.Println("~~~~~~~~~~~~~~~~~~")
			fn(*n.Endpoint)
		} else {
			fmt.Println("Error", n.Endpoint.Description(), err)
		}
	}
}
