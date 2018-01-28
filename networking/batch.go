package networking

import (
	"fmt"
	"gotrading/core"
)

type Batch struct {
}

type orderbooksFetched func(orderbooks []*core.Orderbook)

type sortedOrderbook struct {
	Index     int
	Orderbook *core.Orderbook
}

type sortedOrder struct {
	Index int
	Order *core.Order
}

func (b *Batch) GetOrderbooks(hits []*core.Hit, fn orderbooksFetched) {
	g := SharedGatling()

	orderbooks := make([]*core.Orderbook, len(hits))
	c := make(chan sortedOrderbook, len(hits))

	for i, h := range hits {
		if len(g.Clients) > 1 {
			go b.GetOrderbook(h, i, c)
		} else {
			b.GetOrderbook(h, i, c)
		}
	}

	for range hits {
		elem := <-c
		orderbooks[elem.Index] = elem.Orderbook
	}
	close(c)
	fn(orderbooks)
}

func (b *Batch) GetOrderbook(hit *core.Hit, i int, c chan sortedOrderbook) {
	exchange := hit.Endpoint.Exchange

	o, _ := exchange.GetOrderbook(*hit)
	c <- sortedOrderbook{i, &o}
}

func (b *Batch) PostOrders(orders []core.Order) {
	// g := SharedGatling()

	c := make(chan sortedOrder, len(orders))

	for i, o := range orders {
		// if len(g.Clients) > 1 {
		// 	go b.PostOrder(o, i, c)
		// } else {
		b.PostOrder(o, i, c)
		// }
	}
	// for range orders {
	<-c
	// path.Hits[sortedOrder.Index] = sortedOrder.Order
	// }
	// fn(path)
}

func (b *Batch) PostOrder(order core.Order, i int, c chan sortedOrder) {
	exchange := order.Hit.Endpoint.Exchange

	o, err := exchange.PostOrder(order)
	fmt.Println(o)
	fmt.Println(err)
	c <- sortedOrder{i, &o}
}
