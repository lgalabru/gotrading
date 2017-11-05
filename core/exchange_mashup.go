package core

type ExchangeMashup struct {
	Currencies       []Currency
	Exchanges        []Exchange
	CurrenciesLookup map[Currency]int
	ExchangesLookup  map[*Exchange]int
	Links            [][][]bool
}

func (m *ExchangeMashup) Init(exchanges []Exchange) {
	m.CurrenciesLookup = make(map[Currency]int)
	m.Currencies = make([]Currency, 0)
	m.Exchanges = exchanges
	m.ExchangesLookup = make(map[*Exchange]int, len(exchanges))

	for i, exch := range exchanges {
		for _, pair := range exch.AvailablePairs {

			_, ok := m.CurrenciesLookup[pair.From]
			if !ok {
				m.Currencies = append(m.Currencies, pair.From)
				m.CurrenciesLookup[pair.From] = len(m.Currencies) - 1
			}
			_, ok = m.CurrenciesLookup[pair.To]
			if !ok {
				m.Currencies = append(m.Currencies, pair.To)
				m.CurrenciesLookup[pair.To] = len(m.Currencies) - 1

			}
		}
		m.ExchangesLookup[&exch] = i
	}

	m.Links = make([][][]bool, len(m.Currencies))
	for i, _ := range m.Currencies {
		m.Links[i] = make([][]bool, len(m.Currencies))
		for j, _ := range m.Currencies {
			m.Links[i][j] = make([]bool, len(exchanges))
		}
	}

	for _, exch := range exchanges {
		for _, pair := range exch.AvailablePairs {
			m.Links[m.CurrenciesLookup[pair.From]][m.CurrenciesLookup[pair.To]][m.ExchangesLookup[&exch]] = true
		}
	}
}

func (m *ExchangeMashup) LinkExist(from Currency, to Currency, exchange Exchange) bool {
	ok := m.Links[m.CurrenciesLookup[from]][m.CurrenciesLookup[to]][m.ExchangesLookup[&exchange]]
	return ok
}
