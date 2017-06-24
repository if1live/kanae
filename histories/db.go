package histories

import (
	"strconv"
	"time"

	"strings"

	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/thrasher-/gocryptotrader/exchanges/poloniex"

	"github.com/deckarep/golang-set"
)

type Database struct {
	db       *gorm.DB
	exchange *poloniex.Poloniex
}

func NewDatabase(filepath string, exchange *poloniex.Poloniex) (Database, error) {
	db, err := gorm.Open("sqlite3", filepath)
	if err != nil {
		return Database{}, err
	}

	db.AutoMigrate(&PoloniexTradeRow{})

	return Database{
		db:       db,
		exchange: exchange,
	}, nil
}

func (d *Database) GetAllTrades(asset, currency string) []PoloniexTradeRow {
	var rows []PoloniexTradeRow
	d.db.Where(&PoloniexTradeRow{Asset: asset, Currency: currency}).Order("Date desc").Find(&rows)
	return rows
}

// currency example : "all", "BTC_DOGE"
func (d *Database) LoadFromExchange(currency string) (int, error) {
	start := "0"
	end := strconv.FormatInt(time.Now().Unix(), 10)
	retval, err := d.exchange.GetAuthenticatedTradeHistory(currency, start, end)
	if err != nil {
		return -1, err
	}

	if all, ok := retval.(poloniex.PoloniexAuthenticatedTradeHistoryAll); ok {
		rowcount := 0
		for key, histories := range all.Data {
			// key example : BTC_DOGE
			tokens := strings.Split(key, "_")
			asset := tokens[1]
			currency := tokens[0]
			rowcount += d.insertTradeHistories(asset, currency, histories)
		}
		return rowcount, nil
	}

	if resp, ok := retval.(poloniex.PoloniexAuthenticatedTradeHistoryResponse); ok {
		tokens := strings.Split(currency, "_")
		rowcount := d.insertTradeHistories(tokens[1], tokens[0], resp.Data)
		return rowcount, nil
	}

	return -1, errors.New("unknown error : load from exchange")
}

func (d *Database) insertTradeHistories(asset, currency string, histories []poloniex.PoloniexAuthentictedTradeHistory) int {
	var existRows []PoloniexTradeRow
	d.db.Where(&PoloniexTradeRow{Asset: asset, Currency: currency}).Select("TradeID").Find(&existRows)
	tradeIDSet := mapset.NewSet()
	for _, r := range existRows {
		tradeIDSet.Add(r.TradeID)
	}

	rows := []PoloniexTradeRow{}
	for _, history := range histories {
		if tradeIDSet.Contains(history.TradeID) {
			continue
		}
		row := NewPoloniexTradeRow(asset, currency, history)
		rows = append(rows, row)
	}

	for _, row := range rows {
		d.db.Create(&row)
	}
	return len(rows)
}
