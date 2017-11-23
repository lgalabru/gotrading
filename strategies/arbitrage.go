package strategies

import (
	"math"
	"strconv"

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
		factors := make([]string, 0)
		performance := float64(1)
		chain := ArbitrageChain{}
		chain.OrdersToFulfill = make([]core.Order, len(p.ContextualNodes))
		for i, n := range p.ContextualNodes {
			var factor = float64(0)
			order := core.Order{0, 0, 0}
			volume := float64(0)
			if n.Node.Orderbook != nil {
				if n.Inverted {
					// We want to sell the quote, so we match the Ask.
					if len(n.Node.Orderbook.Asks) > 0 {
						order = n.Node.Orderbook.Asks[0]
						factor = 1 / order.Price
						volume = order.Volume
					}
				} else {
					// We want to buy the quote, so we match the Bid.
					if len(n.Node.Orderbook.Bids) > 0 {
						order = n.Node.Orderbook.Bids[0]
						factor = order.Price
						volume = order.Volume
					}
				}
			}
			performance = performance * factor
			if i == 0 {
				chain.Volume = volume
			} else {
				//
				result := math.Min(chain.Volume*performance, volume*order.Price)
				chain.Volume = result / performance
			}
			chain.OrdersToFulfill[i] = order
			factors = append(factors, strconv.FormatFloat(factor, 'f', 6, 64))
		}
		chain.Path = p
		chain.Performance = performance
		chains[j] = chain
	}

	return chains
}
