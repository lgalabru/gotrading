package strategies

import (
	"fmt"
	"strings"

	"gotrading/core"
)

type Arbitrage struct {
	Solutions []Solution
}

type node struct {
	from      core.Currency
	to        core.Currency
	exch      core.Exchange
	orderbook core.Orderbook
	inverted  bool
}

type path struct {
	nodes []node
}

func (arbitrage *Arbitrage) Run(mashup core.ExchangeMashup, depth int, threshold float64, startingCurrency core.Currency) {
	var paths []path
	from := startingCurrency
	for _, to := range mashup.Currencies {
		if to != from {
			for _, exch := range mashup.Exchanges {
				n := nodeFromOrderbook(from, to, exch, mashup)
				if n != nil {
					paths = append(findPaths(mashup, depth, path{[]node{*n}}), paths...)
				}
			}
		}
	}

	// Pour chaque noeud, on regarde les bids / asks du orderbook
	for _, p := range paths {
		performance := float64(1)
		for _, n := range p.nodes {
			var factor float64
			if n.inverted {
				factor = 1 / n.orderbook.Order.Price
			} else {
				factor = n.orderbook.Order.Price
			}
			performance = performance * factor

			// if n.inverted {
			// } else {
			// }
		}

		fmt.Println(p.description(), performance)
	}

	fmt.Println(len(paths))
}

func findPaths(m core.ExchangeMashup, depth int, p path) []path {
	var paths []path
	lastNode := p.nodes[len(p.nodes)-1]
	if len(p.nodes) == depth {
		var from core.Currency
		var to core.Currency
		if lastNode.inverted {
			to = lastNode.from
		} else {
			to = lastNode.to
		}
		if p.nodes[0].inverted {
			from = p.nodes[0].to
		} else {
			from = p.nodes[0].from
		}
		if to == from {
			return []path{p}
		}
	} else if len(p.nodes) < depth {
		var from core.Currency
		if lastNode.inverted {
			from = lastNode.from
		} else {
			from = lastNode.to
		}
		for _, to := range m.Currencies {
			if to != from {
				for _, exch := range m.Exchanges {
					n := nodeFromOrderbook(from, to, exch, m)
					if n != nil {
						if p.contains(*n) == false && len(p.nodes) < depth {
							r := findPaths(m, depth, path{append(p.nodes, *n)})
							paths = append(r, paths...)
						}
					}
				}
			}
		}
	}
	return paths
}

func nodeFromOrderbook(from core.Currency, to core.Currency, exchange core.Exchange, mashup core.ExchangeMashup) *node {
	var n *node = nil
	o := mashup.GetOrderbook(from, to, exchange)
	if o != nil {
		n = &node{from, to, exchange, *o, false}
	} else {
		o = mashup.GetOrderbook(to, from, exchange)
		if o != nil {
			n = &node{to, from, exchange, *o, true}
		}
	}
	return n
}

func (p path) contains(n node) bool {
	found := false
	for _, m := range p.nodes {
		found = n.isEqual(m)
	}
	return found
}

func (n node) display() {
	fmt.Println(n.description())
}

func (n node) isEqual(m node) bool {
	f := (strings.Compare(string(n.from), string(m.from)) == 0)
	t := (strings.Compare(string(n.to), string(m.to)) == 0)
	fi := (strings.Compare(string(n.to), string(m.from)) == 0)
	ti := (strings.Compare(string(n.from), string(m.to)) == 0)
	e := (strings.Compare(n.exch.Name, m.exch.Name) == 0)
	return f && t && e || fi && ti && e
}

func (n node) description() string {
	var str string
	if n.inverted {
		str = string(n.from) + " <- " + string(n.to) + " (" + n.exch.Name + ")"
	} else {
		str = string(n.from) + " -> " + string(n.to) + " (" + n.exch.Name + ")"
	}
	return str
}

func (p path) description() string {
	str := ""
	for _, n := range p.nodes {
		str += n.description() + " /"
	}
	return str
}

func (p path) display() {
	fmt.Println(p.description())
}
