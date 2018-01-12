package arbitrage

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"gotrading/core"
	"gotrading/graph"
)

type Report struct {
	Id                            *string      `json:"id"`
	Path                          graph.Path   `json:"-"`
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
	desc := r.Path.Description()
	h := sha1.New()
	h.Write([]byte(desc))
	enc := hex.EncodeToString(h.Sum(nil))
	r.Id = &enc
	return json.Marshal(r)
}

func (r Report) Description() string {
	str := fmt.Sprintf("%s: %f", r.Path.Description(), r.VolumeOut-r.VolumeIn)
	return str
}
