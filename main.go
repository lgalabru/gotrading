package main

import (
	"fmt"
	"strconv"
	"strings"

	"gotrading/core"
	"gotrading/exchanges"
	"gotrading/graph"
	"gotrading/reporting"
	"gotrading/strategies/arbitrage"

	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	factory := exchanges.Factory{}
	exchanges := []core.Exchange{}

	arbitrageSettings := viper.GetStringMapString("strategies.arbitrage")
	exchangesEnabled := strings.Split(arbitrageSettings["exchanges_enabled"], ",")

	for _, name := range exchangesEnabled {
		exch := factory.BuildExchange(name)
		exchanges = append(exchanges, exch)
	}

	mashup := core.ExchangeMashup{}
	mashup.Init(exchanges)

	from := core.Currency(arbitrageSettings["from_currency"])
	to := core.Currency(arbitrageSettings["to_currency"])
	depth, _ := strconv.Atoi(arbitrageSettings["shifts_count"])
	treeOfPossibles, _, _, _ := graph.PathFinder(mashup, from, to, depth)

	publisher := reporting.Publisher{}
	publisher.Init(viper.GetStringMapString("strategies.arbitrage.reporting.publisher"))
	defer publisher.Close()

	for {
		treeOfPossibles.DepthTraversing(func(hits []*core.Hit) {
			sim := arbitrage.Simulation{}
			sim.Init(hits)
			sim.Run()
			if sim.IsSuccessful() == false {
				if sim.IsIncomplete() == false {
					go publisher.Send(sim.Report)
				}
				return
			}

			exec := arbitrage.Execution{}
			exec.Init(sim)
			exec.Run()
			if exec.IsSuccessful() == false {
				go publisher.Send(exec.Report)
				// Recovery? Rollback?
				return
			}

			valid := arbitrage.Validation{}
			valid.Init(exec)
			valid.Run()
			go publisher.Send(verif.Report)
		})
	}
}
