package core

import (
	"time"
)

type Portfolio struct {
	CurrencyPair CurrencyPair
	LastUpdated  time.Time
}
