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
	defer db.Close()

	// sync all history
	exchangeCount, err := db.SyncRecentExchange()
	if err != nil {
		check(err)
	}
	fmt.Println(exchangeCount, "exchange rows added")

	lendingCount, err := db.SyncRecentLending()
	if err != nil {
		check(err)
	}
	fmt.Println(lendingCount, "lending rows added")

	depositWithdrawCount, err := db.SyncRecentDepositWithdraw()
	if err != nil {
		check(err)
	}
	fmt.Println(depositWithdrawCount, "deposit/withdraw rows added")

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
