package models

import "github.com/jinzhu/gorm"

type Trade struct {
	gorm.Model
	// DateClosed
	// PriceSold
	// ClosingProfit $
	// ClosingProfit %
	// DailyProfit $
	// DailyProfit %
}
