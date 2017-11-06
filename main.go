package main

import (
	"fmt"
	"log"
	"os"

	"gotrading/core"
	"gotrading/graph"

	"github.com/thrasher-/gocryptotrader/config"
	"github.com/thrasher-/gocryptotrader/currency/pair"
	"github.com/thrasher-/gocryptotrader/exchanges/kraken"
	"github.com/thrasher-/gocryptotrader/exchanges/liqui"
	"github.com/thrasher-/gocryptotrader/exchanges/orderbook"
	"github.com/thrasher-/gocryptotrader/exchanges/poloniex"
	"github.com/thrasher-/gocryptotrader/exchanges/ticker"
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

	kraken := LoadExchange(cfg, "Kraken", new(kraken.Kraken))
	poloniex := LoadExchange(cfg, "Poloniex", new(poloniex.Poloniex))
	liqui := LoadExchange(cfg, "Liqui", new(liqui.Liqui))
	exchanges := []core.Exchange{kraken, poloniex, liqui}

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
		}
	}
	fmt.Println("Observing", len(paths), "combinations, distributed over", len(nodes), "pairs.")

	pairsLookup := make(map[string][]graph.NodeLookup)
	op := 0
	for _, n := range nodes {
		paths := pathsLookup[n.ID()]
		lookups, ok := pairsLookup[n.Exchange.Name]
		if !ok {
			lookups = make([]graph.NodeLookup, 0)
			op += 1
		}
		lookup := graph.NodeLookup{n, len(paths)}
		pairsLookup[n.Exchange.Name] = append(lookups, lookup)

	}
	fmt.Println(op)

	for _, exch := range exchanges {
		pairsLookup[exch.Name] = graph.MergeSort(pairsLookup[exch.Name])
		fmt.Println("Merge sorting", exch.Name, len(pairsLookup[exch.Name]))
		for _, nodeLookup := range pairsLookup[exch.Name] {
			fmt.Println(nodeLookup)
		}
	}

	// [–]

	// combinationsLookup := make(map[*core.Exchange][]graph.Path)
	// pairsLookup: Exchange => [] sorted pairs,
	// pathsLookup: Node => [] sorted paths,
	// sortedNodes:= make([]graph.NodeLookup, len(nodes))

	// On souhaite obtenir les Pairs

	// for k, v := range pathsLookup {
	// 	fmt.Println(k, len(v))
	// }

	for _, p := range pathsLookup["XID-ETH@Liqui"] {
		p.Display()
	}

	// fmt.Println(pathsLookup["XIDETHLiqui"])

	// cascading := []
	// reverseLookup[node] -> path

	// c1 := pair.NewCurrencyPair("BTC", "USD")
	// base := orderbook.Base{
	// 	Pair:         c1,
	// 	CurrencyPair: c1.Pair().String(),
	// 	Asks:         []orderbook.Item{orderbook.Item{Price: 100, Amount: 10}},
	// 	Bids:         []orderbook.Item{orderbook.Item{Price: 200, Amount: 10}},
	// }

	// o1 := orderbook.CreateNewOrderbook("Kraken", c1, base, orderbook.Spot)
	// fmt.Println(o1.Orderbook)
	// o, err := krakenExchange.UpdateOrderbook(rawKrakenPairs[0], "SPOT")
	// if err != nil {
	// 	fmt.Println(o)
	// }

	<-interrupt
}

type ExchangeInterface interface {
	Setup(exch config.ExchangeConfig)
	Start()
	SetDefaults()
	GetName() string
	IsEnabled() bool
	GetTickerPrice(currency pair.CurrencyPair, assetType string) (ticker.Price, error)
	UpdateTicker(currency pair.CurrencyPair, assetType string) (ticker.Price, error)
	GetOrderbookEx(currency pair.CurrencyPair, assetType string) (orderbook.Base, error)
	UpdateOrderbook(currency pair.CurrencyPair, assetType string) (orderbook.Base, error)
	GetEnabledCurrencies() []pair.CurrencyPair
	GetAuthenticatedAPISupport() bool
	GetAvailableCurrencies() []pair.CurrencyPair
}

func LoadExchange(cfg *config.Config, name string, exch ExchangeInterface) core.Exchange {
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
	return core.Exchange{name, pairs}
}
