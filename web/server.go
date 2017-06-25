package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/if1live/kanae/histories"
	"github.com/if1live/kanae/kanaelib"
	"github.com/if1live/kanae/reports"
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
	tickers  *histories.TickerCache
}

var svr *Server

func NewServer(addr string, port int, s kanaelib.Settings) *Server {
	if svr != nil {
		panic("already server exists!")
	}

	api := s.MakePoloniex()
	tickers := histories.NewTickerCache(api)
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
		TradeSync   *histories.TradeSync
		LendingSync *histories.LendingSync
		BalanceSync *histories.BalanceSync

		BalanceReport *reports.BalanceReport

		Tickers *histories.TickerCache
	}

	balanceView := svr.db.MakeBalanceView()
	balanceReport := reports.NewBalanceReport("BTC", balanceView.CurrencyRows("BTC"))

	ctx := Context{
		TradeSync:   svr.db.MakeTradeSync(nil),
		LendingSync: svr.db.MakeLendingSync(nil),
		BalanceSync: svr.db.MakeBalanceSync(nil),

		BalanceReport: &balanceReport,
		Tickers:       svr.tickers,
	}
	renderLayoutTemplate(w, "index.html", ctx)
}

func handlerSample(w http.ResponseWriter, r *http.Request) {
	ctx := map[string]string{
		"Title":   "Hello world",
		"Content": "Hi there",
	}
	renderTemplate(w, "sample.html", ctx)
}

func handlerStatic(w http.ResponseWriter, r *http.Request) {
	targetPath := r.URL.Path[len("/static/"):]
	renderStatic(w, r, targetPath)
}

func (s *Server) Run() {
	http.HandleFunc("/", handlerIndex)
	http.HandleFunc("/static/", handlerStatic)
	http.HandleFunc("/sample/", handlerSample)

	http.HandleFunc("/trade/", handlerTradeDispatch)
	http.HandleFunc("/lending/", handlerLending)
	http.HandleFunc("/balance/", handlerBalance)

	http.HandleFunc("/sync/balance", handlerSyncBalance)
	http.HandleFunc("/sync/trade", handlerSyncTrade)
	http.HandleFunc("/sync/lending", handlerSyncLending)
	http.HandleFunc("/sync/ticker", handlerSyncTicker)

	http.HandleFunc("/histories/lending/", handlerLendingHistories)
	http.HandleFunc("/histories/trade/", handlerTradeHistories)

	addr := s.addr + ":" + strconv.Itoa(s.port)
	fmt.Println("run server on", addr)
	http.ListenAndServe(addr, nil)
}

func (s *Server) Close() {
	s.db.Close()
}
