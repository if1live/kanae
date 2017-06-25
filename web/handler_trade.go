package web

import (
	"net/http"

	"github.com/if1live/kanae/histories"
	"github.com/if1live/kanae/reports"
)

func handlerTradeIndex(w http.ResponseWriter, r *http.Request) {
	type Context struct {
		Sync *histories.TradeSync
		View *histories.TradeView

		UsedAssets []string
		Reports    []reports.TradeReport
	}

	sync := db.MakeTradeSync(nil)
	view := db.MakeTradeView()

	rs := []reports.TradeReport{}
	usedAssets := view.UsedAssets("BTC")
	for _, asset := range usedAssets {
		histories := view.All(asset, "BTC")
		report := reports.NewTradeReport(asset, "BTC", histories)
		rs = append(rs, report)
	}

	ctx := Context{
		Sync:       sync,
		View:       view,
		UsedAssets: usedAssets,
		Reports:    rs,
	}
	renderLayoutTemplate(w, "trade.html", ctx)
}
