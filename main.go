package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gotrading/core"
	"gotrading/graph"
	"gotrading/strategies"

	"github.com/thrasher-/gocryptotrader/config"
	"github.com/thrasher-/gocryptotrader/currency/pair"
	"github.com/thrasher-/gocryptotrader/exchanges/kraken"
)

func main() {

	cfg := &config.Cfg
	err := cfg.LoadConfig("config.dat")
	if err != nil {
		log.Fatal(err)
	}

	interrupt := make(chan os.Signal, 1)

	// portfolio := core.Portfolio{}
	// portfolio.Init(currencies, exchanges)
	// portfolio.DidBuy(0, 7000, core.CurrencyPair{core.Currency("USD"), core.Currency("USD")}, core.Exchange{"Alpha"})

	// order1 := core.Order{6000, 1, core.Sell}
	// portfolio.Fullfill(order1, 1, currencyPair, core.Exchange{"Alpha"})
	//
	// order2 := core.Order{8000, 1, core.Buy}
	// portfolio.Fullfill(order2, 1, currencyPair, core.Exchange{"Alpha"})

	// BTC/USD: 6950
	// ETH/USD: 280
	// ETH/BTC: 0.040

	// portfolio.DidBuy(0, 10, )

	// portfolio.DidSold(0, 10, core.CurrencyPair{core.Currency("BTC"), core.Currency("USD")}, core.Exchange{"Alpha"})
	// portfolio.DisplayBalances()
	krakenEngine := new(kraken.Kraken)
	// poloniexEngine := new(poloniex.Poloniex)
	// liquiEngine := new(liqui.Liqui)

	kraken := LoadExchange(cfg, "Kraken", krakenEngine)
	// poloniex := LoadExchange(cfg, "Poloniex", poloniexEngine)
	// liqui := LoadExchange(cfg, "Liqui", liquiEngine)
	// exchanges := []core.Exchange{poloniex, liqui, kraken}
	exchanges := []core.Exchange{kraken}

	mashup := core.ExchangeMashup{}
	mashup.Init(exchanges)

	from := core.Currency("ETH")
	to := core.Currency("ETH")
	depth := 3
	paths := graph.PathFinder(mashup, from, to, depth)

	nodes := make([]*graph.Node, 0)
	pathsLookup := make(map[string][]graph.Path)
	for _, path := range paths {
		for _, cn := range path.ContextualNodes {
			paths, ok := pathsLookup[cn.Node.ID()]
			if !ok {
				nodes = append(nodes, cn.Node)
				paths = make([]graph.Path, 0)
			}
			pathsLookup[cn.Node.ID()] = append(paths, path)
			// fmt.Println(path.Description())
		}
	}
	fmt.Println("Observing", len(paths), "combinations, distributed over", len(nodes), "pairs.")

	pairsLookup := make(map[string][]graph.NodeLookup)
	for _, n := range nodes {
		paths := pathsLookup[n.ID()]
		lookups, ok := pairsLookup[n.Exchange.Name]
		if !ok {
			lookups = make([]graph.NodeLookup, 0)
		}
		lookup := graph.NodeLookup{n, len(paths)}
		pairsLookup[n.Exchange.Name] = append(lookups, lookup)
	}

	for _, exch := range exchanges {
		pairsLookup[exch.Name] = graph.MergeSort(pairsLookup[exch.Name])
	}

	arbitrage := strategies.Arbitrage{}

	for i := 0; i < 10; i++ {
		for _, exch := range exchanges {
			sortedNodes := pairsLookup[exch.Name]
			if i >= len(sortedNodes)-1 {
				continue
			}
			n := sortedNodes[i]
			cp := pair.NewCurrencyPair(string(n.Node.From), string(n.Node.To))
			fmt.Println("===================")
			fmt.Println(n.Node.From, n.Node.To, exch.Name)
			fmt.Println("===================")
			src, err := exch.Engine.UpdateOrderbook(cp, "SPOT")
			if err == nil {
				dst := n.Node.Orderbook
				if dst == nil {
					dst = &core.Orderbook{}
					dst.CurrencyPair = core.CurrencyPair{n.Node.From, n.Node.To}
					dst.Bids = make([]core.Order, 0)
					dst.Asks = make([]core.Order, 0)
					n.Node.Orderbook = dst
				}
				// fmt.Println("1 ------------------")
				// fmt.Println(src.Asks)
				for _, ask := range src.Asks {
					dst.Asks = append(dst.Asks, core.Order{ask.Price, ask.Amount, core.Sell})
				}
				// fmt.Println("2 ------------------")
				// fmt.Println(src.Bids)
				for _, bid := range src.Bids {
					dst.Bids = append(dst.Bids, core.Order{bid.Price, bid.Amount, core.Buy})
				}
				// fmt.Println("~~~~~~~~~~~~~~~~~~")

				arbitrage.Run(pathsLookup[n.Node.ID()])
			} else {
				fmt.Println("Error", err)
			}
		}
		time.Sleep(5000 * time.Millisecond)
	}
	<-interrupt
}

func LoadExchange(cfg *config.Config, name string, exch core.ExchangeInterface) core.Exchange {
	config, _ := cfg.GetExchangeConfig(name)
	exch.SetDefaults()
	exch.Setup(config)
	var rawPairs = exch.GetAvailableCurrencies()
	pairs := make([]core.CurrencyPair, len(rawPairs))
	for i, c := range rawPairs {
		pairs[i] = core.CurrencyPair{
			core.Currency(c.GetFirstCurrency()),
			core.Currency(c.GetSecondCurrency())}
	}
	return core.Exchange{name, pairs, exch}
}
