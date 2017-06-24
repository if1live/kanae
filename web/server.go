package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/if1live/kanae/kanaelib"
)

type Server struct {
	addr string
	port int

	settings kanaelib.Settings
}

// use singleton
// only one web server exist
var svr *Server

func NewServer(addr string, port int, s kanaelib.Settings) *Server {
	if svr != nil {
		panic("already server exists!")
	}

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
		Name string
	}
	ctx := Context{
		Name: "todo",
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
