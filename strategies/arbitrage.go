package strategies

import (
	"fmt"
	"strconv"
	"strings"

	"gotrading/graph"
)

type Arbitrage struct {
	Solutions []Solution
}

func (arbitrage *Arbitrage) Run(paths []graph.Path) {

	// Pour chaque noeud, on regarde les bids / asks du orderbook
	for _, p := range paths {
		factors := make([]string, 0)
		performance := float64(1)
		for _, n := range p.ContextualNodes {
			var factor = float64(0)
			if n.Node.Orderbook != nil {
				if n.Inverted {
					// We want to sell, so we match the Bid.
					if len(n.Node.Orderbook.Bids) > 0 {
						factor = 1 / n.Node.Orderbook.Bids[0].Price
					}
				} else {
					// We want to buy, so we match the Ask.
					if len(n.Node.Orderbook.Asks) > 0 {
						factor = n.Node.Orderbook.Asks[0].Price
					}
				}
			}
			factors = append(factors, strconv.FormatFloat(factor, 'f', 6, 64))
			performance = performance * factor
			// fmt.Println(n.Description(), strconv.FormatFloat(factor, 'f', 6, 64))
		}
		if performance > 0 {
			fmt.Println(p.Description(), performance, "//", strings.Join(factors, ", "))
		}
	}
}
