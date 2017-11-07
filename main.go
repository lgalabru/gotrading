package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gotrading/core"
	"gotrading/graph"
	"gotrading/services"
	"gotrading/strategies"

	"github.com/thrasher-/gocryptotrader/config"
	"github.com/thrasher-/gocryptotrader/exchanges/kraken"
	"github.com/thrasher-/gocryptotrader/exchanges/liqui"
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
	liquiEngine := new(liqui.Liqui)

	kraken := services.LoadExchange(cfg, "Kraken", krakenEngine)
	// poloniex := services.LoadExchange(cfg, "Poloniex", poloniexEngine)
	liqui := services.LoadExchange(cfg, "Liqui", liquiEngine)
	exchanges := []core.Exchange{liqui, kraken}
	// exchanges := []core.Exchange{kraken}

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
	delayBetweenReqs := make(map[string]time.Duration, len(exchanges))
	delayBetweenReqs["Kraken"] = time.Duration(100)
	delayBetweenReqs["Liqui"] = time.Duration(500)

	for _, exch := range exchanges {
		nodes := pairsLookup[exch.Name]
		go services.StartPollingOrderbooks(exch, nodes, delayBetweenReqs[exch.Name], func(n graph.Node) {
			arbitrage.Run(pathsLookup[n.ID()])
		})
	}

	<-interrupt
}
