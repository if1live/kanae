package reports

import (
	"github.com/if1live/kanae/histories"
	"github.com/if1live/kanae/kanaelib"
)

type TradeReport struct {
	Asset    string
	Currency string
	Rows     []histories.TradeRow
}

func NewTradeReport(asset, currency string, rows []histories.TradeRow) TradeReport {
	return TradeReport{
		Asset:    asset,
		Currency: currency,
		Rows:     rows,
	}
}

func (r *TradeReport) CurrentAsset() float64 {
	var total float64
	total += r.TotalAssetBuys()
	total -= r.TotalAssetSells()
	return total
}
func (r *TradeReport) ProfitLoss() float64 {
	var total float64
	for _, r := range r.Rows {
		switch r.Type {
		case histories.TradeBuy:
			total -= r.Total
		case histories.TradeSell:
			total += r.Total
		}
	}
	return total
}

func (r *TradeReport) TotalAssetBuys() float64 {
	var total float64
	for _, r := range r.Rows {
		if r.Type == histories.TradeBuy {
			total += r.MyAmount()
		}
	}
	return total
}
func (r *TradeReport) TotalAssetSells() float64 {
	var total float64
	for _, r := range r.Rows {
		if r.Type == histories.TradeSell {
			total += r.Amount
		}
	}
	return total
}

func (r *TradeReport) CurrentFixedAsset() float64 {
	return kanaelib.ToPoloniexFixed(r.CurrentAsset())
}

func (r *TradeReport) FixedProfitLoss() float64 {
	return kanaelib.ToPoloniexFixed(r.ProfitLoss())
}
