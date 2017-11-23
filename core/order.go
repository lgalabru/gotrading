package core

type OrderType uint

const (
	Buy OrderType = iota
	Sell
)

type Order struct {
	Price  float64 `json:"price"`
	Volume float64 `json:"volume"`
	Type   OrderType
}
