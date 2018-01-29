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
	Id                            *string      `json:"id"`
	Orders                        []core.Order `json:"orders"`
	Performance                   float64      `json:"performance"`
	Rates                         []float64    `json:"rates"`
	AdjustedVolumes               []float64    `json:"volumes"`
	VolumeToEngage                float64      `json:"volumeToEngage"`
	VolumeIn                      float64      `json:"volumeIn"`
	VolumeOut                     float64      `json:"volumeOut"`
	Cost                          float64      `json:"cost"`
	Results                       []string     `json:"results"`
	IsTradedVolumeEnough          bool         `json:"isTradedVolumeEnough"`
	SimulationStartedAt           time.Time    `json:"simulationStartedAt"`
	SimulationComputingStartedAt  time.Time    `json:"simulationComputingStartedAt"`
	SimulationEndedAt             time.Time    `json:"simulationEndedAt"`
	IsSimulationIncomplete        bool         `json:"isSimulationIncomplete"`
	IsSimulationSuccessful        bool         `json:"isSimulationSuccessful"`
	PreExecutionPortfolioStateID  string       `json:"preExecutionPortfolioStateID"`
	ExecutionStartedAt            time.Time    `json:"executionStartedAt"`
	ExecutionEndedAt              time.Time    `json:"executionEndedAt"`
	PostExecutionPortfolioStateID string       `json:"postExecutionPortfolioStateID"`
	IsExecutionSuccessful         bool         `json:"isExecutionSuccessful"`
	ValidationStartedAt           time.Time    `json:"validationStartedAt"`
	ValidationEndedAt             time.Time    `json:"validationEndedAt"`
	SimulationMinusExecution      float64      `json:"simulationMinusExecution"`
}

func (r Report) Encode() ([]byte, error) {
	desc := "Report"
	for _, o := range r.Orders {
		if o.Hit == nil {
			return nil, fmt.Errorf("Order incomplete")
		}
		desc = desc + " -> " + o.Hit.Endpoint.Description()
	}
	h := sha1.New()
	h.Write([]byte(desc))
	enc := hex.EncodeToString(h.Sum(nil))
	r.Id = &enc
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
