package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"gotrading/core"
	"gotrading/graph"
	"gotrading/services"
	"gotrading/strategies"

	"github.com/olekukonko/tablewriter"
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

	liquiEngine := new(liqui.Liqui)
	krakenEngine := new(kraken.Kraken)
	// bittrexEngine := new(bittrex.Bittrex)
	// gdaxEngine := new(gdax.GDAX)
	// poloniexEngine := new(poloniex.Poloniex)

	liqui := services.LoadExchange(cfg, "Liqui", liquiEngine)
	kraken := services.LoadExchange(cfg, "Kraken", krakenEngine)
	// bittrex := services.LoadExchange(cfg, "Bittrex", bittrexEngine)
	// poloniex := services.LoadExchange(cfg, "Poloniex", poloniexEngine)
	// gdax := services.LoadExchange(cfg, "GDAX", gdaxEngine)

	// exchanges := []core.Exchange{kraken, liqui, gdax, bittrex}
	exchanges := []core.Exchange{kraken, liqui}

	mashup := core.ExchangeMashup{}
	mashup.Init(exchanges)

	from := core.Currency("ETH")
	to := from
	depth := 3
	nodes, paths, _ := graph.PathFinder(mashup, from, to, depth)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Link 1", "Link 2", "Link 3", "Performance", "Input", "Output", "Result"})
	table.Render()

	// Create a map
	pairsLookup := make(map[string][]graph.NodeLookup)
	for _, n := range nodes {
		paths := paths[n.ID()]
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
	delayBetweenReqs["Kraken"] = time.Duration(500)
	delayBetweenReqs["Liqui"] = time.Duration(500)

	// conn, err := amqp.Dial("amqp://developer:xLae4pzT@gotrading-rabbitmq.dev:5672/gotrading")
	// defer conn.Close()

	for _, exch := range exchanges {
		nodes := pairsLookup[exch.Name]
		go services.StartPollingOrderbooks(exch, nodes, delayBetweenReqs[exch.Name], func(n graph.Node) {
			chains := arbitrage.Run(paths[n.ID()])
			rows := make([][]string, 0)
			for _, chain := range chains {
				if chain.Performance == 0 {
					continue
				}
				ordersCount := len(chain.Path.ContextualNodes)
				row := make([]string, ordersCount+4)
				for j, node := range chain.Path.ContextualNodes {
					row[j] = node.Description()
				}
				row[ordersCount] = strconv.FormatFloat(chain.Performance, 'f', 6, 64)
				row[ordersCount+1] = strconv.FormatFloat(chain.Volume, 'f', 6, 64)
				row[ordersCount+2] = strconv.FormatFloat(chain.Volume*chain.Performance, 'f', 6, 64)
				row[ordersCount+3] = strconv.FormatFloat(chain.Volume*chain.Performance-chain.Volume, 'f', 6, 64)
				rows = append(rows, row)
			}
			if len(rows) > 0 {
				table.AppendBulk(rows)
				table.Render()
			}
		})
	}

	<-interrupt
}
