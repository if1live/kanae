package histories

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/thrasher-/gocryptotrader/exchanges/poloniex"
)

func convertPoloniexDate(val string) time.Time {
	// date example : 2017-06-18 04:31:08
	t, _ := time.Parse("2006-01-02 15:04:05", val)
	return t
}

// PoloniexAuthenticatedTradeHistory
type PoloniexTradeRow struct {
	gorm.Model

	// from zenbot format
	// {exchange_slug}.{asset}-{currency}
	// poloniex.AMP-BTC
	Asset    string
	Currency string

	GlobalTradeID int64 `gorm:"unique"`
	TradeID       int64 `gorm:"unique"`
	Date          time.Time
	Rate          float64
	Amount        float64
	Total         float64
	Fee           float64
	OrderNumber   int64
	Type          string
	Category      string
}

func NewPoloniexTradeRow(asset, currency string, h poloniex.PoloniexAuthentictedTradeHistory) PoloniexTradeRow {
	return PoloniexTradeRow{
		Asset:    asset,
		Currency: currency,

		GlobalTradeID: h.GlobalTradeID,
		TradeID:       h.TradeID,
		Date:          convertPoloniexDate(h.Date),
		Rate:          h.Rate,
		Amount:        h.Amount,
		Total:         h.Total,
		Fee:           h.Fee,
		OrderNumber:   h.OrderNumber,
		Type:          h.Type,
		Category:      h.Category,
	}
}

type PoloniexLendingRow struct {
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

func NewPoloniexLendingRow(h PoloniexLendingHistory) PoloniexLendingRow {
	return PoloniexLendingRow{
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
