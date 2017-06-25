package histories

import (
	"time"

	"github.com/if1live/kanae/kanaelib"
	"github.com/jinzhu/gorm"
	"github.com/thrasher-/gocryptotrader/exchanges/poloniex"
)

const (
	TradeSell = "sell"
	TradeBuy  = "buy"
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

func (r *TradeRow) MakeHistory() poloniex.PoloniexAuthentictedTradeHistory {
	return poloniex.PoloniexAuthentictedTradeHistory{
		GlobalTradeID: r.GlobalTradeID,
		TradeID:       r.TradeID,
		Date:          r.Date.Format(time.RFC3339),
		Rate:          r.Rate,
		Amount:        r.Amount,
		Total:         r.Total,
		Fee:           r.Fee,
		OrderNumber:   r.OrderNumber,
		Type:          r.Type,
		Category:      r.Category,
	}
}

func (r *TradeRow) FeeAmount() float64 {
	switch r.Type {
	case TradeBuy:
		return r.buyFeeAmount()
	case TradeSell:
		return r.sellFeeAmount()
	}
	return -1
}

// sell example
// rate: 0.00007900
// amount : 137.43455498
// fee : 0.00001629 BTC (0.15%)
// total in db : 0.01085732
// total in poloniex : 0.01084103 BTC
// 0.00001629 BTC = 0.01085732 BTC * (0.01) * (0.15)
// fee amount = (total in db) * fee
func (r *TradeRow) sellFeeAmount() float64 {
	return r.Total * r.Fee
}

// buy example
// amount : 13.00373802
// fee : 0.03250935 SYS (0.25%)
// total : 0.00094368 BTC
// 0.03250935 SYS = 13.00373802 SYS * (0.01) * (0.25)
// fee amount = amount * fee
func (r *TradeRow) buyFeeAmount() float64 {
	return r.Amount * r.Fee
}

func (r *TradeRow) MyTotal() float64 {
	switch r.Type {
	case TradeBuy:
		return r.Total
	case TradeSell:
		return r.Total - r.FeeAmount()
	}
	return -1
}
func (r *TradeRow) MyAmount() float64 {
	switch r.Type {
	case TradeBuy:
		return r.Amount - r.FeeAmount()
	case TradeSell:
		return r.Amount
	}
	return -1
}

func (r *TradeRow) AmountStr() string {
	return kanaelib.ToFloatStr(r.Amount)
}

func (r *TradeRow) FeeAmountStr() string {
	return kanaelib.ToFloatStr(r.FeeAmount())
}

func (r *TradeRow) MyAmountStr() string {
	return kanaelib.ToFloatStr(r.MyAmount())
}

func (r *TradeRow) RateStr() string {
	return kanaelib.ToFloatStr(r.Rate)
}
