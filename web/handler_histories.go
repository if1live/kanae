package web

import (
	"encoding/json"
	"net/http"

	"github.com/if1live/kanae/histories"
)

func handlerLendingHistories(w http.ResponseWriter, r *http.Request) {
	db, err := histories.NewDatabase(svr.settings.DatabaseFileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer db.Close()

	type Response struct {
		Histories []histories.LendingRow `json:"histories"`
	}
	resp := Response{
		Histories: db.GetLendings("BTC"),
	}

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
