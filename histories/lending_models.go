package histories

import (
	"time"

	"github.com/jinzhu/gorm"
)

type LendingRow struct {
	gorm.Model

	LendingID int64 `gorm:"unique"`
	Currency  string
	Rate      float64
	Amount    float64
	Duration  float64
	Interest  float64
	Fee       float64
	Earned    float64
	Open      time.Time
	Close     time.Time
}

func NewLendingRow(h PoloniexLendingHistory) LendingRow {
	return LendingRow{
		LendingID: h.ID,
		Currency:  h.Currency,
		Rate:      h.Rate,
		Amount:    h.Amount,
		Duration:  h.Duration,
		Interest:  h.Interest,
		Fee:       h.Fee,
		Earned:    h.Earned,
		Open:      convertPoloniexDate(h.Open),
		Close:     convertPoloniexDate(h.Close),
	}
}
