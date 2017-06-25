package commands

import (
	"fmt"

	"github.com/if1live/kanae/histories"
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
	db, err := histories.NewDatabase(cmd.settings.DatabaseFileName)
	if err != nil {
		return nil
	}

	view := db.MakeTradeView()
	assets := view.UsedAssets("BTC")
	fmt.Println("used assets :", assets)

	//rows := db.GetLendings("BTC")
	//fmt.Println(rows)

	/*
		exchange := s.PoloniexExchange()
		db, err := histories.NewDatabase("histories.db", exchange)


		// get all
		rows := db.AllTrades("DOGE", "BTC")
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
