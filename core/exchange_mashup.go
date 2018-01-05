package core

// ExchangeMashup allowing to identify available bridges between currencies accross exchanges
type ExchangeMashup struct {
	Currencies       []Currency
	Exchanges        []Exchange
	CurrenciesLookup map[Currency]int
	ExchangesLookup  map[string]int
	Links            [][][]bool
}

// Init initializes a mashup
func (m *ExchangeMashup) Init(exchanges []Exchange) {
	m.CurrenciesLookup = make(map[Currency]int)
	m.Currencies = make([]Currency, 0)
	m.Exchanges = exchanges
	m.ExchangesLookup = make(map[string]int, len(exchanges))

	for i, exch := range exchanges {
		for _, pair := range exch.PairsEnabled {

			_, ok := m.CurrenciesLookup[pair.Base]
			if !ok {
				m.Currencies = append(m.Currencies, pair.Base)
				m.CurrenciesLookup[pair.Base] = len(m.Currencies) - 1
			}
			_, ok = m.CurrenciesLookup[pair.Quote]
			if !ok {
				m.Currencies = append(m.Currencies, pair.Quote)
				m.CurrenciesLookup[pair.Quote] = len(m.Currencies) - 1

			}
		}
		m.ExchangesLookup[exch.Name] = i
	}

	m.Links = make([][][]bool, len(m.Currencies))
	for i := range m.Currencies {
		m.Links[i] = make([][]bool, len(m.Currencies))
		for j := range m.Currencies {
			m.Links[i][j] = make([]bool, len(exchanges))
			for z := range exchanges {
				m.Links[i][j][z] = false
			}
		}
	}

	for _, exch := range exchanges {
		for _, pair := range exch.PairsEnabled {
			m.Links[m.CurrenciesLookup[pair.Base]][m.CurrenciesLookup[pair.Quote]][m.ExchangesLookup[exch.Name]] = true
		}
	}
}

// LinkExist returns true if a currency pair exists for a given exchange
func (m *ExchangeMashup) LinkExist(base Currency, quote Currency, exch Exchange) bool {
	ok := m.Links[m.CurrenciesLookup[base]][m.CurrenciesLookup[quote]][m.ExchangesLookup[exch.Name]]
	return ok
}
