package web

import (
	"net/http"
	"strings"
	"time"

	"github.com/if1live/kanae/histories"
	"github.com/if1live/kanae/reports"
	"github.com/thrasher-/gocryptotrader/exchanges/poloniex"
)

func handlerTradeDispatch(w http.ResponseWriter, r *http.Request) {
	asset := strings.ToUpper(r.URL.Path[len("/trade/"):])
	if asset == "" {
		handlerTradeIndex(w, r)
	} else {
		handlerTradeAsset(w, r, asset)
	}
}

func handlerTradeIndex(w http.ResponseWriter, r *http.Request) {
	type Context struct {
		Sync *histories.TradeSync
		View *histories.TradeView

		UsedAssets []string
	}

	sync := svr.db.MakeTradeSync(nil)
	view := svr.db.MakeTradeView()
	usedAssets := view.UsedAssets("BTC")

	ctx := Context{
		Sync:       sync,
		View:       view,
		UsedAssets: usedAssets,
	}
	renderLayoutTemplate(w, "trade.html", ctx)
}

func handlerTradeAsset(w http.ResponseWriter, r *http.Request, asset string) {
	type Context struct {
		Asset string

		Sync   *histories.TradeSync
		View   *histories.TradeView
		Report *reports.TradeReport

		Ticker          poloniex.PoloniexTicker
		TickerUpdatedAt time.Time
	}

	ticker, err := svr.tickers.Get(asset, "BTC")
	if err != nil {
		renderErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	sync := svr.db.MakeTradeSync(nil)
	view := svr.db.MakeTradeView()
	histories := view.All(asset, "BTC")
	report := reports.NewTradeReport(asset, "BTC", histories)
	ctx := Context{
		Asset:           asset,
		Sync:            sync,
		View:            view,
		Report:          &report,
		Ticker:          ticker,
		TickerUpdatedAt: svr.tickers.UpdatedAt,
	}
	renderLayoutTemplate(w, "trade_detail.html", ctx)
}
