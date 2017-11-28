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
	Path            graph.Path   `json:"path"`
	Orders          []core.Order `json:"orders"`
	Performance     float64      `json:"performance"`
	Rates           []float64    `json:"rates"`
	AdjustedVolumes []float64    `json:"volumes"`
	VolumeToEngage  float64      `json:"volumeToEngage"`
}

func (arbitrage *Arbitrage) Run(paths []graph.Path) []ArbitrageChain {

	chains := make([]ArbitrageChain, len(paths))
	// fromCurrentToInitial := float64(1)
	fromInitialToCurrent := float64(1)

	// rateForInitialCurrency := float64(1) // How many INITIAL_CURRENCY are we getting for 1 CURRENT_CURRENCY
	for j, p := range paths {
		chain := ArbitrageChain{}
		chain.Rates = make([]float64, len(p.Nodes))
		chain.AdjustedVolumes = make([]float64, len(p.Nodes))

		chain.Orders = make([]core.Order, len(p.Nodes))

		for i, n := range p.Nodes {

			var priceOfCurrencyToSell float64
			var volumeOfCurrencyToSell float64

			var order core.Order
			if n.Endpoint.Orderbook != nil {
				if n.IsBaseToQuote {
					// We are selling the base -> we match the Bid.
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
					// We are selling the quote <=> we are buying the base -> we match the Ask.
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
			fromInitialToCurrent = fromInitialToCurrent * priceOfCurrencyToSell
			chain.Rates[i] = fromInitialToCurrent
			chain.Performance = fromInitialToCurrent
			if i == 0 {
				chain.VolumeToEngage = volumeOfCurrencyToSell
			} else {
				limitingAmount := chain.VolumeToEngage * fromInitialToCurrent
				currentAmount := volumeOfCurrencyToSell * priceOfCurrencyToSell
				newLimitingAmount := math.Min(limitingAmount, currentAmount)
				chain.VolumeToEngage = newLimitingAmount / fromInitialToCurrent
			}
			chain.Orders[i] = order
		}
		for i, n := range p.Nodes {
			chain.AdjustedVolumes[i] = chain.VolumeToEngage * chain.Rates[i]
			if n.IsBaseToQuote {
				chain.Orders[i].UpdateQuoteVolume(chain.AdjustedVolumes[i])
			} else {
				chain.Orders[i].UpdateBaseVolume(chain.AdjustedVolumes[i])
			}
		}
		chain.Path = p
		chains[j] = chain
	}

	return chains
}
