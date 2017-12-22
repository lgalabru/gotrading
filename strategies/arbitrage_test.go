package strategies

import (
	"fmt"
	"gotrading/core"
	"gotrading/graph"
	"time"

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
Considering the combination: [ABC/DEF]@Exhange1 -> [XYZ/DEF]@Exhange1 and the orderbooks:
[ABC/DEF]@Exhange1 -> Best Bid: 1ABC = 2DEF, Best Ask: 1ABC = 2DEF
[XYZ/DEF]@Exhange1 -> Best Bid: 1XYZ = 3DEF, Best Ask: 1XYZ = 3DEF`, func() {
		Context(`
When I fulfill all the orders, running the arbitrage`, func() {
			var (
				chains []ArbitrageChain
				ob1    core.Orderbook
				ob2    core.Orderbook
				// ob3    core.Orderbook
				paths []graph.Path
			)

			BeforeEach(func() {
				exchange1 := core.Exchange{"Exchange1", make([]core.CurrencyPair, 0), nil}

				abc := core.Currency("ABC")
				def := core.Currency("DEF")
				xyz := core.Currency("XYZ")

				pair1 := core.CurrencyPair{abc, def}
				pair2 := core.CurrencyPair{xyz, def}

				bids1 := append([]core.Order{}, core.NewBid(pair1, 2, 6))
				asks1 := append([]core.Order{}, core.NewAsk(pair1, 2, 6))
				bids2 := append([]core.Order{}, core.NewBid(pair2, 3, 2))
				asks2 := append([]core.Order{}, core.NewAsk(pair2, 3, 2))

				ob1 = core.Orderbook{pair1, bids1, asks1, time.Now()}
				ob2 = core.Orderbook{pair2, bids2, asks2, time.Now()}

				endpoint1 := graph.Endpoint{abc, def, exchange1, &ob1}
				endpoint2 := graph.Endpoint{xyz, def, exchange1, &ob2}

				nodes := make([]*graph.Node, 2)
				nodes[0] = &(graph.Node{&endpoint1, true, &abc, &def})
				nodes[1] = &(graph.Node{&endpoint2, false, &xyz, &def})

				paths = make([]graph.Path, 1)
				paths[0] = graph.Path{nodes, nil, nil, 0}

				chains = arbitrage.Run(paths)
			})

			It("should return one chain", func() {
				Expect(len(chains)).To(Equal(1))
			})

			It("should return one chain enforcing the initial volume to 3", func() {
				c := chains[0]
				fmt.Println(c)
				Expect(c.VolumeToEngage).To(Equal(3.0))
			})

			It("should return one chain announcing a performance equal to 1x", func() {
				c := chains[0]
				Expect(c.Performance).To(Equal(2. / 3.))
			})

			It("should return 2 as a rate for the node #1", func() {
				c := chains[0]
				Expect(c.Rates[0]).To(Equal(2.0))
			})

			It("should return 2./3. as an adjusted volume for the node #2", func() {
				c := chains[0]
				Expect(c.Rates[1]).To(Equal(2. / 3.))
			})

			It("should return 6 as an adjusted volume for the node #1", func() {
				c := chains[0]
				Expect(c.AdjustedVolumes[0]).To(Equal(6.))
			})

			It("should return 2 as an adjusted volume for the node #2", func() {
				c := chains[0]
				Expect(c.AdjustedVolumes[1]).To(Equal(2.))
			})

			It("should return 3 for the volume of the order corresponding to node #1", func() {
				c := chains[0]
				Expect(c.Orders[0].BaseVolume).To(Equal(3.))
			})

			It("should return 2 for the volume of the order corresponding to node #2", func() {
				c := chains[0]
				Expect(c.Orders[1].BaseVolume).To(Equal(2.))
			})

		})
	})

	// 	Describe(`
	// Considering the combination: [ABC/DEF]@Exhange1 -> [DEF/XYZ]@Exhange1 -> [XYZ/ABC]@Exhange1, and the orderbooks:
	// [ABC/DEF]@Exhange1 -> Best Bid: 1ABC = 10DEF, Best Ask: 1ABC = 10DEF #ABC=0, DEF=10, XYZ=0
	// [DEF/XYZ]@Exhange1 -> Best Bid: 1DEF = 10XYZ, Best Ask: 1DEF = 10XYZ #ABC=0, DEF=0, XYZ=100
	// [XYZ/ABC]@Exhange1 -> Best Bid: 1XYZ = 0.01ABC, Best Ask: 1XYZ = 0.01ABC`, func() {
	// 		Context(`
	// When I fulfill all the orders, running the arbitrage`, func() {
	// 			var (
	// 				chains []ArbitrageChain
	// 				ob1    core.Orderbook
	// 				ob2    core.Orderbook
	// 				ob3    core.Orderbook
	// 				paths  []graph.Path
	// 			)

	// 			BeforeEach(func() {
	// 				exchange1 := core.Exchange{"Exchange1", make([]core.CurrencyPair, 0), nil}

	// 				abc := core.Currency("ABC")
	// 				def := core.Currency("DEF")
	// 				xyz := core.Currency("XYZ")

	// 				pair1 := core.CurrencyPair{abc, def}
	// 				pair2 := core.CurrencyPair{def, xyz}
	// 				pair3 := core.CurrencyPair{xyz, abc}

	// 				bids1 := append([]core.Order{}, core.NewBid(pair1, 10, 1))
	// 				asks1 := append([]core.Order{}, core.NewAsk(pair1, 10, 1))
	// 				bids2 := append([]core.Order{}, core.NewBid(pair2, 10, 10))
	// 				asks2 := append([]core.Order{}, core.NewAsk(pair2, 10, 1))
	// 				bids3 := append([]core.Order{}, core.NewBid(pair3, 0.01, 100))
	// 				asks3 := append([]core.Order{}, core.NewAsk(pair3, 0.01, 1))

	// 				ob1 = core.Orderbook{pair1, bids1, asks1, time.Now()}
	// 				ob2 = core.Orderbook{pair2, bids2, asks2, time.Now()}
	// 				ob3 = core.Orderbook{pair3, bids3, asks3, time.Now()}

	// 				endpoint1 := graph.Endpoint{abc, def, exchange1, &ob1}
	// 				endpoint2 := graph.Endpoint{def, xyz, exchange1, &ob2}
	// 				endpoint3 := graph.Endpoint{xyz, abc, exchange1, &ob3}

	// 				nodes := make([]*graph.Node, 3)
	// 				nodes[0] = &(graph.Node{&endpoint1, true, &abc, &def})
	// 				nodes[1] = &(graph.Node{&endpoint2, true, &def, &xyz})
	// 				nodes[2] = &(graph.Node{&endpoint3, true, &xyz, &abc})

	// 				paths = make([]graph.Path, 1)
	// 				paths[0] = graph.Path{nodes, nil, nil}

	// 				chains = arbitrage.Run(paths)
	// 			})

	// 			It("should return one chain", func() {
	// 				Expect(len(chains)).To(Equal(1))
	// 			})

	// 			It("should return one chain enforcing the initial volume to 1", func() {
	// 				c := chains[0]
	// 				fmt.Println(c)
	// 				Expect(c.VolumeToEngage).To(Equal(1.0))
	// 			})

	// 			It("should return one chain announcing a performance equal to 1x", func() {
	// 				c := chains[0]
	// 				Expect(c.Performance).To(Equal(1.0))
	// 			})

	// 			It("should return one chain announcing a performance equal to 10x if 1XYZ = 0.10ABC instead of 1XYZ = 0.01ABC", func() {
	// 				ob3.Bids[0].Price = 0.10
	// 				chains = arbitrage.Run(paths)
	// 				c := chains[0]
	// 				Expect(c.Performance).To(Equal(10.0))
	// 			})

	// 			It("should return one chain announcing a performance equal to 10x if 1XYZ = 0.10ABC instead of 1XYZ = 0.01ABC", func() {
	// 				ob3.Bids[0].Price = 0.10
	// 				chains = arbitrage.Run(paths)
	// 				c := chains[0]
	// 				Expect(c.Performance).To(Equal(10.0))
	// 			})

	// 			It("should return one chain enforcing the initial volume to 0.1 if only 10 XYZ are available", func() {
	// 				ob3.Bids[0].BaseVolume = 10
	// 				chains = arbitrage.Run(paths)
	// 				c := chains[0]
	// 				Expect(c.VolumeToEngage).To(Equal(0.1))
	// 			})
	// 		})
	// 	})

	// 	Describe(`
	// Considering the combination: [ABC/DEF]@Exhange1 -> [XYZ/DEF]@Exhange1 -> [XYZ/ABC]@Exhange1, and the orderbooks:
	// [ABC/DEF]@Exhange1 -> Best Bid: 1ABC = 10DEF, Best Ask: 1ABC = 10DEF #ABC=0, DEF=10, XYZ=0
	// [XYZ/DEF]@Exhange1 -> Best Bid: 1XYZ = 0.01DEF, Best Ask: 1DEF = 0.1XYZ #ABC=0, DEF=0, XYZ=100
	// [XYZ/ABC]@Exhange1 -> Best Bid: 1XYZ = 0.1ABC, Best Ask: 1XYZ = 0.01ABC`, func() {
	// 		Context(`
	// When I fulfill all the orders, running the arbitrage`, func() {
	// 			var (
	// 				chains []ArbitrageChain
	// 				ob1    core.Orderbook
	// 				ob2    core.Orderbook
	// 				ob3    core.Orderbook
	// 				paths  []graph.Path
	// 			)

	// 			BeforeEach(func() {
	// 				exchange1 := core.Exchange{"Exchange1", make([]core.CurrencyPair, 0), nil}

	// 				abc := core.Currency("ABC")
	// 				def := core.Currency("DEF")
	// 				xyz := core.Currency("XYZ")

	// 				pair1 := core.CurrencyPair{abc, def}
	// 				pair2 := core.CurrencyPair{xyz, def}
	// 				pair3 := core.CurrencyPair{xyz, abc}

	// 				bids1 := append([]core.Order{}, core.NewBid(pair1, 10, 1))
	// 				asks1 := append([]core.Order{}, core.NewAsk(pair1, 10, 1))
	// 				bids2 := append([]core.Order{}, core.NewBid(pair2, 0.01, 1000))
	// 				asks2 := append([]core.Order{}, core.NewAsk(pair2, 0.01, 1000))
	// 				bids3 := append([]core.Order{}, core.NewBid(pair3, 0.001, 1000))
	// 				asks3 := append([]core.Order{}, core.NewAsk(pair3, 0.001, 1000))

	// 				ob1 = core.Orderbook{pair1, bids1, asks1, time.Now()}
	// 				ob2 = core.Orderbook{pair2, bids2, asks2, time.Now()}
	// 				ob3 = core.Orderbook{pair3, bids3, asks3, time.Now()}

	// 				endpoint1 := graph.Endpoint{abc, def, exchange1, &ob1}
	// 				endpoint2 := graph.Endpoint{xyz, def, exchange1, &ob2}
	// 				endpoint3 := graph.Endpoint{xyz, abc, exchange1, &ob3}

	// 				nodes := make([]*graph.Node, 3)
	// 				nodes[0] = &(graph.Node{&endpoint1, true, &abc, &def})
	// 				nodes[1] = &(graph.Node{&endpoint2, false, &xyz, &def})
	// 				nodes[2] = &(graph.Node{&endpoint3, true, &xyz, &abc})

	// 				paths = make([]graph.Path, 1)
	// 				paths[0] = graph.Path{nodes, nil, nil}

	// 				chains = arbitrage.Run(paths)
	// 			})

	// 			It("should return one chain", func() {
	// 				Expect(len(chains)).To(Equal(1))
	// 			})

	// 			It("should return one chain enforcing the initial volume to 1", func() {
	// 				c := chains[0]
	// 				fmt.Println(c)
	// 				Expect(c.VolumeToEngage).To(Equal(1.0))
	// 			})

	// 			It("should return one chain announcing a performance equal to 1x", func() {
	// 				c := chains[0]
	// 				Expect(c.Performance).To(Equal(1.0))
	// 			})

	// 			It("should return one chain announcing a performance equal to 10x if 1XYZ = 0.10ABC instead of 1XYZ = 0.01ABC", func() {
	// 				ob3.Bids[0].Price = 0.10
	// 				chains = arbitrage.Run(paths)
	// 				c := chains[0]
	// 				Expect(c.Performance).To(Equal(10.0))
	// 			})

	// 			It("should return one chain announcing a performance equal to 10x if 1XYZ = 0.10ABC instead of 1XYZ = 0.01ABC", func() {
	// 				ob3.Bids[0].Price = 0.10
	// 				chains = arbitrage.Run(paths)
	// 				c := chains[0]
	// 				Expect(c.Performance).To(Equal(10.0))
	// 			})

	// 			It("should return one chain enforcing the initial volume to 0.1 if only 10 XYZ are available", func() {
	// 				ob3.Bids[0].BaseVolume = 10
	// 				chains = arbitrage.Run(paths)
	// 				c := chains[0]
	// 				Expect(c.VolumeToEngage).To(Equal(0.1))
	// 			})
	// 		})
	// 	})
})
