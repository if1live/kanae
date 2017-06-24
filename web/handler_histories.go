package web

import (
	"net/http"

	"strings"

	"github.com/if1live/kanae/histories"
	"github.com/thrasher-/gocryptotrader/exchanges/poloniex"
)

func handlerLendingHistories(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Histories []histories.PoloniexLendingHistory `json:"histories"`
	}

	q := db.MakeLendingQuery()
	rows := q.GetAll("BTC")
	histories := []histories.PoloniexLendingHistory{}
	for _, row := range rows {
		histories = append(histories, row.MakeHistory())
	}
	resp := Response{
		Histories: histories,
	}
	renderJSON(w, resp)
}

func handlerTradeHistories(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Asset     string                                      `json:"asset"`
		Histories []poloniex.PoloniexAuthentictedTradeHistory `json:"histories"`
	}

	asset := strings.ToUpper(r.URL.Path[len("/histories/trade/"):])
	q := db.MakeTradeQuery()
	rows := q.GetAll(asset, "BTC")
	histories := []poloniex.PoloniexAuthentictedTradeHistory{}
	for _, row := range rows {
		histories = append(histories, row.MakeHistory())
	}
	resp := Response{
		Histories: histories,
		Asset:     asset,
	}
	renderJSON(w, resp)
}
