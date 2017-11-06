package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gotrading/core"
	"gotrading/exchanges/orderbook"
	"gotrading/graph"

	"github.com/thrasher-/gocryptotrader/config"
	"github.com/thrasher-/gocryptotrader/currency/pair"
	"github.com/thrasher-/gocryptotrader/exchanges/kraken"
	"github.com/thrasher-/gocryptotrader/exchanges/liqui"
	"github.com/thrasher-/gocryptotrader/exchanges/poloniex"
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
	poloniexEngine := new(poloniex.Poloniex)
	liquiEngine := new(liqui.Liqui)

	kraken := LoadExchange(cfg, "Kraken", krakenEngine)
	poloniex := LoadExchange(cfg, "Poloniex", poloniexEngine)
	liqui := LoadExchange(cfg, "Liqui", liquiEngine)
	exchanges := []core.Exchange{poloniex, liqui, kraken}

	mashup := core.ExchangeMashup{}
	mashup.Init(exchanges)

	from := core.Currency("BTC")
	to := core.Currency("BTC")
	depth := 3
	paths := graph.PathFinder(mashup, from, to, depth)

	nodes := make([]graph.Node, 0)
	pathsLookup := make(map[string][]graph.Path)
	for _, path := range paths {
		for _, node := range path.Nodes {
			paths, ok := pathsLookup[node.ID()]
			if !ok {
				nodes = append(nodes, node)
				paths = make([]graph.Path, 0)
			}
			pathsLookup[node.ID()] = append(paths, path)
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

	for _, exch := range exchanges {
		sortedNodes := pairsLookup[exch.Name]
		for i, n := range sortedNodes {
			if i > 5 {
				break
			}
			cp := pair.NewCurrencyPair(string(n.Node.From), string(n.Node.To))

			base := orderbook.Base{
				Pair:         cp,
				CurrencyPair: cp.Pair().String(),
				Asks:         []orderbook.Item{orderbook.Item{Price: 0, Amount: 0}},
				Bids:         []orderbook.Item{orderbook.Item{Price: 0, Amount: 0}},
			}

			o1 := orderbook.CreateNewOrderbook(exch.Name, cp, base, orderbook.Spot)
			fmt.Println(o1.Orderbook)
			o, err := exch.Engine.UpdateOrderbook(cp, "SPOT")
			if err != nil {
				fmt.Println(o)
			}
			time.Sleep(10000 * time.Millisecond)
		}
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
