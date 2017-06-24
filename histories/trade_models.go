package histories

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/thrasher-/gocryptotrader/exchanges/poloniex"
)

// PoloniexAuthenticatedTradeHistory
type TradeRow struct {
	gorm.Model

	// from zenbot format
	// {exchange_slug}.{asset}-{currency}
	// poloniex.AMP-BTC
	Asset    string
	Currency string

	GlobalTradeID int64 `gorm:"unique"`
	TradeID       int64
	Date          time.Time
	Rate          float64
	Amount        float64
	Total         float64
	Fee           float64
	OrderNumber   int64
	Type          string
	Category      string
}

func NewTradeRow(asset, currency string, h poloniex.PoloniexAuthentictedTradeHistory) TradeRow {
	return TradeRow{
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
