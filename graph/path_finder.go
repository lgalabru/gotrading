package graph

import (
	"fmt"
	"strings"

	"gotrading/core"
)

type node struct {
	from     core.Currency
	to       core.Currency
	exch     core.Exchange
	inverted bool
}

type path struct {
	nodes []node
}

func PathFinder(mashup core.ExchangeMashup, from core.Currency, to core.Currency, depth int) []path {
	var paths []path
	for _, to := range mashup.Currencies {
		if to != from {
			for _, exch := range mashup.Exchanges {
				n := nodeFromMashup(from, to, exch, mashup)
				if n != nil {
					paths = append(findPaths(mashup, depth, path{[]node{*n}}), paths...)
				}
			}
		}
	}
	return paths
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
					n := nodeFromMashup(from, to, exch, m)
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

func nodeFromMashup(from core.Currency, to core.Currency, exchange core.Exchange, mashup core.ExchangeMashup) *node {
	var n *node = nil
	ok := mashup.LinkExist(from, to, exchange)
	if ok {
		n = &node{from, to, exchange, false}
	} else {
		ok := mashup.LinkExist(to, from, exchange)
		if ok {
			n = &node{to, from, exchange, true}
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

func (p path) Display() {
	fmt.Println(p.description())
}
