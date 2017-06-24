package main

import (
	"fmt"
	"path"

	"github.com/if1live/kanae/histories"
	"github.com/if1live/kanae/kanaelib"
)

func main() {
	filename := "config.yaml"
	filepath := path.Join(GetExecutablePath(), filename)
	s, err := kanaelib.LoadSettings(filepath)
	if err != nil {
		check(err)
	}

	exchange := s.PoloniexExchange()
	db, err := histories.NewDatabase("histories.db", exchange)
	if err != nil {
		check(err)
	}

	// get all
	rows := db.GetAllTrades("DOGE", "BTC")
	fmt.Println(rows)

	//if len(rows) == 0 {
	//rowcount, err := db.LoadFromExchange("all")
	//if err != nil {
	//	check(err)
	//}
	//fmt.Printf("%d row added\n", rowcount)
	//}
}
