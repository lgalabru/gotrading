package graph

import (
	"fmt"
	"gotrading/core"
)

func PathFinder(mashup core.ExchangeMashup, from core.Currency, to core.Currency, depth int) ([]*Endpoint, map[string][]Path, []Path) {
	var rawPaths []Path
	endpointLookup := make(map[string]*Endpoint)
	nodeLookup := make(map[string]*Node)

	for _, to := range mashup.Currencies {
		if to != from {
			for _, exch := range mashup.Exchanges {
				n, endpointLookup, nodeLookup := nodeFromMashup(from, to, exch, mashup, endpointLookup, nodeLookup)
				if n != nil {
					n.SoldCurrency = &from
					n.BoughtCurrency = &to
					rawPaths = append(findPaths(mashup, depth, Path{[]*Node{n}, nil, nil}, endpointLookup, nodeLookup), rawPaths...)
				}
			}
		}
	}

	endpoints := make([]*Endpoint, 0)
	paths := make(map[string][]Path)
	for _, path := range rawPaths {
		for _, n := range path.Nodes {
			p, ok := paths[n.Endpoint.ID()]
			if !ok {
				endpoints = append(endpoints, n.Endpoint)
				p = make([]Path, 0)
			}
			path.encode()
			paths[n.Endpoint.ID()] = append(p, path)
		}
	}

	fmt.Println("Observing", len(rawPaths), "paths")

	return endpoints, paths, rawPaths
}

func findPaths(m core.ExchangeMashup, depth int, p Path, endpointLookup map[string]*Endpoint, nodeLookup map[string]*Node) []Path {
	var paths []Path
	// fmt.Println(p.Description())
	lastNode := p.Nodes[len(p.Nodes)-1]
	if len(p.Nodes) == depth {
		from := p.Nodes[0].SoldCurrency
		to := lastNode.BoughtCurrency
		if *from == *to {
			return []Path{p}
		}
	} else if len(p.Nodes) < depth {
		from := lastNode.BoughtCurrency
		for _, to := range m.Currencies {
			if to != *from {
				for _, exch := range m.Exchanges {
					n, endpointLookup, nodeLookup := nodeFromMashup(*from, to, exch, m, endpointLookup, nodeLookup)
					if n != nil {
						n.BoughtCurrency = &to
						n.SoldCurrency = from
						if p.contains(*n) == false && len(p.Nodes) < depth {
							r := findPaths(m, depth, Path{append(p.Nodes, n), nil, nil}, endpointLookup, nodeLookup)
							paths = append(r, paths...)
						}
					}
				}
			}
		}
	}
	return paths
}

func nodeFromMashup(from core.Currency, to core.Currency, exchange core.Exchange, mashup core.ExchangeMashup, endpointLookup map[string]*Endpoint, nodeLookup map[string]*Node) (*Node, map[string]*Endpoint, map[string]*Node) {
	var n *Node = nil
	ok := mashup.LinkExist(from, to, exchange)
	if ok {
		proto := Endpoint{from, to, exchange, nil}
		endpoint, ok := endpointLookup[proto.ID()]
		if !ok {
			endpointLookup[proto.ID()] = &proto
			endpoint = &proto
		}
		cproto := Node{endpoint, true, nil, nil}
		node, ok := nodeLookup[cproto.ID()]
		if !ok {
			nodeLookup[cproto.ID()] = &cproto
			node = &cproto
		}
		n = node
	} else {
		ok := mashup.LinkExist(to, from, exchange)
		if ok {
			proto := Endpoint{to, from, exchange, nil}
			endpoint, ok := endpointLookup[proto.ID()]
			if !ok {
				endpointLookup[proto.ID()] = &proto
				endpoint = &proto
			}
			cproto := Node{endpoint, false, nil, nil}
			node, ok := nodeLookup[cproto.ID()]
			if !ok {
				nodeLookup[cproto.ID()] = &cproto
				node = &cproto
			}
			n = node
		}
	}
	return n, endpointLookup, nodeLookup
}
