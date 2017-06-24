package histories

import (
	"time"

	"github.com/jinzhu/gorm"
)

type LendingRow struct {
	gorm.Model

	LendingID int64     `gorm:"unique" json:"lending_id"`
	Currency  string    `json:"currency"`
	Rate      float64   `json:"rate,string"`
	Amount    float64   `json:"amount,string"`
	Duration  float64   `json:"duration,string"`
	Interest  float64   `json:"interest,string"`
	Fee       float64   `json:"fee,string"`
	Earned    float64   `json:"earned,string"`
	Open      time.Time `json:"open"`
	Close     time.Time `json:"close"`
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

func (r *LendingRow) MakeHistory() PoloniexLendingHistory {
	return PoloniexLendingHistory{
		ID:       r.LendingID,
		Currency: r.Currency,
		Rate:     r.Rate,
		Amount:   r.Amount,
		Duration: r.Duration,
		Interest: r.Interest,
		Fee:      r.Fee,
		Earned:   r.Earned,
		Open:     r.Open.Format(time.RFC3339),
		Close:    r.Close.Format(time.RFC3339),
	}
}
