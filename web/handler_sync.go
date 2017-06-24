package web

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/if1live/kanae/histories"
)

func handlerSync(w http.ResponseWriter, r *http.Request, sync histories.Synchronizer) {
	rowcount, err := sync.SyncRecent()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func handlerSyncBalance(w http.ResponseWriter, r *http.Request) {
	db, err := histories.NewDatabase(svr.settings.DatabaseFileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer db.Close()

	api := svr.settings.MakePoloniex()
	sync := db.MakeBalanceSync(api)
	handlerSync(w, r, sync)
}

func handlerSyncTrade(w http.ResponseWriter, r *http.Request) {
	db, err := histories.NewDatabase(svr.settings.DatabaseFileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer db.Close()

	api := svr.settings.MakePoloniex()
	sync := db.MakeTradeSync(api)
	handlerSync(w, r, sync)
}

func handlerSyncLending(w http.ResponseWriter, r *http.Request) {
	db, err := histories.NewDatabase(svr.settings.DatabaseFileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer db.Close()

	api := svr.settings.MakePoloniex()
	sync := db.MakeLendingSync(api)
	handlerSync(w, r, sync)
}
