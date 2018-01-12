package arbitrage

import (
	"gotrading/core"
	"time"
)

type Validation struct {
	Report Report
}

func (v *Validation) Init(exec Execution) {
	v.Report = exec.Report
}

func (v *Validation) Run() {
	m := core.SharedPortfolioManager()
	r.ValidationStartedAt = time.Now()

	r := &v.Report

	simIn := r.VolumeIn
	simOut := r.VolumeOut

	firstHit := r.Orders[0].Hit
	lastHit := r.Orders[2].Hit
	execIn := m.Position(r.PreExecutionPortfolioStateID, firstHit.Endpoint.Exchange.Name, core.Currency("BTC"))
	execOut := m.Position(r.PostExecutionPortfolioStateID, lastHit.Endpoint.Exchange.Name, core.Currency("BTC"))
	r.SimulationMinusExecution = (simOut - simIn) - (execOut - execIn)
	r.ValidationEndedAt = time.Now()
}