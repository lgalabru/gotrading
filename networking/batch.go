package networking

import (
	"fmt"
	"gotrading/core"
	"gotrading/graph"
)

type Batch struct {
}

type pathFetched func(path graph.Path)

type sortedHit struct {
	Index int
	Hit   *core.Hit
}

type sortedOrder struct {
	Index int
	Order *core.Order
}

func (b *Batch) UpdateOrderbooks(hits []*core.Hit, fn pathFetched) {
	g := SharedGatling()

	path := graph.Path{}
	path.Hits = hits
	c := make(chan sortedHit, len(hits))

	for i, n := range hits {
		if len(g.Clients) > 1 {
			go b.GetOrderbook(n, i, c)
		} else {
			b.GetOrderbook(n, i, c)
		}
	}
	for range hits {
		sortedHit := <-c
		path.Hits[sortedHit.Index] = sortedHit.Hit
	}
	path.Encode()
	fn(path)
}

func (b *Batch) GetOrderbook(hit *core.Hit, i int, c chan sortedHit) {
	exchange := hit.Endpoint.Exchange

	o, _ := exchange.GetOrderbook(*hit)
	hit.Endpoint.Orderbook = &o
	c <- sortedHit{i, hit}
}

func (b *Batch) PostOrders(orders []core.Order) {
	g := SharedGatling()

	c := make(chan sortedOrder, len(orders))

	for i, o := range orders {
		if len(g.Clients) > 1 {
			go b.PostOrder(o, i, c)
		} else {
			b.PostOrder(o, i, c)
		}
	}
	for range orders {
		sortedOrder := <-c
		fmt.Println(sortedOrder)
		// path.Hits[sortedOrder.Index] = sortedOrder.Order
	}
	// fn(path)
}

func (b *Batch) PostOrder(order core.Order, i int, c chan sortedOrder) {
	exchange := order.Hit.Endpoint.Exchange

	o, err := exchange.PostOrder(order)
	fmt.Println(o)
	fmt.Println(err)
	// hit.Endpoint.Orderbook = &o
	// c <- indexedHit{i, hit}
}
