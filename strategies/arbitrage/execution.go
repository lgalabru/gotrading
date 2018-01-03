package arbitrage

import (
	"fmt"
	"math"
	"strings"
	"time"

	"gotrading/core"
)

type Execution struct {
	Report     Report
	simulation Simulation
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func (exec *Execution) Init(sim Simulation) {
	exec.Report = sim.Report
}

func (exec *Execution) Run() {
	r := &exec.Report
	r.IsExecutionSuccessful = true
	r.Results = make([]string, len(r.Orders))
	r.ExecutionStartedAt = time.Now()
	for i, o := range r.Orders {
		exchange := r.Path.Hits[i].Endpoint.Exchange
		pair := strings.ToLower(string(r.Path.Hits[i].Endpoint.From)) + "_" + strings.ToLower(string(r.Path.Hits[i].Endpoint.To))
		var orderType string
		var amount float64

		if o.TransactionType == core.Ask {
			orderType = "sell"
			amount = o.BaseVolumeIn
		} else {
			orderType = "buy"
			amount = o.QuoteVolumeIn / o.Price
		}
		price := o.Price
		// decimals := exec.chain.Path.Hits[i].Endpoint.Exchange.Liqui.Info.Pairs[pair].DecimalPlaces
		decimals := 8
		res, error := exchange.PostOrder(o)

		// res, error := exchange.Trade(pair, orderType, toFixed(amount, decimals), price)
		fmt.Println("Executing order:", pair, orderType, decimals, toFixed(amount, decimals), price, res, error)
		if error != nil {
			r.Results[i] = error.Error()
			r.ExecutionEndedAt = time.Now()
			r.IsExecutionSuccessful = false
			return
		} else {
			r.Results[i] = "hit" // Taker? Maker?
		}
	}
	r.ExecutionEndedAt = time.Now()
}

func (exec *Execution) IsSuccessful() bool {
	return exec.Report.IsExecutionSuccessful
}
