package histories

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/thrasher-/gocryptotrader/exchanges/poloniex"
)

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
	// date example : 2017-06-18 04:31:08
	t, _ := time.Parse("2006-01-02 15:04:05", h.Date)

	return PoloniexTradeRow{
		Asset:    asset,
		Currency: currency,

		GlobalTradeID: h.GlobalTradeID,
		TradeID:       h.TradeID,
		Date:          t,
		Rate:          h.Rate,
		Amount:        h.Amount,
		Total:         h.Total,
		Fee:           h.Fee,
		OrderNumber:   h.OrderNumber,
		Type:          h.Type,
		Category:      h.Category,
	}
}

type LendRow struct {
}
