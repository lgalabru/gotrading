package arbitrage

import (
	"gotrading/core"
	"gotrading/networking"
	"math"
	"time"
)

type Execution struct {
	Report Report
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
	m := core.SharedPortfolioManager()

	r := &exec.Report
	r.IsExecutionSuccessful = true
	r.ExecutionStartedAt = time.Now()
	r.PreExecutionPortfolioStateID = m.LastStateID

	batch := networking.Batch{}
	batch.PostOrders(r.Orders, func(dispatched []core.OrderDispatched) {

		r.DispatchedOrders = dispatched
		r.PostExecutionPortfolioStateID = m.LastStateID
		r.ExecutionEndedAt = time.Now()
	})

	// for i, o := range r.Orders {
	// 	exchange := o.Hit.Endpoint.Exchange
	// 	exchange.PostOrder(o)

	// 	if error != nil {
	// 		r.Results[i] = error.Error()
	// 		r.ExecutionEndedAt = time.Now()
	// 		r.IsExecutionSuccessful = false
	// 		return
	// 	} else {
	// 		r.Results[i] = "hit" // Taker? Maker?
	// 	}
	// }
}

func (exec *Execution) IsSuccessful() bool {
	return exec.Report.IsExecutionSuccessful
}
