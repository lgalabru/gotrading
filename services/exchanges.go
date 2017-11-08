package services

import (
	"fmt"
	"time"

	"gotrading/core"
	"gotrading/graph"

	"github.com/thrasher-/gocryptotrader/config"
	"github.com/thrasher-/gocryptotrader/currency/pair"
)

type orderbookReceived func(node graph.Node)

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
	return core.Exchange{name, pairs, exch}
}

func StartPollingOrderbooks(exch core.Exchange, nodes []graph.NodeLookup, delayBetweenReqs time.Duration, fn orderbookReceived) {

	for _, n := range nodes {
		time.Sleep(delayBetweenReqs * time.Millisecond)

		cp := pair.NewCurrencyPair(string(n.Node.From), string(n.Node.To))
		src, err := exch.Engine.UpdateOrderbook(cp, "SPOT")
		if err == nil {
			dst := n.Node.Orderbook
			if dst == nil {
				dst = &core.Orderbook{}
				dst.CurrencyPair = core.CurrencyPair{n.Node.From, n.Node.To}
				dst.Bids = make([]core.Order, 0)
				dst.Asks = make([]core.Order, 0)
				n.Node.Orderbook = dst
			}
			// fmt.Println("1 ------------------")
			// fmt.Println(src.Asks)
			for _, ask := range src.Asks {
				// if exch.Name == "Poloniex" {
				// 	dst.Bids = append(dst.Bids, core.Order{1 / ask.Price, ask.Amount, core.Buy})
				// } else {
				dst.Asks = append(dst.Asks, core.Order{ask.Price, ask.Amount, core.Sell})
				// }
			}
			// fmt.Println("2 ------------------")
			// fmt.Println(src.Bids)
			for _, bid := range src.Bids {
				// if exch.Name == "Poloniex" {
				// 	dst.Asks = append(dst.Asks, core.Order{1 / bid.Price, bid.Amount, core.Sell})
				// } else {
				dst.Bids = append(dst.Bids, core.Order{bid.Price, bid.Amount, core.Buy})
				// }
			}

			// fmt.Println("~~~~~~~~~~~~~~~~~~")
			fn(*n.Node)
		} else {
			fmt.Println("Error", n.Node.Description(), err)
		}
	}
}
