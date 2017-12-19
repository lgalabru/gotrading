package services

import (
	"fmt"
	"gotrading/core"
	"gotrading/graph"
	"math"
	"strings"
)

type ChainedOrders struct {
	Path            graph.Path   `json:"path"`
	Orders          []core.Order `json:"orders"`
	Performance     float64      `json:"performance"`
	Rates           []float64    `json:"rates"`
	AdjustedVolumes []float64    `json:"volumes"`
	VolumeToEngage  float64      `json:"volumeToEngage"`
	VolumeIn        float64      `json:"volumeIn"`
	VolumeOut       float64      `json:"volumeOut"`
	Cost            float64      `json:"cost"`
	Results         []string     `json:"results"`
	IsBroken        bool         `json:"is_broken"`
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func (c *ChainedOrders) Execute() bool {
	chainExecutedWithoutError := true
	c.Results = make([]string, len(c.Orders))
	for i, o := range c.Orders {
		engine := *c.Path.Nodes[i].Endpoint.Exchange.Engine
		pair := strings.ToLower(string(c.Path.Nodes[i].Endpoint.From)) + "_" + strings.ToLower(string(c.Path.Nodes[i].Endpoint.To))
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
		decimals := c.Path.Nodes[i].Endpoint.Exchange.Liqui.Info.Pairs[pair].DecimalPlaces

		res, error := engine.Trade(pair, orderType, toFixed(amount, decimals), price)
		fmt.Println("Executing order:", pair, orderType, toFixed(amount, decimals), price, res, error)
		if error != nil {
			c.Results[i] = error.Error()
			return false
		} else {
			c.Results[i] = "ok" // Taker? Maker?
		}
	}
	return chainExecutedWithoutError
}
