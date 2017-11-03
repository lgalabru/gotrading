package core

type OrderType uint

const (
	Buy OrderType = iota
	Sell
)

type Order struct {
	Price  float64
	Volume float64
	Type   OrderType
}
