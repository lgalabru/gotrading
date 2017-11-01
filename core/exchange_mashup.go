package core

type ExchangeMashup struct {
	Currencies map[Currency]int
	Exchanges  map[Exchange]int
	Orderbooks [][][]*Orderbook
}

func (mashup *ExchangeMashup) Init(currencies []Currency, exchanges []Exchange) {
	mashup.Currencies = make(map[Currency]int, len(currencies))
	for i, curr := range currencies {
		mashup.Currencies[curr] = i
	}

	mashup.Exchanges = make(map[Exchange]int, len(exchanges))
	for i, exch := range exchanges {
		mashup.Exchanges[exch] = i
	}

	mashup.Orderbooks = make([][][]*Orderbook, len(currencies))
	for i, _ := range currencies {
		mashup.Orderbooks[i] = make([][]*Orderbook, len(currencies))
		for j, _ := range currencies {
			mashup.Orderbooks[i][j] = make([]*Orderbook, len(exchanges))
		}
	}
}

func (mashup *ExchangeMashup) GetOrderbook(from Currency, to Currency, exchange Exchange) *Orderbook {
	return mashup.Orderbooks[mashup.Currencies[from]][mashup.Currencies[to]][mashup.Exchanges[exchange]]
}
