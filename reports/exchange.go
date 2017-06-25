package reports

import (
	"github.com/if1live/kanae/histories/exchanges"
	"github.com/if1live/kanae/kanaelib"
)

type ExchangeReport struct {
	Asset    string
	Currency string
	Rows     []exchanges.Exchange
}

func NewExchangeReport(asset, currency string, rows []exchanges.Exchange) ExchangeReport {
	return ExchangeReport{
		Asset:    asset,
		Currency: currency,
		Rows:     rows,
	}
}

func (r *ExchangeReport) CurrentAsset() float64 {
	var total float64
	total += r.TotalAssetBuys()
	total -= r.TotalAssetSells()
	return total
}
func (r *ExchangeReport) ProfitLoss() float64 {
	var total float64
	for _, r := range r.Rows {
		switch r.Type {
		case exchanges.ExchangeBuy:
			total -= r.Total
		case exchanges.ExchangeSell:
			total += r.Total
		}
	}
	return total
}

func (r *ExchangeReport) TotalAssetBuys() float64 {
	var total float64
	for _, r := range r.Rows {
		if r.Type == exchanges.ExchangeBuy {
			total += r.MyAmount()
		}
	}
	return total
}
func (r *ExchangeReport) TotalAssetSells() float64 {
	var total float64
	for _, r := range r.Rows {
		if r.Type == exchanges.ExchangeSell {
			total += r.Amount
		}
	}
	return total
}

func (r *ExchangeReport) CurrentFixedAsset() float64 {
	return kanaelib.ToPoloniexFixed(r.CurrentAsset())
}

func (r *ExchangeReport) FixedProfitLoss() float64 {
	return kanaelib.ToPoloniexFixed(r.ProfitLoss())
}
