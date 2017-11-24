package graph

import (
	"fmt"
	"gotrading/core"
	"strings"
)

type Endpoint struct {
	From      core.Currency   `json:"from"`
	To        core.Currency   `json:"to"`
	Exchange  core.Exchange   `json:"exchange"`
	Orderbook *core.Orderbook `json:"-"`
}

type EndpointLookup struct {
	Endpoint   *Endpoint
	PathsCount int
}

func (e Endpoint) display() {
	fmt.Println(e.Description())
}

func (e Endpoint) isEqual(m Endpoint) bool {
	f := (strings.Compare(string(e.From), string(m.From)) == 0)
	t := (strings.Compare(string(e.To), string(m.To)) == 0)
	fi := (strings.Compare(string(e.To), string(m.From)) == 0)
	ti := (strings.Compare(string(e.From), string(m.To)) == 0)
	exch := (strings.Compare(e.Exchange.Name, m.Exchange.Name) == 0)
	return f && t && exch || fi && ti && exch
}

func (e Endpoint) Description() string {
	var str string
	str = string(e.From) + " / " + string(e.To) + " (" + e.Exchange.Name + ")"
	return str
}

func (e Endpoint) ID() string {
	var str string
	str = string(e.From) + "+" + string(e.To) + "@" + e.Exchange.Name
	return str
}
