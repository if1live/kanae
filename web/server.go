package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/if1live/kanae/histories"
	"github.com/if1live/kanae/kanaelib"
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
}

// use singleton
// only one web server exist
var svr *Server
var db *histories.Database
var api *poloniex.Poloniex

func NewServer(addr string, port int, s kanaelib.Settings) *Server {
	if svr != nil {
		panic("already server exists!")
	}

	// share single orm
	dbobj, err := histories.NewDatabase(s.DatabaseFileName)
	if err != nil {
		panic(err)
	}
	db = &dbobj

	api = s.MakePoloniex()

	svr = &Server{
		addr:     addr,
		port:     port,
		settings: s,
	}
	return svr
}

func handlerIndex(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "hello world : %s", r.URL.Path[1:])

	type Context struct {
		TradeSync   *histories.TradeSync
		LendingSync *histories.LendingSync
		BalanceSync *histories.BalanceSync
	}

	ctx := Context{
		TradeSync:   db.MakeTradeSync(nil),
		LendingSync: db.MakeLendingSync(nil),
		BalanceSync: db.MakeBalanceSync(nil),
	}

	renderLayoutTemplate(w, "index.html", ctx)
}

func handlerStatic(w http.ResponseWriter, r *http.Request) {
	targetPath := r.URL.Path[len("/static/"):]
	renderStatic(w, r, targetPath)
}

func (s *Server) Run() {
	http.HandleFunc("/", handlerIndex)
	http.HandleFunc("/static/", handlerStatic)
	http.HandleFunc("/trade/", handlerTradeIndex)
	http.HandleFunc("/lending/", handlerLending)

	http.HandleFunc("/sync/balance", handlerSyncBalance)
	http.HandleFunc("/sync/trade", handlerSyncTrade)
	http.HandleFunc("/sync/lending", handlerSyncLending)

	http.HandleFunc("/histories/lending/", handlerLendingHistories)
	http.HandleFunc("/histories/trade/", handlerTradeHistories)

	addr := s.addr + ":" + strconv.Itoa(s.port)
	fmt.Println("run server on", addr)
	http.ListenAndServe(addr, nil)
}

func (s *Server) Close() {
	db.Close()
}
