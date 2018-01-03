package networking

import (
	"fmt"
	"gotrading/core"
	"gotrading/graph"
)

type Batch struct {
}

type pathFetched func(path graph.Path)

type indexedHit struct {
	Index int
	Hit   *core.Hit
}

func (b *Batch) UpdateOrderbooks(hits []*core.Hit, fn pathFetched) {
	g := SharedGatling()

	path := graph.Path{}
	path.Hits = hits
	c := make(chan indexedHit, len(hits))

	for i, n := range hits {
		if len(g.Clients) > 1 {
			go b.GetOrderbook(n, i, c)
		} else {
			b.GetOrderbook(n, i, c)
		}
	}
	for range hits {
		indexedHit := <-c
		path.Hits[indexedHit.Index] = indexedHit.Hit
	}
	path.Encode()
	fn(path)
}

func (b *Batch) GetOrderbook(hit *core.Hit, i int, c chan indexedHit) {
	exchange := hit.Endpoint.Exchange

	o, _ := exchange.GetOrderbook(*hit)
	hit.Endpoint.Orderbook = &o
	c <- indexedHit{i, hit}
}
