package histories

import (
	"errors"
	"net/url"

	"github.com/thrasher-/gocryptotrader/exchanges/poloniex"
)

// some apis are missing in github.com/thrasher-/gocryptotrader
// * returnLendingHistory

type PoloniexLendingHistory struct {
	ID       int64   `json:"id"`
	Currency string  `json:"currency"`
	Rate     float64 `json:"rate,string"`
	Amount   float64 `json:"amount,string"`
	Duration float64 `json:"duration,string"`
	Interest float64 `json:"interest,string"`
	Fee      float64 `json:"fee,string"`
	Earned   float64 `json:"earned,string"`
	Open     string  `json:"open"`
	Close    string  `json:"close"`
}

const (
	POLONIEX_LENDING_HISTORY = "returnLendingHistory"
)

func GetLendingHistory(p *poloniex.Poloniex, start, end string) ([]PoloniexLendingHistory, error) {
	vals := url.Values{}

	if start != "" {
		vals.Set("start", start)
	}

	if end != "" {
		vals.Set("end", end)
	}

	type Response struct {
		Data []PoloniexLendingHistory
	}
	result := Response{}

	err := p.SendAuthenticatedHTTPRequest("POST", POLONIEX_LENDING_HISTORY, vals, &result.Data)

	if err != nil {
		return nil, err
	}

	if result.Data == nil {
		return nil, errors.New("There are no lending history")
	}

	return result.Data, nil
}
