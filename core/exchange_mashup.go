package core

import "fmt"

type ExchangeMashup struct {
	Currencies       []Currency
	Exchanges        []Exchange
	CurrenciesLookup map[Currency]int
	ExchangesLookup  map[Exchange]int
	Orderbooks       [][][]*Orderbook
}

func (mashup *ExchangeMashup) Init(currencies []Currency, exchanges []Exchange) {
	mashup.Currencies = currencies
	mashup.CurrenciesLookup = make(map[Currency]int, len(currencies))
	for i, curr := range currencies {
		mashup.CurrenciesLookup[curr] = i
	}

	mashup.Exchanges = exchanges
	mashup.ExchangesLookup = make(map[Exchange]int, len(exchanges))
	for i, exch := range exchanges {
		mashup.ExchangesLookup[exch] = i
	}

	mashup.Orderbooks = make([][][]*Orderbook, len(currencies))
	for i, _ := range currencies {
		mashup.Orderbooks[i] = make([][]*Orderbook, len(currencies))
		for j, _ := range currencies {
			mashup.Orderbooks[i][j] = make([]*Orderbook, len(exchanges))
		}
	}
}

func (mashup *ExchangeMashup) AddOrderbook(o Orderbook, exchange Exchange) {
	fmt.Println("Inserting item at", mashup.CurrenciesLookup[o.CurrencyPair.From], mashup.CurrenciesLookup[o.CurrencyPair.To])
	mashup.Orderbooks[mashup.CurrenciesLookup[o.CurrencyPair.From]][mashup.CurrenciesLookup[o.CurrencyPair.To]][mashup.ExchangesLookup[exchange]] = &o
}

func (mashup *ExchangeMashup) GetOrderbook(from Currency, to Currency, exchange Exchange) *Orderbook {
	fmt.Println("Looking for item at", mashup.CurrenciesLookup[from], mashup.CurrenciesLookup[to], exchange.Name)
	return mashup.Orderbooks[mashup.CurrenciesLookup[from]][mashup.CurrenciesLookup[to]][mashup.ExchangesLookup[exchange]]
}
