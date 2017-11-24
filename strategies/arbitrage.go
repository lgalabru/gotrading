package strategies

import (
	"math"

	"gotrading/core"
	"gotrading/graph"
)

type Arbitrage struct {
	Solutions []Solution
}

type ArbitrageChain struct {
	Path            graph.Path
	OrdersToFulfill []core.Order
	Performance     float64 `json:"performance"`
	Volume          float64 `json:"volume"`
}

func (arbitrage *Arbitrage) Run(paths []graph.Path) []ArbitrageChain {

	chains := make([]ArbitrageChain, len(paths))

	for j, p := range paths {
		initialCurrencyToLastCurrency := float64(1)
		chain := ArbitrageChain{}
		chain.OrdersToFulfill = make([]core.Order, len(p.Nodes))
		for i, n := range p.Nodes {
			var priceOfCurrencyToSell float64
			var volumeOfCurrencyToSell float64
			var order core.Order
			if n.Endpoint.Orderbook != nil {
				if n.IsBaseToQuote {
					// We want to sell the base, so we match the Bid.
					if len(n.Endpoint.Orderbook.Bids) > 0 {
						bestBid := n.Endpoint.Orderbook.Bids[0]
						o, err := bestBid.CreateMatchingAsk()
						if err == nil {
							order = *o
							priceOfCurrencyToSell = order.Price
							volumeOfCurrencyToSell = order.BaseVolume
						}
					}
				} else {
					// We want to sell the quote, so we match the Ask.
					if len(n.Endpoint.Orderbook.Asks) > 0 {
						bestAsk := n.Endpoint.Orderbook.Asks[0]
						o, err := bestAsk.CreateMatchingBid()
						if err == nil {
							order = *o
							priceOfCurrencyToSell = order.PriceOfQuoteToBase
							volumeOfCurrencyToSell = order.QuoteVolume
						}
					}
				}
			}
			initialCurrencyToLastCurrency = initialCurrencyToLastCurrency * priceOfCurrencyToSell
			if i == 0 {
				chain.Volume = volumeOfCurrencyToSell
			} else {
				result := math.Min(
					chain.Volume*initialCurrencyToLastCurrency,
					volumeOfCurrencyToSell*priceOfCurrencyToSell)

				chain.Volume = result / initialCurrencyToLastCurrency
			}
			chain.OrdersToFulfill[i] = order
		}
		chain.Path = p
		chain.Performance = initialCurrencyToLastCurrency
		chains[j] = chain
	}

	return chains
}
