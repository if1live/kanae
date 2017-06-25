package reports

import (
	"github.com/if1live/kanae/histories"
	"github.com/if1live/kanae/kanaelib"
	"github.com/thrasher-/gocryptotrader/exchanges/poloniex"
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
	total += r.TotalBuyAsset()
	total -= r.TotalSellAsset()
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

func (r *TradeReport) TotalBuyAsset() float64 {
	var total float64
	for _, r := range r.Rows {
		if r.Type == histories.TradeBuy {
			total += r.MyAmount()
		}
	}
	return total
}
func (r *TradeReport) TotalSellAsset() float64 {
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

type CurrencyState struct {
	FirstCurrency  string  `json:"firstCurrency"`
	SecondCurrency string  `json:"secondCurrency"`
	Available      float64 `json:"available,string"`

	Ticker poloniex.PoloniexTicker
}

/*
func (r *PoloniexReport) GetStates() []CurrencyState {
	balances, err := r.api.GetBalances()
	if err != nil {
		check(err)
	}

	tickers, err := r.api.GetTicker()
	if err != nil {
		check(err)
	}

	states := []CurrencyState{}

	for key, balance := range balances.Currency {
		if balance == 0 {
			continue
		}
		if key == "BTC" {
			continue
		}

		// example: BTC_XRP
		tickerKey := "BTC_" + key
		ticker := tickers[tickerKey]

		s := CurrencyState{
			FirstCurrency:  "BTC",
			SecondCurrency: key,
			Available:      balance,

			Ticker: ticker,
		}
		states = append(states, s)
	}
	return states
}
*/
