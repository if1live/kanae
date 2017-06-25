package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/if1live/kanae/balances"
	"github.com/if1live/kanae/exchanges"
	"github.com/if1live/kanae/histories"
	"github.com/if1live/kanae/kanaelib"
	"github.com/if1live/kanae/lendings"
	"github.com/thrasher-/gocryptotrader/exchanges/poloniex"
)

/*
reference
http://www.alexedwards.net/blog/golang-response-snippets
*/

type Server struct {
	addr string
	port int

	settings kanaelib.Settings
	db       *histories.Database
	api      *poloniex.Poloniex
	tickers  *kanaelib.TickerCache
}

var svr *Server

func NewServer(addr string, port int, s kanaelib.Settings) *Server {
	if svr != nil {
		panic("already server exists!")
	}

	api := s.MakePoloniex()
	tickers := kanaelib.NewTickerCache(api)
	tickers.Refresh()

	// share single orm
	db, err := histories.NewDatabase(s.DatabaseFileName)
	if err != nil {
		panic(err)
	}

	svr = &Server{
		addr:     addr,
		port:     port,
		settings: s,
		db:       &db,
		api:      api,
		tickers:  tickers,
	}

	return svr
}

func handlerIndex(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "hello world : %s", r.URL.Path[1:])

	type Context struct {
		ExchangeSync *exchanges.Sync
		LendingSync  *lendings.Sync
		BalanceSync  *balances.Sync

		BalanceReport *balances.Report

		Tickers *kanaelib.TickerCache
	}

	balanceView := svr.db.MakeBalanceView()
	balanceReport := balances.NewReport("BTC", balanceView.CurrencyRows("BTC"))

	ctx := Context{
		ExchangeSync: svr.db.MakeExchangeSync(nil),
		LendingSync:  svr.db.MakeLendingSync(nil),
		BalanceSync:  svr.db.MakeBalanceSync(nil),

		BalanceReport: &balanceReport,
		Tickers:       svr.tickers,
	}
	err := renderLayoutTemplate(w, "layout.html", "index.html", ctx)
	if err != nil {
		renderErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
}

func handlerStatic(w http.ResponseWriter, r *http.Request) {
	targetPath := r.URL.Path[len("/static/"):]
	renderStatic(w, r, targetPath)
}

func (s *Server) Run() {
	http.HandleFunc("/", handlerIndex)
	http.HandleFunc("/static/", handlerStatic)

	http.HandleFunc("/exchange/", handlerExchangeDispatch)
	http.HandleFunc("/lending/", handlerLending)
	http.HandleFunc("/balance/", handlerBalance)

	http.HandleFunc("/sync/balance", handlerSyncBalance)
	http.HandleFunc("/sync/exchange", handlerSyncExchange)
	http.HandleFunc("/sync/lending", handlerSyncLending)
	http.HandleFunc("/sync/ticker", handlerSyncTicker)

	//http.HandleFunc("/histories/lending/", handlerLendingHistories)
	//http.HandleFunc("/histories/trade/", handlerTradeHistories)

	addr := s.addr + ":" + strconv.Itoa(s.port)
	fmt.Println("run server on", addr)
	http.ListenAndServe(addr, nil)
}

func (s *Server) Close() {
	s.db.Close()
}
