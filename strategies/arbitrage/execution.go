package arbitrage

import (
	"fmt"
	"math"
	"strings"
	"time"

	"gotrading/core"
	"gotrading/graph"
	"gotrading/networking"
)

type Execution struct {
	StartedAt    time.Time `json:"startedAt"`
	EndedAt      time.Time `json:"endedAt"`
	Report       Report
	gatling      *networking.Gatling
	simulation   Simulation
	chain        ChainedOrders
	IsSuccessful bool
}

type ChainedOrders struct {
	Path             graph.Path   `json:"path"`
	Orders           []core.Order `json:"orders"`
	Performance      float64      `json:"performance"`
	Rates            []float64    `json:"rates"`
	AdjustedVolumes  []float64    `json:"volumes"`
	VolumeToEngage   float64      `json:"volumeToEngage"`
	VolumeIn         float64      `json:"volumeIn"`
	VolumeOut        float64      `json:"volumeOut"`
	Cost             float64      `json:"cost"`
	Results          []string     `json:"results"`
	IsBroken         bool         `json:"is_broken"`
	CreatedAt        time.Time    `json:"createdAt"`
	DiagnosedAt      time.Time    `json:"diagnosedAt"`
	StartedTradingAt time.Time    `json:"startedTradingAt"`
	EndedTradingAt   time.Time    `json:"endedTradingAt"`
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func (exec *Execution) Init(sim Simulation) {
	exec.simulation = sim
	exec.chain = sim.chain
}

func (exec *Execution) Run() {

	exec.IsSuccessful = true
	exec.chain.Results = make([]string, len(exec.chain.Orders))
	exec.chain.StartedTradingAt = time.Now()
	for i, o := range exec.chain.Orders {
		exchange := exec.chain.Path.Hits[i].Endpoint.Exchange
		pair := strings.ToLower(string(exec.chain.Path.Hits[i].Endpoint.From)) + "_" + strings.ToLower(string(exec.chain.Path.Hits[i].Endpoint.To))
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
			exec.chain.Results[i] = error.Error()
			exec.chain.EndedTradingAt = time.Now()
			exec.IsSuccessful = false
			return
		} else {
			exec.chain.Results[i] = "ok" // Taker? Maker?
		}
	}
	exec.chain.EndedTradingAt = time.Now()
}

func (exec *Execution) BuildReport() Report {
	return Report{}
}
