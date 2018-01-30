package arbitrage

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"gotrading/core"
)

type Report struct {
	Id                            *string                `json:"id"`
	Orders                        []core.Order           `json:"orders"`
	DispatchedOrders              []core.OrderDispatched `json:"dispatchedOrders"`
	Performance                   float64                `json:"performance"`
	Rates                         []float64              `json:"rates"`
	AdjustedVolumes               []float64              `json:"volumes"`
	VolumeToEngage                float64                `json:"volumeToEngage"`
	VolumeIn                      float64                `json:"volumeIn"`
	VolumeOut                     float64                `json:"volumeOut"`
	Cost                          float64                `json:"cost"`
	IsTradedVolumeEnough          bool                   `json:"isTradedVolumeEnough"`
	SimulationStartedAt           time.Time              `json:"simulationStartedAt"`
	SimulationFetchingStartedAt   time.Time              `json:"simulationFetchingStartedAt"`
	SimulationComputingStartedAt  time.Time              `json:"simulationComputingStartedAt"`
	SimulationEndedAt             time.Time              `json:"simulationEndedAt"`
	IsSimulationIncomplete        bool                   `json:"isSimulationIncomplete"`
	IsSimulationSuccessful        bool                   `json:"isSimulationSuccessful"`
	ExecutionStartedAt            time.Time              `json:"executionStartedAt"`
	ExecutionEndedAt              time.Time              `json:"executionEndedAt"`
	IsExecutionSuccessful         bool                   `json:"isExecutionSuccessful"`
	ValidationStartedAt           time.Time              `json:"validationStartedAt"`
	ValidationEndedAt             time.Time              `json:"validationEndedAt"`
	SimulationMinusExecution      float64                `json:"simulationMinusExecution"`
	PreExecutionPortfolioStateID  string                 `json:"-"`
	PostExecutionPortfolioStateID string                 `json:"-"`
	PreExecutionPortfolioState    core.PortfolioState    `json:"statePreExecution"`
	PostExecutionPortfolioState   core.PortfolioState    `json:"statePostExecution"`
}

func (r Report) Encode() ([]byte, error) {
	desc := "Report"
	r.SimulationFetchingStartedAt = r.SimulationComputingStartedAt
	for _, o := range r.Orders {
		if o.Hit == nil {
			return nil, fmt.Errorf("Order incomplete")
		}
		desc = desc + " -> " + o.Hit.Endpoint.Description()
		if o.Hit.Endpoint.Orderbook.StartedLastUpdateAt.Before(r.SimulationFetchingStartedAt) {
			r.SimulationFetchingStartedAt = o.Hit.Endpoint.Orderbook.StartedLastUpdateAt
		}
	}
	h := sha1.New()
	h.Write([]byte(desc))
	enc := hex.EncodeToString(h.Sum(nil))
	r.Id = &enc

	r.PreExecutionPortfolioState = core.SharedPortfolioManager().States[r.PreExecutionPortfolioStateID]
	r.PostExecutionPortfolioState = core.SharedPortfolioManager().States[r.PostExecutionPortfolioStateID]

	return json.Marshal(r)
}

func (r Report) Description() string {
	desc := "Report"
	for _, o := range r.Orders {
		var link string
		if o.Hit == nil {
			link = "Missing link"
		} else {
			link = o.Hit.Endpoint.Description()
		}
		desc = desc + " -> " + link
	}
	return desc
}
