package graph

import (
	"fmt"
	"gotrading/core"
	"strings"
)

type Node struct {
	From      core.Currency   `json:"from"`
	To        core.Currency   `json:"to"`
	Exchange  core.Exchange   `json:"exchange"`
	Orderbook *core.Orderbook `json:"-"`
}

type ContextualNode struct {
	Node           *Node          `json:"content"`
	Inverted       bool           `json:"inverted"`
	SoldCurrency   *core.Currency `json:"sold_currency"`
	BoughtCurrency *core.Currency `json:"bought_currency"`
}

type NodeLookup struct {
	Node       *Node
	PathsCount int
}

// func (n NodeLookup) String() string {
// 	return n.Node.Description() + " " + strconv.Itoa(n.PathsCount)
// }

func (n Node) display() {
	fmt.Println(n.Description())
}

func (n Node) isEqual(m Node) bool {
	f := (strings.Compare(string(n.From), string(m.From)) == 0)
	t := (strings.Compare(string(n.To), string(m.To)) == 0)
	fi := (strings.Compare(string(n.To), string(m.From)) == 0)
	ti := (strings.Compare(string(n.From), string(m.To)) == 0)
	e := (strings.Compare(n.Exchange.Name, m.Exchange.Name) == 0)
	return f && t && e || fi && ti && e
}

func (n ContextualNode) isEqual(m ContextualNode) bool {
	return n.Node.isEqual(*m.Node)
}

func (n ContextualNode) Description() string {
	var str string
	if n.Inverted {
		// str = "Use " + string(n.Node.From) + " to buy " + string(n.Node.To) + " on " + n.Node.Exchange.Name + "."
		str = "[" + string(n.Node.From) + "+/-" + string(n.Node.To) + "]@" + n.Node.Exchange.Name
	} else {
		// str = "Sell " + string(n.Node.From) + " for " + string(n.Node.To) + " on " + n.Node.Exchange.Name + "."
		str = "[" + string(n.Node.From) + "-/+" + string(n.Node.To) + "]@" + n.Node.Exchange.Name
	}
	return str
}

func (n Node) Description() string {
	var str string
	str = string(n.From) + " / " + string(n.To) + " (" + n.Exchange.Name + ")"
	return str
}

func (n Node) ID() string {
	var str string
	str = string(n.From) + "+" + string(n.To) + "@" + n.Exchange.Name
	return str
}

func (n ContextualNode) ID() string {
	var str string
	if n.Inverted {
		str = string(n.Node.To) + "-" + string(n.Node.From) + "@" + n.Node.Exchange.Name
	} else {
		str = string(n.Node.From) + "-" + string(n.Node.To) + "@" + n.Node.Exchange.Name
	}
	return str
}
