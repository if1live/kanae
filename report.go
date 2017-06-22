package main

import (
	"fmt"

	"github.com/thrasher-/gocryptotrader/config"
	"github.com/thrasher-/gocryptotrader/exchanges/poloniex"
)

type PoloniexReport struct {
	api *poloniex.Poloniex
}

type CurrencyState struct {
	FirstCurrency  string  `json:"firstCurrency"`
	SecondCurrency string  `json:"secondCurrency"`
	Available      float64 `json:"available,string"`

	Ticker poloniex.PoloniexTicker
}

func NewPoloniexReport(s Settings) *PoloniexReport {
	conf := config.ExchangeConfig{
		Enabled:                 true,
		APIKey:                  s.APIKey,
		APISecret:               s.APISecret,
		AuthenticatedAPISupport: true,
		Verbose:                 true,
	}
	p := poloniex.Poloniex{}
	p.Setup(conf)

	return &PoloniexReport{
		api: &p,
	}
}

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

func (r *PoloniexReport) Generate() {
	states := r.GetStates()
	fmt.Println(states)
}
