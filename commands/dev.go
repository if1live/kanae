package commands

import (
	"fmt"

	"github.com/if1live/kanae/kanaelib"
)

type Dev struct {
	settings kanaelib.Settings
}

func NewDev(s kanaelib.Settings) *Dev {
	return &Dev{
		settings: s,
	}
}

func (cmd *Dev) Execute() error {
	fmt.Println("TODO: Dev")
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
	return nil
}
