package commands

import (
	"html/template"
	"path"
	"strings"

	"os"

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
	/*
		db, err := histories.NewDatabase(cmd.settings.DatabaseFileName)
		if err != nil {
			return nil
		}

		view := db.MakeTradeView()
		assets := view.UsedAssets("BTC")
		fmt.Println("used assets :", assets)
	*/

	fmap := template.FuncMap{
		"title": strings.Title,
	}
	ctx := map[string]string{
		"Title":   "Hello world",
		"Content": "Hi there",
	}
	basePath := kanaelib.GetExecutablePath()
	fp := path.Join(basePath, "web", "templates", "sample.html")
	tpl, err := template.New("sample.html").Funcs(fmap).ParseFiles(fp)
	if err != nil {
		panic(err)
	}
	if err := tpl.Execute(os.Stdout, ctx); err != nil {
		panic(err)
	}

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
