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

	sync := svr.db.MakeBalanceSync(nil)
	view := svr.db.MakeBalanceView()
	currency := "BTC"
	report := reports.NewBalanceReport(currency, view.CurrencyRows(currency))

	ctx := Context{
		Sync:   sync,
		View:   view,
		Report: &report,
	}
	err := renderLayoutTemplate(w, "layout.html", "balance.html", ctx)
	if err != nil {
		renderErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
}
