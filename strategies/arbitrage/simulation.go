package arbitrage

import (
	"math"
	"time"

	"gotrading/core"
	"gotrading/graph"
	"gotrading/networking"
)

type Simulation struct {
	StartedAt    time.Time `json:"startedAt"`
	EndedAt      time.Time `json:"endedAt"`
	hits         []*core.Hit
	chain        ChainedOrders
	IsSuccessful bool
}

func (sim *Simulation) Init(hits []*core.Hit) {
	sim.hits = hits
}

func (sim *Simulation) Run() {

	batch := networking.Batch{}
	batch.UpdateOrderbooks(sim.hits, func(path graph.Path) {

		sim.chain = ChainedOrders{}
		fromInitialToCurrent := float64(1)

		// rateForInitialCurrency := float64(1) // How many INITIAL_CURRENCY are we getting for 1 CURRENT_CURRENCY

		sim.chain.CreatedAt = time.Now()
		sim.chain.Cost = 0
		sim.chain.Rates = make([]float64, len(path.Hits))
		sim.chain.AdjustedVolumes = make([]float64, len(path.Hits))
		sim.chain.IsBroken = false
		sim.chain.Orders = make([]core.Order, len(path.Hits))

		for i, n := range path.Hits {
			if n.Endpoint.Orderbook == nil {
				sim.chain.IsBroken = true
				continue
			}

			var priceOfCurrencyToSell float64
			var volumeOfCurrencyToSell float64
			var order core.Order

			if n.IsBaseToQuote {
				// We are selling the base -> we match the Bid.
				if len(n.Endpoint.Orderbook.Bids) > 0 {
					bestBid := n.Endpoint.Orderbook.Bids[0]
					o, err := bestBid.CreateMatchingAsk()
					if err == nil {
						order = *o
						priceOfCurrencyToSell = order.Price
						volumeOfCurrencyToSell = order.BaseVolume
					} else {
						sim.chain.IsBroken = true
						continue
					}
				} else {
					sim.chain.IsBroken = true
					continue
				}
			} else {
				// We are selling the quote <=> we are buying the base -> we match the Ask.
				if len(n.Endpoint.Orderbook.Asks) > 0 {
					bestAsk := n.Endpoint.Orderbook.Asks[0]
					o, err := bestAsk.CreateMatchingBid()
					if err == nil {
						order = *o
						priceOfCurrencyToSell = order.PriceOfQuoteToBase
						volumeOfCurrencyToSell = order.QuoteVolume
					} else {
						sim.chain.IsBroken = true
						continue
					}
				} else {
					sim.chain.IsBroken = true
					continue
				}
			}

			fromInitialToCurrent = fromInitialToCurrent * priceOfCurrencyToSell
			sim.chain.Rates[i] = fromInitialToCurrent
			sim.chain.Performance = fromInitialToCurrent

			if i == 0 {
				sim.chain.VolumeToEngage = volumeOfCurrencyToSell
			} else {
				limitingAmount := sim.chain.VolumeToEngage * fromInitialToCurrent
				currentAmount := volumeOfCurrencyToSell * priceOfCurrencyToSell
				newLimitingAmount := math.Min(limitingAmount, currentAmount)
				sim.chain.VolumeToEngage = newLimitingAmount / fromInitialToCurrent
			}
			sim.chain.Orders[i] = order
		}

		for i, n := range path.Hits {
			var currentVolumeToEngage float64
			if i == 0 {
				currentVolumeToEngage = sim.chain.VolumeToEngage
			} else {
				if path.Hits[i-1].IsBaseToQuote {
					currentVolumeToEngage = sim.chain.Orders[i-1].QuoteVolumeOut
				} else if sim.chain.Orders[i-1].TransactionType == core.Bid {
					currentVolumeToEngage = sim.chain.Orders[i-1].BaseVolumeOut
				}
			}
			if n.IsBaseToQuote {
				sim.chain.AdjustedVolumes[i] = currentVolumeToEngage * sim.chain.Orders[i].Price
			} else {
				sim.chain.AdjustedVolumes[i] = currentVolumeToEngage * sim.chain.Orders[i].PriceOfQuoteToBase
			}
			if n.IsBaseToQuote {
				sim.chain.Orders[i].UpdateQuoteVolume(sim.chain.AdjustedVolumes[i])
			} else {
				sim.chain.Orders[i].UpdateBaseVolume(sim.chain.AdjustedVolumes[i])
			}
			sim.chain.Cost = sim.chain.Cost + sim.chain.Orders[i].Fee*sim.chain.Rates[i]
		}
		sim.chain.Path = path

		firstOrder := sim.chain.Orders[0]
		if firstOrder.TransactionType == core.Bid {
			sim.chain.VolumeIn = firstOrder.QuoteVolumeIn
		} else {
			sim.chain.VolumeIn = firstOrder.BaseVolumeIn
		}

		lastOrder := sim.chain.Orders[len(sim.chain.Orders)-1]
		if lastOrder.TransactionType == core.Bid {
			sim.chain.VolumeOut = lastOrder.BaseVolumeOut
		} else {
			sim.chain.VolumeOut = lastOrder.QuoteVolumeOut
		}
		if sim.chain.VolumeIn < 0.0001 || sim.chain.VolumeOut < 0.0001 {
			sim.chain.IsBroken = true
		}
		sim.chain.Performance = sim.chain.VolumeOut / sim.chain.VolumeIn
		sim.chain.DiagnosedAt = time.Now()

	})
}

func (sim *Simulation) BuildReport() Report {
	return Report{}
}
