package core

type ExchangeMashup struct {
	Currencies []Currency
	Exchanges  []Exchange
	Orderbooks [][][]*Orderbook
}

func (mashup *ExchangeMashup) Init(currencies []Currency, exchanges []Exchange) {
	mashup.Currencies = currencies
	mashup.Exchanges = exchanges

	mashup.Orderbooks = make([][][]*Orderbook, len(currencies))
	for i, _ := range currencies {
		mashup.Orderbooks[i] = make([][]*Orderbook, len(currencies))
		for j, _ := range currencies {
			mashup.Orderbooks[i][j] = make([]*Orderbook, len(exchanges))
		}
	}
}
