package arbitrage

import (
	"encoding/json"
	"time"

	"gotrading/core"
	"gotrading/graph"
)

type Report struct {
	Path                     graph.Path   `json:"-"`
	Orders                   []core.Order `json:"orders"`
	Performance              float64      `json:"performance"`
	Rates                    []float64    `json:"rates"`
	AdjustedVolumes          []float64    `json:"volumes"`
	VolumeToEngage           float64      `json:"volumeToEngage"`
	VolumeIn                 float64      `json:"volumeIn"`
	VolumeOut                float64      `json:"volumeOut"`
	Cost                     float64      `json:"cost"`
	Results                  []string     `json:"results"`
	IsTradedVolumeEnough     bool         `json:"isTradedVolumeEnough"`
	SimulationStartedAt      time.Time    `json:"simulationStartedAt"`
	SimulationEndedAt        time.Time    `json:"simulationEndedAt"`
	IsSimulationIncomplete   bool         `json:"isSimulationIncomplete"`
	IsSimulationSuccessful   bool         `json:"isSimulationSuccessful"`
	ExecutionStartedAt       time.Time    `json:"executionStartedAt"`
	ExecutionEndedAt         time.Time    `json:"executionEndedAt"`
	IsExecutionSuccessful    bool         `json:"isExecutionSuccessful"`
	VerificationStartedAt    time.Time    `json:"verificationStartedAt"`
	VerificationEndedAt      time.Time    `json:"verificationEndedAt"`
	IsVerificationSuccessful bool         `json:"isVerificationSuccessful"`
}

func (r Report) Encode() ([]byte, error) {
	return json.Marshal(r)
}
