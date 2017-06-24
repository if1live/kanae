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

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world : %s", r.URL.Path[1:])
}

func (s *Server) Run() {
	http.HandleFunc("/", handler)

	http.HandleFunc("/sync/balance", handlerSyncBalance)
	http.HandleFunc("/sync/trade", handlerSyncTrade)
	http.HandleFunc("/sync/lending", handlerSyncLending)

	http.HandleFunc("/histories/lending", handlerLendingHistories)

	addr := s.addr + ":" + strconv.Itoa(s.port)
	fmt.Println("run server on", addr)
	http.ListenAndServe(addr, nil)
}
