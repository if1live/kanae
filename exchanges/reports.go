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

func reverseExchanges(a []Exchange) {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
}

func NewReport(asset, currency string, ticker poloniex.PoloniexTicker, rows []Exchange) (closed, opened *Report) {
	// input data : date desc
	ascRows := append([]Exchange(nil), rows...)
	reverseExchanges(ascRows)

	var myAmountAccum float64
	lastZeroSumIdx := -1
	for i, r := range ascRows {
		if r.Type == ExchangeBuy {
			myAmountAccum += r.MyAmount()
		} else {
			myAmountAccum -= r.Amount
		}
		if myAmountAccum <= float64(0) {
			lastZeroSumIdx = i
		}
	}

	if lastZeroSumIdx == -1 {
		// closed exchange not exists
		closed = nil
		opened = &Report{
			Asset:    asset,
			Currency: currency,
			Rows:     rows,
			ticker:   ticker,
		}
		return
	}

	closedRows := ascRows[:lastZeroSumIdx+1]
	reverseExchanges(closedRows)

	openedRows := ascRows[lastZeroSumIdx+1:]
	reverseExchanges(openedRows)

	if len(closedRows) == 0 {
		closed = nil
	} else {
		closed = &Report{
			Asset:    asset,
			Currency: currency,
			Rows:     closedRows,
			ticker:   ticker,
		}
	}

	if len(openedRows) == 0 {
		opened = nil
	} else {
		opened = &Report{
			Asset:    asset,
			Currency: currency,
			Rows:     openedRows,
			ticker:   ticker,
		}
	}

	return
}

func NewReports(tickers *kanaelib.TickerCache, rows []Exchange) (closedList, openedList []*Report) {
	// find all asset-currency pairs
	currencyPairSet := mapset.NewSet()
	for _, e := range rows {
		currencyPairSet.Add(e.CurrencyPair())
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
	closedList = []*Report{}
	openedList = []*Report{}

	for _, currencyPair := range currencyPairs {
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
		closed, opened := NewReport(asset, currency, ticker, selected)
		if closed != nil {
			closedList = append(closedList, closed)
		}
		if opened != nil {
			openedList = append(openedList, opened)
		}
	}
	return
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
