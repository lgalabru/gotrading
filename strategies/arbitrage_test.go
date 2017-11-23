package strategies

import (
	"gotrading/core"
	"gotrading/graph"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Arbitrage in 3 steps, starting and finishing with ABC", func() {

	var (
		arbitrage Arbitrage
	)

	BeforeEach(func() {
		arbitrage = Arbitrage{}
	})

	Describe(`
Considering the combination: [ABC/DEF]@Exhange1 -> [DEF/XYZ]@Exhange1 -> [XYZ/ABC]@Exhange1, and the orderbooks:
[ABC/DEF]@Exhange1 -> Best Bid: 1ABC = 10DEF, Best Ask: 1ABC = 10DEF #ABC=0, DEF=10, XYZ=0
[DEF/XYZ]@Exhange1 -> Best Bid: 1DEF = 10XYZ, Best Ask: 1DEF = 10XYZ #ABC=0, DEF=0, XYZ=100
[XYZ/ABC]@Exhange1 -> Best Bid: 1XYZ = 0.01ABC, Best Ask: 1XYZ = 0.01ABC`, func() {
		Context(`
When I fulfill all the orders, running the arbitrage`, func() {
			var (
				chains []ArbitrageChain
				ob1    core.Orderbook
				ob2    core.Orderbook
				ob3    core.Orderbook
				paths  []graph.Path
			)

			BeforeEach(func() {
				exchange1 := core.Exchange{"Exchange1", make([]core.CurrencyPair, 0), nil}

				abc := core.Currency("ABC")
				def := core.Currency("DEF")
				xyz := core.Currency("XYZ")

				bids1 := append([]core.Order{}, core.Order{10, 1, core.Buy})
				asks1 := append([]core.Order{}, core.Order{10, 1, core.Sell})
				bids2 := append([]core.Order{}, core.Order{10, 10, core.Buy})
				asks2 := append([]core.Order{}, core.Order{10, 1, core.Sell})
				bids3 := append([]core.Order{}, core.Order{0.01, 100, core.Buy})
				asks3 := append([]core.Order{}, core.Order{0.01, 1, core.Sell})

				ob1 = core.Orderbook{core.CurrencyPair{abc, def}, bids1, asks1}
				ob2 = core.Orderbook{core.CurrencyPair{def, xyz}, bids2, asks2}
				ob3 = core.Orderbook{core.CurrencyPair{xyz, abc}, bids3, asks3}

				node1 := graph.Node{abc, def, exchange1, &ob1}
				node2 := graph.Node{def, xyz, exchange1, &ob2}
				node3 := graph.Node{xyz, abc, exchange1, &ob3}

				cnodes := make([]*graph.ContextualNode, 3)
				cnodes[0] = &(graph.ContextualNode{&node1, false, &abc, &def})
				cnodes[1] = &(graph.ContextualNode{&node2, false, &def, &xyz})
				cnodes[2] = &(graph.ContextualNode{&node3, false, &xyz, &abc})

				paths = make([]graph.Path, 1)
				paths[0] = graph.Path{cnodes, nil, nil}

				chains = arbitrage.Run(paths)
			})

			It("should return one chain", func() {
				Expect(len(chains)).To(Equal(1))
			})

			It("should return one chain enforcing the initial volume to 1", func() {
				c := chains[0]
				Expect(c.Volume).To(Equal(1.0))
			})

			It("should return one chain announcing a performance equal to 1x", func() {
				c := chains[0]
				Expect(c.Performance).To(Equal(1.0))
			})

			It("should return one chain announcing a performance equal to 10x if 1XYZ = 0.10ABC instead of 1XYZ = 0.01ABC", func() {
				ob3.Bids[0].Price = 0.10
				chains = arbitrage.Run(paths)
				c := chains[0]
				Expect(c.Performance).To(Equal(10.0))
			})

			It("should return one chain announcing a performance equal to 10x if 1XYZ = 0.10ABC instead of 1XYZ = 0.01ABC", func() {
				ob3.Bids[0].Price = 0.10
				chains = arbitrage.Run(paths)
				c := chains[0]
				Expect(c.Performance).To(Equal(10.0))
			})

			It("should return one chain enforcing the initial volume to 0.1 if only 10 XYZ are available", func() {
				ob3.Bids[0].Volume = 10
				chains = arbitrage.Run(paths)
				c := chains[0]
				Expect(c.Volume).To(Equal(0.1))
			})
		})
	})
})

// Describe(`
// Considering the combination: [BTC/USD]@Exhange1 -> [ETH/USD]@Exhange1 -> [ETH/BTC]@Exhange1`, func() {
// 		Context(`
// [BTC/USD]@Exhange1 -> Best Bid: 1BTC = 5,999USD / Best Ask: 1BTC = 6,000USD
// [ETH/USD]@Exhange1 -> Best Bid: 1ETH = 299USD / Best Ask: 1ETH = 300USD
// [ETH/BTC]@Exhange1 -> Best Bid: 1ETH = 0.0442BTC / Best Ask: 1ETH = 0.0443BTC
// `, func() {
// 			It("should work", func() {
// 				// func (arbitrage *Arbitrage) Run(paths []graph.Path) []ArbitrageChain {
// 				Expect(1).To(Equal(0))
// 			})
// 		})
// 	})
