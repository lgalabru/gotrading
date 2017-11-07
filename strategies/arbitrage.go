package strategies

import (
	"fmt"
	"strconv"

	"gotrading/graph"
)

type Arbitrage struct {
	Solutions []Solution
}

func (arbitrage *Arbitrage) Run(paths []graph.Path) {

	// Pour chaque noeud, on regarde les bids / asks du orderbook

	for _, p := range paths {
		performance := float64(1)
		for _, n := range p.ContextualNodes {
			var factor = float64(0)
			if n.Node.Orderbook != nil {
				if n.Inverted {
					if len(n.Node.Orderbook.Asks) > 0 {
						factor = 1 / n.Node.Orderbook.Asks[0].Price
					}
				} else {
					if len(n.Node.Orderbook.Bids) > 0 {
						factor = n.Node.Orderbook.Bids[0].Price
					}
				}
			}
			performance = performance * factor
			fmt.Println(n.Description(), strconv.FormatFloat(factor, 'f', 6, 64))

		}
		fmt.Println(p.Description(), performance)
	}
}
