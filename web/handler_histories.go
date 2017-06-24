package web

import (
	"net/http"

	"strings"

	"github.com/if1live/kanae/histories"
	"github.com/thrasher-/gocryptotrader/exchanges/poloniex"
)

func handlerLendingHistories(w http.ResponseWriter, r *http.Request) {
	db, err := histories.NewDatabase(svr.settings.DatabaseFileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	type Response struct {
		Histories []histories.PoloniexLendingHistory `json:"histories"`
	}

	rows := db.GetLendings("BTC")
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
	db, err := histories.NewDatabase(svr.settings.DatabaseFileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	type Response struct {
		Asset     string                                      `json:"asset"`
		Histories []poloniex.PoloniexAuthentictedTradeHistory `json:"histories"`
	}

	asset := strings.ToUpper(r.URL.Path[len("/histories/trade/"):])
	rows := db.GetAllTrades(asset, "BTC")
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
