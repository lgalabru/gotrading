package arbitrage

import (
	"fmt"
	"gotrading/core"
)

type Validation struct {
	Report Report
}

func (v *Validation) Init(exec Execution) {
	v.Report = exec.Report
}

func (v *Validation) Run() {
	m := core.SharedPortfolioManager()

	r := &v.Report
	firstHit := r.Orders[0].Hit
	lastHit := r.Orders[2].Hit

	simIn := r.VolumeIn
	simOut := r.VolumeOut

	// execIn := m.Position(r.PreExecutionPortfolioStateID, firstHit.Endpoint.Exchange.Name, firstHit.SoldCurrency)
	// execOut := m.Position(r.PreExecutionPortfolioStateID, lastHit.Endpoint.Exchange.Name, firstHit.BoughtCurrency)
	execIn := m.Position(r.PreExecutionPortfolioStateID, firstHit.Endpoint.Exchange.Name, core.Currency("BTC"))
	execOut := m.Position(r.PostExecutionPortfolioStateID, lastHit.Endpoint.Exchange.Name, core.Currency("BTC"))

	fmt.Println("Sim:", simIn, simOut, simOut-simIn)
	fmt.Println("Exec:", execIn, execOut, execOut-execIn)
	fmt.Println("---------")
	fmt.Println("Pre:", r.PreExecutionPortfolioStateID)
	fmt.Println("Post:", r.PostExecutionPortfolioStateID)

	fmt.Println(m)
}

func (v *Validation) IsSuccessful() bool {
	return v.Report.IsVerificationSuccessful
}
