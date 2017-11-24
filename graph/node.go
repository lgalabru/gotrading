package graph

import (
	"gotrading/core"
)

type Node struct {
	Endpoint       *Endpoint      `json:"endpoint"`
	IsBaseToQuote  bool           `json:"isBaseToQuote"`
	SoldCurrency   *core.Currency `json:"soldCurrency"`
	BoughtCurrency *core.Currency `json:"boughtCurrency"`
}

func (n Node) isEqual(m Node) bool {
	return n.Endpoint.isEqual(*m.Endpoint)
}

func (n Node) Description() string {
	var str string
	if n.IsBaseToQuote {
		str = "[" + string(n.Endpoint.From) + "-/+" + string(n.Endpoint.To) + "]@" + n.Endpoint.Exchange.Name
	} else {
		str = "[" + string(n.Endpoint.From) + "+/-" + string(n.Endpoint.To) + "]@" + n.Endpoint.Exchange.Name
	}
	return str
}

func (n Node) ID() string {
	var str string
	if n.IsBaseToQuote {
		str = string(n.Endpoint.From) + "-" + string(n.Endpoint.To) + "@" + n.Endpoint.Exchange.Name
	} else {
		str = string(n.Endpoint.To) + "-" + string(n.Endpoint.From) + "@" + n.Endpoint.Exchange.Name
	}
	return str
}
