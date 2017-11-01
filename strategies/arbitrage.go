package strategies

import (
	"fmt"

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

func (arbitrage *Arbitrage) Run(mashup core.ExchangeMashup, depth uint, threshold float64, startingCurrency core.Currency) {
	var paths []path
	from := startingCurrency
	for to, _ := range mashup.Currencies {
		if to != from {
			for exch, _ := range mashup.Exchanges {
				n := nodeFromOrderbook(from, to, exch, mashup)
				if n != nil {
					paths = append(findPaths(mashup, depth-1, path{[]node{*n}}), paths...)
				}
			}
		}
	}
	for _, p := range paths {
		p.display()
	}
	fmt.Println(len(paths))
}

func findPaths(m core.ExchangeMashup, depthLeft uint, p path) []path {
	var paths []path
	if depthLeft == 0 {
		return []path{p}
	} else {
		from := p.nodes[len(p.nodes)-1].to
		recursion := func(from core.Currency, to core.Currency) {
			for exch, _ := range m.Exchanges {
				n := nodeFromOrderbook(from, to, exch, m)
				if n != nil && p.contains(*n) == false {
					p.nodes = append(p.nodes, *n)
					r := findPaths(m, depthLeft-1, p)
					paths = append(r, paths...)
				}
			}
		}
		if depthLeft == 1 {
			to := p.nodes[0].from
			recursion(from, to)
		} else {
			for to, _ := range m.Currencies {
				if to != from {
					recursion(from, to)
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
		found = (n == m)
	}
	return found
}

func (n node) display() {
	fmt.Println(n.description())
}

func (n node) description() string {
	return string(n.from) + "_" + string(n.to) + " (" + n.exch.Name + " - " + n.orderbook.Id + ")"
}

func (p path) display() {
	str := ""
	for _, n := range p.nodes {
		str += n.description() + " /"
	}
	fmt.Println(str)
}
