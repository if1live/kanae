package web

import (
	"net/http"

	"github.com/if1live/kanae/histories"
	"github.com/if1live/kanae/reports"
)

func handlerBalance(w http.ResponseWriter, r *http.Request) {
	type Context struct {
		Sync   *histories.BalanceSync
		View   *histories.BalanceView
		Report *reports.BalanceReport
	}

	sync := db.MakeBalanceSync(nil)
	view := db.MakeBalanceView()
	currency := "BTC"
	report := reports.NewBalanceReport(currency, view.CurrencyRows(currency))

	ctx := Context{
		Sync:   sync,
		View:   view,
		Report: &report,
	}
	renderLayoutTemplate(w, "balance.html", ctx)
}
