package core

import "fmt"

type Portfolio struct {
	Currencies map[Currency]int
	Exchanges  map[Exchange]int
	Balances   [][]float64
}

func (p *Portfolio) Init(currencies []Currency, exchanges []Exchange) {
	p.Currencies = make(map[Currency]int, len(currencies))
	for i, curr := range currencies {
		p.Currencies[curr] = i
	}

	p.Exchanges = make(map[Exchange]int, len(exchanges))
	for i, exch := range exchanges {
		p.Exchanges[exch] = i
	}

	p.Balances = make([][]float64, len(currencies))
	for i, _ := range currencies {
		p.Balances[i] = make([]float64, len(exchanges))
		for j, _ := range exchanges {
			p.Balances[i][j] = 0
		}
	}
}

func (p *Portfolio) Fullfill(order Order, volume float64, pair CurrencyPair, exchange Exchange) {
	if order.Type == Sell {
		p.Bid(order, volume, pair, exchange)
	} else if order.Type == Buy {
		p.Ask(order, volume, pair, exchange)
	}
}

func (p *Portfolio) Ask(order Order, volume float64, pair CurrencyPair, exchange Exchange) {
	p.DidSold(order.Price, volume, pair, exchange)
}

func (p *Portfolio) Bid(order Order, volume float64, pair CurrencyPair, exchange Exchange) {
	p.DidBuy(order.Price, volume, pair, exchange)
}

func (p *Portfolio) DidSold(price float64, volume float64, pair CurrencyPair, exchange Exchange) {
	p.Balances[p.Currencies[pair.From]][p.Exchanges[exchange]] -= volume
	p.Balances[p.Currencies[pair.To]][p.Exchanges[exchange]] += (volume * price)
}

func (p *Portfolio) DidBuy(price float64, volume float64, pair CurrencyPair, exchange Exchange) {
	p.Balances[p.Currencies[pair.From]][p.Exchanges[exchange]] += volume
	p.Balances[p.Currencies[pair.To]][p.Exchanges[exchange]] -= (volume * price)
}

func (p *Portfolio) DisplayBalances() {
	for curr, i := range p.Currencies {
		for exch, j := range p.Exchanges {
			fmt.Println(exch, curr, p.Balances[i][j])
		}
	}
}
