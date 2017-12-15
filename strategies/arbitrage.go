package strategies

import (
	"math"

	"gotrading/core"
	"gotrading/graph"
	"gotrading/services"
)

type Arbitrage struct {
	Solutions []Solution
}

func (arbitrage *Arbitrage) Run(paths []graph.Path) []services.ChainedOrders {

	chains := make([]services.ChainedOrders, len(paths))
	// fromCurrentToInitial := float64(1)
	fromInitialToCurrent := float64(1)

	// rateForInitialCurrency := float64(1) // How many INITIAL_CURRENCY are we getting for 1 CURRENT_CURRENCY
	for j, p := range paths {

		chain := services.ChainedOrders{}
		chain.Cost = 0
		chain.Rates = make([]float64, len(p.Nodes))
		chain.AdjustedVolumes = make([]float64, len(p.Nodes))
		chain.IsBroken = false
		chain.Orders = make([]core.Order, len(p.Nodes))

		for i, n := range p.Nodes {
			if n.Endpoint.Orderbook == nil {
				chain.IsBroken = true
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
						chain.IsBroken = true
						continue
					}
				} else {
					chain.IsBroken = true
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
						chain.IsBroken = true
						continue
					}
				} else {
					chain.IsBroken = true
					continue
				}
			}

			fromInitialToCurrent = fromInitialToCurrent * priceOfCurrencyToSell
			chain.Rates[i] = fromInitialToCurrent
			chain.Performance = fromInitialToCurrent

			if i == 0 {
				chain.VolumeToEngage = volumeOfCurrencyToSell
			} else {
				limitingAmount := chain.VolumeToEngage * fromInitialToCurrent
				currentAmount := volumeOfCurrencyToSell * priceOfCurrencyToSell
				newLimitingAmount := math.Min(limitingAmount, currentAmount)
				chain.VolumeToEngage = newLimitingAmount / fromInitialToCurrent
			}
			chain.Orders[i] = order
		}

		for i, n := range p.Nodes {
			var currentVolumeToEngage float64
			if i == 0 {
				currentVolumeToEngage = chain.VolumeToEngage
			} else {
				if p.Nodes[i-1].IsBaseToQuote {
					currentVolumeToEngage = chain.Orders[i-1].QuoteVolumeOut
				} else if chain.Orders[i-1].TransactionType == core.Bid {
					currentVolumeToEngage = chain.Orders[i-1].BaseVolumeOut
				}
			}
			if n.IsBaseToQuote {
				chain.AdjustedVolumes[i] = currentVolumeToEngage * chain.Orders[i].Price
			} else {
				chain.AdjustedVolumes[i] = currentVolumeToEngage * chain.Orders[i].PriceOfQuoteToBase
			}
			if n.IsBaseToQuote {
				chain.Orders[i].UpdateQuoteVolume(chain.AdjustedVolumes[i])
			} else {
				chain.Orders[i].UpdateBaseVolume(chain.AdjustedVolumes[i])
			}
			chain.Cost = chain.Cost + chain.Orders[i].Fee*chain.Rates[i]
		}
		chain.Path = p

		firstOrder := chain.Orders[0]
		if firstOrder.TransactionType == core.Bid {
			chain.VolumeIn = firstOrder.QuoteVolumeIn
		} else {
			chain.VolumeIn = firstOrder.BaseVolumeIn
		}

		lastOrder := chain.Orders[len(chain.Orders)-1]
		if lastOrder.TransactionType == core.Bid {
			chain.VolumeOut = lastOrder.BaseVolumeOut
		} else {
			chain.VolumeOut = lastOrder.QuoteVolumeOut
		}
		if chain.VolumeIn < 0.0001 || chain.VolumeOut < 0.0001 {
			chain.IsBroken = true
		}
		chain.Performance = chain.VolumeOut / chain.VolumeIn
		chains[j] = chain
	}
	return chains
}
