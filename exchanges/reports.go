package exchanges

import (
	"sort"

	"strings"

	"github.com/deckarep/golang-set"
	"github.com/if1live/kanae/kanaelib"
	"github.com/thrasher-/gocryptotrader/exchanges/poloniex"
)

type Report struct {
	Asset    string
	Currency string
	Rows     []Exchange
	ticker   poloniex.PoloniexTicker
}

func NewReport(asset, currency string, ticker poloniex.PoloniexTicker, rows []Exchange) *Report {
	return &Report{
		Asset:    asset,
		Currency: currency,
		Rows:     rows,
		ticker:   ticker,
	}
}

func NewReports(tickers *kanaelib.TickerCache, rows []Exchange) []*Report {
	// find all asset-currency pairs
	currencyPairSet := mapset.NewSet()
	for _, e := range rows {
		currencyPair := e.Asset + "_" + e.Currency
		currencyPairSet.Add(currencyPair)
	}

	// sort asset-currency pairs
	currencyPairs := []string{}
	it := currencyPairSet.Iterator()
	for elem := range it.C {
		s := elem.(string)
		currencyPairs = append(currencyPairs, s)
	}
	sort.Strings(currencyPairs)

	// generate reports
	reports := make([]*Report, len(currencyPairs))
	for i, currencyPair := range currencyPairs {
		tokens := strings.Split(currencyPair, "_")
		asset := tokens[0]
		currency := tokens[1]

		selected := []Exchange{}
		for _, r := range rows {
			if r.Asset == asset && r.Currency == currency {
				selected = append(selected, r)
			}
		}
		ticker, _ := tickers.Get(asset, currency)
		reports[i] = NewReport(asset, currency, ticker, selected)
	}
	return reports
}

func (r *Report) CurrentAsset() float64 {
	var total float64
	total += r.TotalAssetBuys()
	total -= r.TotalAssetSells()
	if total < 0 {
		total = 0
	}
	return total
}

func (r *Report) EquivalentCurrency() float64 {
	asset := r.CurrentAsset()
	rate := r.ticker.Last
	return asset * rate
}

func (r *Report) ProfitLoss() float64 {
	var total float64
	for _, r := range r.Rows {
		switch r.Type {
		case ExchangeBuy:
			total -= r.MyTotal()
		case ExchangeSell:
			total += r.MyTotal()
		}
	}
	return total
}

func (r *Report) TotalAssetBuys() float64 {
	var total float64
	for _, r := range r.Rows {
		if r.Type == ExchangeBuy {
			total += r.MyAmount()
		}
	}
	return total
}
func (r *Report) TotalAssetSells() float64 {
	var total float64
	for _, r := range r.Rows {
		if r.Type == ExchangeSell {
			total += r.Amount
		}
	}
	return total
}

func (r *Report) TotalCurrencyBuys() float64 {
	var total float64
	for _, r := range r.Rows {
		if r.Type == ExchangeBuy {
			total += r.MyTotal()
		}
	}
	return total
}

func (r *Report) EarningRate() float64 {
	return r.ProfitLoss() / r.TotalCurrencyBuys() * 100
}

type ClosedSummaryReport struct {
	reports []*Report
}

func NewClosedSummaryReport(rs []*Report) *ClosedSummaryReport {
	return &ClosedSummaryReport{
		reports: rs,
	}
}
func (r *ClosedSummaryReport) ProfitLoss() float64 {
	var total float64
	for _, r := range r.reports {
		total += r.ProfitLoss()
	}
	return total
}
