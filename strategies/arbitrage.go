package strategies

import (
	"fmt"
	"strconv"
	"strings"

	"gotrading/core"
	"gotrading/graph"
)

type Arbitrage struct {
	Solutions []Solution
}

type ArbitrageChain struct {
	Path            graph.Path
	OrdersToFulfill []core.Order
	Performance     float64
}

func (arbitrage *Arbitrage) Run(paths []graph.Path) []ArbitrageChain {

	chains := make([]ArbitrageChain, len(paths))

	for _, p := range paths {
		factors := make([]string, 0)
		performance := float64(1)
		chain := ArbitrageChain{}
		chain.OrdersToFulfill = make([]core.Order, len(p.ContextualNodes))
		for i, n := range p.ContextualNodes {
			var factor = float64(0)
			order := core.Order{0, 0, 0}
			if n.Node.Orderbook != nil {
				if n.Inverted {
					// We want to sell, so we match the Bid.
					if len(n.Node.Orderbook.Bids) > 0 {
						order = n.Node.Orderbook.Bids[0]
						factor = 1 / order.Price
					}
				} else {
					// We want to buy, so we match the Ask.
					if len(n.Node.Orderbook.Asks) > 0 {
						order = n.Node.Orderbook.Asks[0]
						factor = order.Price
					}
				}
			}
			chain.OrdersToFulfill[i] = order
			factors = append(factors, strconv.FormatFloat(factor, 'f', 6, 64))
			performance = performance * factor
			// fmt.Println(n.Description(), strconv.FormatFloat(factor, 'f', 6, 64))
		}
		chain.Path = p
		chain.Performance = performance
		chains = append(chains, chain)
		if performance > 0 {
			fmt.Println(p.Description(), performance, "//", strings.Join(factors, ", "))
		}
	}

	return chains
}
