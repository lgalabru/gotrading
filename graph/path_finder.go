package graph

import (
	"gotrading/core"
)

func PathFinder(mashup core.ExchangeMashup, from core.Currency, to core.Currency, depth int) []Path {
	var paths []Path
	for _, to := range mashup.Currencies {
		if to != from {
			for _, exch := range mashup.Exchanges {
				n := nodeFromMashup(from, to, exch, mashup)
				if n != nil {
					paths = append(findPaths(mashup, depth, Path{[]Node{*n}}), paths...)
				}
			}
		}
	}
	return paths
}

func findPaths(m core.ExchangeMashup, depth int, p Path) []Path {
	var paths []Path
	lastNode := p.Nodes[len(p.Nodes)-1]
	if len(p.Nodes) == depth {
		var from core.Currency
		var to core.Currency
		if lastNode.Inverted {
			to = lastNode.From
		} else {
			to = lastNode.To
		}
		if p.Nodes[0].Inverted {
			from = p.Nodes[0].To
		} else {
			from = p.Nodes[0].From
		}
		if to == from {
			return []Path{p}
		}
	} else if len(p.Nodes) < depth {
		var from core.Currency
		if lastNode.Inverted {
			from = lastNode.From
		} else {
			from = lastNode.To
		}
		for _, to := range m.Currencies {
			if to != from {
				for _, exch := range m.Exchanges {
					n := nodeFromMashup(from, to, exch, m)
					if n != nil {
						if p.contains(*n) == false && len(p.Nodes) < depth {
							r := findPaths(m, depth, Path{append(p.Nodes, *n)})
							paths = append(r, paths...)
						}
					}
				}
			}
		}
	}
	return paths
}

func nodeFromMashup(from core.Currency, to core.Currency, exchange core.Exchange, mashup core.ExchangeMashup) *Node {
	var n *Node = nil
	ok := mashup.LinkExist(from, to, exchange)
	if ok {
		n = &Node{from, to, exchange, false}
	} else {
		ok := mashup.LinkExist(to, from, exchange)
		if ok {
			n = &Node{to, from, exchange, true}
		}
	}
	return n
}
