package graph

import (
	"fmt"
	"gotrading/core"
)

func PathFinder(mashup core.ExchangeMashup, from core.Currency, to core.Currency, depth int) (Tree, []*core.Endpoint, map[string][]Path, []Path) {
	var rawPaths []Path
	endpointLookup := make(map[string]*core.Endpoint)
	nodeLookup := make(map[string]*core.Hit)

	for _, to := range mashup.Currencies {
		if to != from {
			for _, exch := range mashup.Exchanges {
				n, endpointLookup, nodeLookup := nodeFromMashup(from, to, exch, mashup, endpointLookup, nodeLookup)
				if n != nil {
					rawPaths = append(findPaths(mashup, depth, Path{[]*core.Hit{n}, nil, nil, 0}, endpointLookup, nodeLookup), rawPaths...)
				}
			}
		}
	}

	// arbitrage := arbitrage.From(endpoint1).To(endpoint2).To(endpoint3).To(endpoint4)
	// result := arbitrage.Run()

	// Behind the scene, arbitrage, is going to deal with the fetching.
	// vertices := treeOfPossibles.Roots()
	// treeNodes := vertice.Children()

	treeOfPossibles := Tree{}
	endpoints := make([]*core.Endpoint, 0)
	paths := make(map[string][]Path)

	for _, path := range rawPaths {
		treeOfPossibles.InsertPath(path)

		for _, n := range path.Hits {

			p, ok := paths[n.Endpoint.ID()]
			if !ok {
				endpoints = append(endpoints, n.Endpoint)
				p = make([]Path, 0)
			}
			path.Encode()
			paths[n.Endpoint.ID()] = append(p, path)
		}
	}
	fmt.Println("Observing", len(rawPaths), "paths")
	return treeOfPossibles, endpoints, paths, rawPaths
}

func findPaths(m core.ExchangeMashup, depth int, p Path, endpointLookup map[string]*core.Endpoint, nodeLookup map[string]*core.Hit) []Path {
	var paths []Path
	firstNode := p.Hits[0]
	lastNode := p.Hits[len(p.Hits)-1]
	if len(p.Hits) == depth {
		from := firstNode.SoldCurrency
		to := lastNode.BoughtCurrency
		if from == to {
			paths = []Path{p}
		}
	} else if len(p.Hits) < depth {
		from := lastNode.BoughtCurrency
		for _, to := range m.Currencies {
			if to != from {
				for _, exch := range m.Exchanges {
					n, endpointLookup, nodeLookup := nodeFromMashup(from, to, exch, m, endpointLookup, nodeLookup)
					if n != nil {
						firstFrom := firstNode.SoldCurrency
						nextTo := n.BoughtCurrency
						if (nextTo == firstFrom) && p.contains(*n) == false {
							pathToEvaluate := Path{append(p.Hits, n), nil, nil, 0}
							candidates := findPaths(m, depth, pathToEvaluate, endpointLookup, nodeLookup)
							if len(candidates) > 0 {
								paths = append(paths, candidates...)
							}
						} else if len(p.Hits) < depth-1 {
							if p.contains(*n) == false {
								pathToEvaluate := Path{append(p.Hits, n), nil, nil, 0}
								candidates := findPaths(m, depth, pathToEvaluate, endpointLookup, nodeLookup)
								if len(candidates) > 0 {
									paths = append(paths, candidates...)
								}
							}
						}
					}
				}
			}
		}
	}
	return paths
}

func nodeFromMashup(from core.Currency, to core.Currency, exchange core.Exchange, mashup core.ExchangeMashup, endpointLookup map[string]*core.Endpoint, nodeLookup map[string]*core.Hit) (*core.Hit, map[string]*core.Endpoint, map[string]*core.Hit) {
	var n *core.Hit
	ok := mashup.LinkExist(from, to, exchange)
	if ok {
		var base, quote core.Currency
		base = from
		quote = to
		proto := core.Endpoint{base, quote, exchange, nil}
		endpoint, ok := endpointLookup[proto.ID()]
		if !ok {
			endpointLookup[proto.ID()] = &proto
			endpoint = &proto
		}
		cproto := core.Hit{endpoint, true, from, to}
		node, ok := nodeLookup[cproto.ID()]
		if !ok {
			nodeLookup[cproto.ID()] = &cproto
			node = &cproto
		}
		n = node
	} else {
		ok := mashup.LinkExist(to, from, exchange)
		if ok {
			var base, quote core.Currency
			base = to
			quote = from
			proto := core.Endpoint{base, quote, exchange, nil}
			endpoint, ok := endpointLookup[proto.ID()]
			if !ok {
				endpointLookup[proto.ID()] = &proto
				endpoint = &proto
			}
			cproto := core.Hit{endpoint, false, from, to}
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
