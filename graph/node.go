package graph

import (
	"fmt"
	"gotrading/core"
	"strconv"
	"strings"
)

type Node struct {
	From     core.Currency
	To       core.Currency
	Exchange core.Exchange
	Inverted bool
}

type NodeLookup struct {
	Node       Node
	PathsCount int
}

func (n NodeLookup) String() string {
	return n.Node.Description() + " " + strconv.Itoa(n.PathsCount)
}

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

func (n Node) Description() string {
	var str string
	if n.Inverted {
		str = string(n.From) + " <- " + string(n.To) + " (" + n.Exchange.Name + ")"
	} else {
		str = string(n.From) + " -> " + string(n.To) + " (" + n.Exchange.Name + ")"
	}
	return str
}

func (n Node) ID() string {
	return string(n.From) + "-" + string(n.To) + "@" + n.Exchange.Name
}
