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

	from := core.Currency("BTC")
	to := from
	depth := 3
	nodes, paths := graph.PathFinder(mashup, from, to, depth)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Link 1", "Link 2", "Link 3", "Values", "Performance", "Trigger"})
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
	delayBetweenReqs["Kraken"] = time.Duration(100)
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
				row := make([]string, ordersCount+3)
				for j, node := range chain.Path.ContextualNodes {
					row[j] = node.Description()
				}
				row[ordersCount] = "-"
				row[ordersCount+1] = strconv.FormatFloat(chain.Performance, 'f', 6, 64)
				row[ordersCount+2] = "Refresh"
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
