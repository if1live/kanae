package web

import (
	"net/http"
	"time"

	"github.com/if1live/kanae/histories"
)

func handlerSyncPage(w http.ResponseWriter, r *http.Request) {
	type Context struct {
		Trade   *histories.TradeSync
		Lending *histories.LendingSync
		Balance *histories.BalanceSync
	}

	ctx := Context{
		Trade:   db.MakeTradeSync(nil),
		Lending: db.MakeLendingSync(nil),
		Balance: db.MakeBalanceSync(nil),
	}

	renderLayoutTemplate(w, "sync.html", ctx)
}

func handlerSync(w http.ResponseWriter, r *http.Request, sync histories.Synchronizer) {
	rowcount, err := sync.SyncRecent()
	if err != nil {
		renderErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	type Response struct {
		RowCount  int       `json:"row_count"`
		UpdatedAt time.Time `json:"updated_at"`
		LastTime  time.Time `json:"last_time"`
	}
	resp := Response{
		RowCount:  rowcount,
		UpdatedAt: time.Now(),
		LastTime:  sync.GetLastTime(),
	}
	renderJSON(w, resp)
}

func handlerSyncBalance(w http.ResponseWriter, r *http.Request) {
	if ok := checkPostRequest(w, r); !ok {
		return
	}

	sync := db.MakeBalanceSync(api)
	handlerSync(w, r, sync)
}

func handlerSyncTrade(w http.ResponseWriter, r *http.Request) {
	if ok := checkPostRequest(w, r); !ok {
		return
	}

	sync := db.MakeTradeSync(api)
	handlerSync(w, r, sync)
}

func handlerSyncLending(w http.ResponseWriter, r *http.Request) {
	if ok := checkPostRequest(w, r); !ok {
		return
	}

	sync := db.MakeLendingSync(api)
	handlerSync(w, r, sync)
}
