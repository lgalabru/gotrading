package arbitrage

import (
	"fmt"
	"math"
	"time"

	"gotrading/core"
	"gotrading/graph"
	"gotrading/networking"
)

type Simulation struct {
	StartedAt time.Time `json:"startedAt"`
	EndedAt   time.Time `json:"endedAt"`
	hits      []*core.Hit
	Report    Report
}

func (sim *Simulation) Init(hits []*core.Hit) {
	sim.hits = hits
	sim.Report = Report{}
}

func (sim *Simulation) Run() {
	r := &sim.Report
	r.SimulationStartedAt = time.Now()
	r.IsSimulationIncomplete = false
	batch := networking.Batch{}
	batch.UpdateOrderbooks(sim.hits, func(path graph.Path) {

		fromInitialToCurrent := float64(1)

		// rateForInitialCurrency := float64(1) // How many INITIAL_CURRENCY are we getting for 1 CURRENT_CURRENCY

		r.Cost = 0
		r.Rates = make([]float64, len(path.Hits))
		r.AdjustedVolumes = make([]float64, len(path.Hits))
		r.IsSimulationSuccessful = true
		r.Orders = make([]core.Order, len(path.Hits))

		for i, n := range path.Hits {
			if n.Endpoint.Orderbook == nil {
				r.IsSimulationSuccessful = false
				r.SimulationEndedAt = time.Now()
				r.IsSimulationIncomplete = true
				return
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
						r.IsSimulationSuccessful = false
						r.SimulationEndedAt = time.Now()
						r.IsSimulationIncomplete = true
						return
					}
				} else {
					r.IsSimulationSuccessful = false
					r.SimulationEndedAt = time.Now()
					r.IsSimulationIncomplete = true
					return
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
						r.IsSimulationSuccessful = false
						r.SimulationEndedAt = time.Now()
						r.IsSimulationIncomplete = true
						return
					}
				} else {
					r.IsSimulationSuccessful = false
					r.SimulationEndedAt = time.Now()
					r.IsSimulationIncomplete = true
					return
				}
			}

			fromInitialToCurrent = fromInitialToCurrent * priceOfCurrencyToSell
			r.Rates[i] = fromInitialToCurrent
			r.Performance = fromInitialToCurrent

			if i == 0 {
				r.VolumeToEngage = volumeOfCurrencyToSell
			} else {
				limitingAmount := r.VolumeToEngage * fromInitialToCurrent
				currentAmount := volumeOfCurrencyToSell * priceOfCurrencyToSell
				newLimitingAmount := math.Min(limitingAmount, currentAmount)
				r.VolumeToEngage = newLimitingAmount / fromInitialToCurrent
			}
			order.Hit = n
			r.Orders[i] = order
		}

		for i, n := range path.Hits {
			var currentVolumeToEngage float64
			if i == 0 {
				currentVolumeToEngage = r.VolumeToEngage
			} else {
				if path.Hits[i-1].IsBaseToQuote {
					currentVolumeToEngage = r.Orders[i-1].QuoteVolumeOut
				} else if r.Orders[i-1].TransactionType == core.Bid {
					currentVolumeToEngage = r.Orders[i-1].BaseVolumeOut
				}
			}
			if n.IsBaseToQuote {
				r.AdjustedVolumes[i] = currentVolumeToEngage * r.Orders[i].Price
			} else {
				r.AdjustedVolumes[i] = currentVolumeToEngage * r.Orders[i].PriceOfQuoteToBase
			}
			if n.IsBaseToQuote {
				r.Orders[i].UpdateQuoteVolume(r.AdjustedVolumes[i])
			} else {
				r.Orders[i].UpdateBaseVolume(r.AdjustedVolumes[i])
			}
			r.Cost = r.Cost + r.Orders[i].Fee*r.Rates[i]
		}
		r.Path = path

		firstOrder := r.Orders[0]
		if firstOrder.TransactionType == core.Bid {
			r.VolumeIn = firstOrder.QuoteVolumeIn
		} else {
			r.VolumeIn = firstOrder.BaseVolumeIn
		}

		lastOrder := r.Orders[len(r.Orders)-1]
		if lastOrder.TransactionType == core.Bid {
			r.VolumeOut = lastOrder.BaseVolumeOut
		} else {
			r.VolumeOut = lastOrder.QuoteVolumeOut
		}

		if r.VolumeIn < 0.0001 || r.VolumeOut < 0.0001 {
			fmt.Println("Traded volume under threshold")
			r.IsTradedVolumeEnough = false
			r.IsSimulationSuccessful = false
			r.SimulationEndedAt = time.Now()
			return
		}
		r.IsTradedVolumeEnough = true
		r.Performance = r.VolumeOut / r.VolumeIn
		r.SimulationEndedAt = time.Now()
		r.IsSimulationSuccessful = r.Performance > 1.0
	})
}

func (sim *Simulation) IsSuccessful() bool {
	return sim.Report.IsSimulationSuccessful
}

func (sim *Simulation) IsExecutable() bool {
	return sim.Report.IsSimulationIncomplete == false && sim.Report.IsTradedVolumeEnough
}
