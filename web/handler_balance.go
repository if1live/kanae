package web

import (
	"net/http"

	"github.com/if1live/kanae/balances"
)

func handlerBalance(w http.ResponseWriter, r *http.Request) {
	type Context struct {
		Sync   *balances.Sync
		View   *balances.View
		Report *balances.Report
	}

	sync := svr.db.MakeBalanceSync(nil)
	view := svr.db.MakeBalanceView()
	currency := "BTC"
	report := balances.NewReport(currency, view.CurrencyRows(currency))

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
