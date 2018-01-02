package arbitrage

import (
  "time"
)

type Report struct {
  StartedAt time.Time    `json:"startedAt"`
	EndedAt   time.Time    `json:"endedAt"`
}

func (r Report) Encoded() {
}
