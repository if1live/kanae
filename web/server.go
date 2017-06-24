package web

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/if1live/kanae/histories"
	"github.com/if1live/kanae/kanaelib"
	"github.com/if1live/kanae/reports"
)

type Server struct {
	addr string
	port int

	settings kanaelib.Settings
}

// use singleton
// only one web server exist
var svr *Server
var db *histories.Database

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
		LastTradeSyncTime   time.Time
		LastLendingSyncTime time.Time
		LastBalanceSyncTime time.Time

		LendingHistories []histories.LendingRow

		TradeUsedAssets []string
		TradeReports    []reports.TradeReport
	}

	tradeSync := db.MakeTradeSync(nil)
	lendingSync := db.MakeLendingSync(nil)
	balanceSync := db.MakeBalanceSync(nil)

	tradeView := db.MakeTradeView()

	tradeReports := []reports.TradeReport{}
	usedAssets := tradeView.UsedAssets("BTC")
	for _, asset := range usedAssets {
		histories := tradeView.All(asset, "BTC")
		report := reports.NewTradeReport(asset, "BTC", histories)
		tradeReports = append(tradeReports, report)
	}

	ctx := Context{
		LastTradeSyncTime:   tradeSync.GetLastTime(),
		LastLendingSyncTime: lendingSync.GetLastTime(),
		LastBalanceSyncTime: balanceSync.GetLastTime(),

		TradeReports:    tradeReports,
		TradeUsedAssets: usedAssets,
	}
	renderTemplate(w, "index.html", ctx)
}

func handlerStatic(w http.ResponseWriter, r *http.Request) {
	targetPath := r.URL.Path[len("/static/"):]
	renderStatic(w, r, targetPath)
}

func (s *Server) Run() {
	http.HandleFunc("/", handlerIndex)
	http.HandleFunc("/static/", handlerStatic)

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
