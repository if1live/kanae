package main

import (
	"fmt"
	"path"

	"reflect"

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

	db, err := histories.NewDatabase("data.db")
	if err != nil {
		check(err)
	}
	defer db.Close()

	// sync all history
	exchange := s.PoloniexExchange()
	syncs := []histories.Synchronizer{
		db.MakeTradeSync(exchange),
		db.MakeLendingSync(exchange),
		db.MakeBalanceSync(exchange),
	}
	for _, sync := range syncs {
		rowcount, err := sync.SyncRecent()
		if err != nil {
			check(err)
		}

		syncName := reflect.TypeOf(sync).String()
		fmt.Printf("%s : %d exchange rows added\n", syncName, rowcount)
	}

	//rows := db.GetLendings("BTC")
	//fmt.Println(rows)

	/*
		exchange := s.PoloniexExchange()
		db, err := histories.NewDatabase("histories.db", exchange)


		// get all
		rows := db.GetAllTrades("DOGE", "BTC")
		fmt.Println(rows)
	*/

	//if len(rows) == 0 {
	//rowcount, err := db.LoadFromExchange("all")
	//if err != nil {
	//	check(err)
	//}
	//fmt.Printf("%d row added\n", rowcount)
	//}
}
