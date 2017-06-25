package web

/*
func handlerLendingHistories(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Histories []lendings.PoloniexLendingHistory `json:"histories"`
	}

	q := svr.db.MakeLendingView()
	rows := q.All("BTC")
	histories := []lendings.PoloniexLendingHistory{}
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
	q := svr.db.MakeExchangeView()
	rows := q.Get(asset, "BTC")
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
*/
