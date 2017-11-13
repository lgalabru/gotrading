package graph

import (
	"gotrading/core"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Finding path between pairs, within and accross exchanges", func() {

	var ()

	BeforeEach(func() {
	})

	Describe(`
Considering the 3 pairs: [ABC/DEF DEF/XYZ XYZ/ABC] available on Exhange1`, func() {
		Context(`
When looking for the available paths within Exhange1, starting and ending with ABC
`, func() {
			var (
				exchanges []core.Exchange
				mashup    core.ExchangeMashup
				from      core.Currency
				to        core.Currency
				nodes     []*Node
				lookups   map[string][]Path
				paths     []Path
			)

			BeforeEach(func() {
				abc := core.Currency("ABC")
				def := core.Currency("DEF")
				xyz := core.Currency("XYZ")

				p1 := core.CurrencyPair{abc, def}
				p2 := core.CurrencyPair{def, xyz}
				p3 := core.CurrencyPair{xyz, abc}

				exchange1 := core.Exchange{"Exchange1", []core.CurrencyPair{p1, p2, p3}, nil}

				exchanges = []core.Exchange{exchange1}
				mashup = core.ExchangeMashup{}
				mashup.Init(exchanges)
				from = abc
				to = from
				depth := 3
				nodes, lookups, paths = PathFinder(mashup, from, to, depth)
			})

			It("should return 2 paths", func() {
				// [ABC-/+DEF] -> [DEF-/+XYZ] -> [XYZ-/+ABC]
				// [XYZ+/-ABC] -> [DEF+/-XYZ] -> [ABC+/-DEF]
				Expect(len(paths)).To(Equal(2))
			})
		})
	})

	Describe(`
Considering the 3 pairs: [BTC/EUR XBT/EUR BTC/XBT] available on Exhange1`, func() {
		Context(`
When looking for the available paths within Exhange1, starting and ending with BTC
`, func() {
			var (
				exchanges []core.Exchange
				mashup    core.ExchangeMashup
				from      core.Currency
				to        core.Currency
				nodes     []*Node
				lookups   map[string][]Path
				paths     []Path
			)

			BeforeEach(func() {
				btc := core.Currency("BTC")
				eur := core.Currency("EUR")
				xbt := core.Currency("XBT")

				p1 := core.CurrencyPair{btc, eur}
				p2 := core.CurrencyPair{xbt, eur}
				p3 := core.CurrencyPair{btc, xbt}

				exchange1 := core.Exchange{"Exchange1", []core.CurrencyPair{p1, p2, p3}, nil}

				exchanges = []core.Exchange{exchange1}
				mashup = core.ExchangeMashup{}
				mashup.Init(exchanges)
				from = btc
				to = from
				depth := 3
				nodes, lookups, paths = PathFinder(mashup, from, to, depth)
			})

			It("should return 2 paths", func() {
				// [BTC-/+XBT] -> [XBT-/+EUR] -> [BTC+/-EUR]
				// [BTC-/+EUR] -> [XBT+/-EUR] -> [BTC+/-XBT]
				Expect(len(paths)).To(Equal(2))
			})
		})
	})
})
