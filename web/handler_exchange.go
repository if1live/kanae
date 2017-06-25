package web

import (
	"net/http"
	"strings"
	"time"

	"github.com/if1live/kanae/histories/exchanges"
	"github.com/if1live/kanae/reports"
	"github.com/thrasher-/gocryptotrader/exchanges/poloniex"
)

func handlerExchangeDispatch(w http.ResponseWriter, r *http.Request) {
	asset := strings.ToUpper(r.URL.Path[len("/exchange/"):])
	if asset == "" {
		handlerExchangeIndex(w, r)
	} else {
		handlerExchangeAsset(w, r, asset)
	}
}

func handlerExchangeIndex(w http.ResponseWriter, r *http.Request) {
	type Context struct {
		Sync *exchanges.Sync
		View *exchanges.View

		OpenedReports []*reports.ExchangeReport
		ClosedReports []*reports.ExchangeReport
	}

	sync := svr.db.MakeExchangeSync(nil)
	view := svr.db.MakeExchangeView()
	rs := reports.NewExchangeReports(svr.tickers, view.All())

	opens := []*reports.ExchangeReport{}
	closes := []*reports.ExchangeReport{}
	for _, r := range rs {
		if r.CurrentAsset() == 0 {
			closes = append(closes, r)
		} else {
			opens = append(opens, r)
		}
	}

	ctx := Context{
		Sync:          sync,
		View:          view,
		OpenedReports: opens,
		ClosedReports: closes,
	}
	err := renderLayoutTemplate(w, "layout.html", "exchange_index.html", ctx)
	if err != nil {
		renderErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
}

func handlerExchangeAsset(w http.ResponseWriter, r *http.Request, asset string) {
	type Context struct {
		Asset string

		Sync   *exchanges.Sync
		View   *exchanges.View
		Report *reports.ExchangeReport

		Ticker          poloniex.PoloniexTicker
		TickerUpdatedAt time.Time
	}

	ticker, err := svr.tickers.Get(asset, "BTC")
	if err != nil {
		renderErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	sync := svr.db.MakeExchangeSync(nil)
	view := svr.db.MakeExchangeView()
	histories := view.Get(asset, "BTC")
	report := reports.NewExchangeReport(asset, "BTC", ticker, histories)
	ctx := Context{
		Asset:           asset,
		Sync:            sync,
		View:            view,
		Report:          report,
		Ticker:          ticker,
		TickerUpdatedAt: svr.tickers.UpdatedAt,
	}
	err = renderLayoutTemplate(w, "layout.html", "exchange_detail.html", ctx)
	if err != nil {
		renderErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
}
