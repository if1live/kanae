package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/if1live/kanae/exchanges"
	"github.com/if1live/kanae/kanaelib"
)

type ExchangeRecord struct {
	CurrencyPair string `json:"currencyPair"`
	Date         string `json:"date"`
	Type         string `json:"type"`
	Category     string `json:"category"`
	Rate         string `json:"rate"`
	Amount       string `json:"amount"`
	Total        string `json:"total"`
	Fee          string `json:"fee"`
}

func NewExchangeRecord(r *exchanges.Exchange) ExchangeRecord {
	feeAmountStr := kanaelib.ToFloatStr(r.FeeAmount())
	feeUnit := ""
	if r.Type == exchanges.ExchangeBuy {
		feeUnit = r.Asset
	} else if r.Type == exchanges.ExchangeSell {
		feeUnit = r.Currency
	}
	feePercent := r.Fee * 100
	fee := fmt.Sprintf("%s %s (%.2f %%)", feeAmountStr, feeUnit, feePercent)

	return ExchangeRecord{
		CurrencyPair: r.Currency + "_" + r.Asset,
		Type:         r.Type,
		Category:     r.Category,
		Amount:       kanaelib.ToFloatStr(r.Amount),
		Rate:         kanaelib.ToFloatStr(r.Rate),
		Date:         r.Date.Format(time.RFC3339),
		Total:        kanaelib.ToFloatStr(r.MyTotal()) + r.Currency,
		Fee:          fee,
	}
}

func handlerPoloniexStaticJS(w http.ResponseWriter, r *http.Request) {
	targetPath := r.URL.Path[len("/js/"):]
	targetPath = path.Join("js", targetPath)
	renderStatic(w, r, targetPath)
}
func handlerPoloniexStaticCSS(w http.ResponseWriter, r *http.Request) {
	targetPath := r.URL.Path[len("/css/"):]
	targetPath = path.Join("css", targetPath)
	renderStatic(w, r, targetPath)
}

func handlerPoloniexIndex(w http.ResponseWriter, r *http.Request) {
	renderStaticHtml(w, r, "trade_history.html")
}

func handlerPoloniexPrivateAPI(w http.ResponseWriter, r *http.Request) {
	cmd := r.FormValue("command")

	if cmd == "returnPaginatedTradeHistory" {
		start := r.FormValue("start")
		end := r.FormValue("end")
		//page := r.FormValue("page")
		//tradesPerPage := r.FormValue("tradesPerPage")
		//typeStr := r.FormValue("type")

		startTimestamp, _ := strconv.ParseInt(start, 10, 64)
		endTimestamp, _ := strconv.ParseInt(end, 10, 64)

		startTime := time.Unix(startTimestamp, 0)
		endTime := time.Unix(endTimestamp, 0)

		view := svr.db.MakeExchangeView()
		rows := view.AllWithRange(startTime, endTime)

		buf := new(bytes.Buffer)
		buf.Write([]byte("[1,"))
		for i, r := range rows {
			record := NewExchangeRecord(&r)
			data, _ := json.Marshal(record)
			buf.Write(data)

			if i < len(rows)-1 {
				buf.Write([]byte(","))
			}
		}
		buf.Write([]byte("]"))

		w.Write(buf.Bytes())
		return
	}
	// else..
	renderStaticHtml(w, r, "test.json")
}
