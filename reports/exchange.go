package reports

import (
	"sort"

	"strings"

	"github.com/deckarep/golang-set"
	"github.com/if1live/kanae/histories"
	"github.com/if1live/kanae/histories/exchanges"
	"github.com/thrasher-/gocryptotrader/exchanges/poloniex"
)

type ExchangeReport struct {
	Asset    string
	Currency string
	Rows     []exchanges.Exchange
	ticker   poloniex.PoloniexTicker
}

func NewExchangeReport(asset, currency string, ticker poloniex.PoloniexTicker, rows []exchanges.Exchange) *ExchangeReport {
	return &ExchangeReport{
		Asset:    asset,
		Currency: currency,
		Rows:     rows,
		ticker:   ticker,
	}
}

func NewExchangeReports(tickers *histories.TickerCache, rows []exchanges.Exchange) []*ExchangeReport {
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
	reports := make([]*ExchangeReport, len(currencyPairs))
	for i, currencyPair := range currencyPairs {
		tokens := strings.Split(currencyPair, "_")
		asset := tokens[0]
		currency := tokens[1]

		selected := []exchanges.Exchange{}
		for _, r := range rows {
			if r.Asset == asset && r.Currency == currency {
				selected = append(selected, r)
			}
		}
		ticker, _ := tickers.Get(asset, currency)
		reports[i] = NewExchangeReport(asset, currency, ticker, selected)
	}
	return reports
}

func (r *ExchangeReport) CurrentAsset() float64 {
	var total float64
	total += r.TotalAssetBuys()
	total -= r.TotalAssetSells()
	if total < 0 {
		total = 0
	}
	return total
}

func (r *ExchangeReport) EquivalentCurrency() float64 {
	asset := r.CurrentAsset()
	rate := r.ticker.Last
	return asset * rate
}

func (r *ExchangeReport) ProfitLoss() float64 {
	var total float64
	for _, r := range r.Rows {
		switch r.Type {
		case exchanges.ExchangeBuy:
			total -= r.MyTotal()
		case exchanges.ExchangeSell:
			total += r.MyTotal()
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

func (r *ExchangeReport) TotalCurrencyBuys() float64 {
	var total float64
	for _, r := range r.Rows {
		if r.Type == exchanges.ExchangeBuy {
			total += r.MyTotal()
		}
	}
	return total
}

func (r *ExchangeReport) EarningRate() float64 {
	return r.ProfitLoss() / r.TotalCurrencyBuys() * 100
}
