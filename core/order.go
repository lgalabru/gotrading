package core

type orderType uint

const (
	buy orderType = iota
	sell
)

type Order struct {
	Type   orderType
	Volume float64
	Price  float64
}
