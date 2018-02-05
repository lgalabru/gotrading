package core

import "time"

// OrderDispatched represents an order
type OrderDispatched struct {
	Order       *Order    `json:"order"`
	SentAt      time.Time `json:"sentAt"`
	ConfirmedAt time.Time `json:"confirmedAt"`
	IsSuccess   bool      `json:"isSuccess"`
	Message     string    `json:"error"`
}
