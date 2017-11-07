package graph

import (
	"gotrading/core"
)

func PathFinder(mashup core.ExchangeMashup, from core.Currency, to core.Currency, depth int) []Path {
	var paths []Path
	lookup := make(map[string]*Node)
	cLookup := make(map[string]*ContextualNode)

	for _, to := range mashup.Currencies {
		if to != from {
			for _, exch := range mashup.Exchanges {
				n, lookup, cLookup := nodeFromMashup(from, to, exch, mashup, lookup, cLookup)
				if n != nil {
					n.SoldCurrency = &from
					n.BoughtCurrency = &to
					paths = append(findPaths(mashup, depth, Path{[]*ContextualNode{n}}, lookup, cLookup), paths...)
				}
			}
		}
	}
	return paths
}

func findPaths(m core.ExchangeMashup, depth int, p Path, lookup map[string]*Node, cLookup map[string]*ContextualNode) []Path {
	var paths []Path
	lastNode := p.ContextualNodes[len(p.ContextualNodes)-1]
	if len(p.ContextualNodes) == depth {
		from := p.ContextualNodes[0].SoldCurrency
		to := lastNode.BoughtCurrency
		if *from == *to {
			return []Path{p}
		}
	} else if len(p.ContextualNodes) < depth {
		from := lastNode.BoughtCurrency
		for _, to := range m.Currencies {
			if to != *from {
				for _, exch := range m.Exchanges {
					n, lookup, cLookup := nodeFromMashup(*from, to, exch, m, lookup, cLookup)
					if n != nil {
						n.BoughtCurrency = &to
						n.SoldCurrency = from
						if p.contains(*n) == false && len(p.ContextualNodes) < depth {
							r := findPaths(m, depth, Path{append(p.ContextualNodes, n)}, lookup, cLookup)
							paths = append(r, paths...)
						}
					}
				}
			}
		}
	}
	return paths
}

func nodeFromMashup(from core.Currency, to core.Currency, exchange core.Exchange, mashup core.ExchangeMashup, lookup map[string]*Node, cLookup map[string]*ContextualNode) (*ContextualNode, map[string]*Node, map[string]*ContextualNode) {
	var cn *ContextualNode = nil
	ok := mashup.LinkExist(from, to, exchange)
	if ok {
		proto := Node{from, to, exchange, nil}
		node, ok := lookup[proto.ID()]
		if !ok {
			lookup[proto.ID()] = &proto
			node = &proto
		}
		cproto := ContextualNode{node, false, nil, nil}
		cnode, ok := cLookup[cproto.ID()]
		if !ok {
			cLookup[cproto.ID()] = &cproto
			cnode = &cproto
		}
		cn = cnode
	} else {
		ok := mashup.LinkExist(to, from, exchange)
		if ok {
			proto := Node{to, from, exchange, nil}
			node, ok := lookup[proto.ID()]
			if !ok {
				lookup[proto.ID()] = &proto
				node = &proto
			}
			cproto := ContextualNode{node, true, nil, nil}
			cnode, ok := cLookup[cproto.ID()]
			if !ok {
				cLookup[cproto.ID()] = &cproto
				cnode = &cproto
			}
			cn = cnode
		}
	}
	return cn, lookup, cLookup
}
