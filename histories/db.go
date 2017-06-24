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
	db.AutoMigrate(&PoloniexLendingRow{})
	db.AutoMigrate(&PoloniexDepositWithdrawRow{})

	return Database{
		db:       db,
		exchange: exchange,
	}, nil
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) GetAllTrades(asset, currency string) []PoloniexTradeRow {
	var rows []PoloniexTradeRow
	d.db.Where(&PoloniexTradeRow{Asset: asset, Currency: currency}).Order("date desc").Find(&rows)
	return rows
}

// currencyPair example : "all", "BTC_DOGE"
func (d *Database) SyncExchange(currencyPair string, start, end time.Time) (int, error) {
	startTime := strconv.FormatInt(start.Unix(), 10)
	endTime := strconv.FormatInt(end.Unix(), 10)
	retval, err := d.exchange.GetAuthenticatedTradeHistory(currencyPair, startTime, endTime)

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
		tokens := strings.Split(currencyPair, "_")
		rowcount := d.insertTradeHistories(tokens[1], tokens[0], resp.Data)
		return rowcount, nil
	}

	return -1, errors.New("unknown error : load from exchange")
}

func (d *Database) SyncAllExchange(currencyPair string) (int, error) {
	start := time.Unix(0, 0)
	end := time.Now()
	return d.SyncExchange(currencyPair, start, end)
}

func (d *Database) SyncRecentExchange() (int, error) {
	start := d.GetLastExchangeTime()
	end := time.Now()
	return d.SyncExchange("all", start, end)
}

func (d *Database) GetLastExchangeTime() time.Time {
	var last PoloniexTradeRow
	d.db.Order("date desc").First(&last)
	if last.ID == 0 {
		return time.Unix(0, 0)
	}
	return last.Date
}

func (d *Database) insertTradeHistories(asset, currency string, histories []poloniex.PoloniexAuthentictedTradeHistory) int {
	var existRows []PoloniexTradeRow
	d.db.Where(&PoloniexTradeRow{Asset: asset, Currency: currency}).Select("trade_id").Find(&existRows)
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

func (d *Database) GetLendings(currency string) []PoloniexLendingRow {
	var rows []PoloniexLendingRow
	d.db.Where(&PoloniexLendingRow{Currency: currency}).Order("close desc").Find(&rows)
	return rows
}

func (d *Database) SyncLending(start, end time.Time) (int, error) {
	startTime := strconv.FormatInt(start.Unix(), 10)
	endTime := strconv.FormatInt(end.Unix(), 10)
	retval, err := GetLendingHistory(d.exchange, startTime, endTime)
	if err != nil {
		return -1, err
	}

	var existRows []PoloniexLendingRow
	d.db.Select("lending_id").Find(&existRows)
	idSet := mapset.NewSet()
	for _, r := range existRows {
		idSet.Add(r.LendingID)
	}

	rows := []PoloniexLendingRow{}
	for _, history := range retval {
		if idSet.Contains(history.ID) {
			continue
		}
		row := NewPoloniexLendingRow(history)
		rows = append(rows, row)
	}
	for _, row := range rows {
		d.db.Create(&row)
	}
	return len(rows), nil
}

func (d *Database) SyncAllLending() (int, error) {
	start := time.Unix(0, 0)
	end := time.Now()
	return d.SyncLending(start, end)
}

func (d *Database) SyncRecentLending() (int, error) {
	start := d.GetLastLendingTime()
	end := time.Now()
	return d.SyncLending(start, end)
}
func (d *Database) GetLastLendingTime() time.Time {
	var last PoloniexLendingRow
	d.db.Order("open desc").First(&last)
	if last.ID == 0 {
		return time.Unix(0, 0)
	}
	return last.Open
}

func (d *Database) SyncDepositWithdraw(start, end time.Time) (int, error) {
	startTime := strconv.FormatInt(start.Unix(), 10)
	endTime := strconv.FormatInt(end.Unix(), 10)
	retval, err := d.exchange.GetDepositsWithdrawals(startTime, endTime)
	if err != nil {
		return -1, err
	}

	var existRows []PoloniexDepositWithdrawRow
	d.db.Select("transaction_id").Find(&existRows)
	idSet := mapset.NewSet()
	for _, r := range existRows {
		idSet.Add(r.TransactionID)
	}

	rows := []PoloniexDepositWithdrawRow{}
	retvals := NewPoloniexDepositWithdrawRows(retval)
	for _, history := range retvals {
		if idSet.Contains(history.TransactionID) {
			continue
		}
		rows = append(rows, history)
	}
	for _, row := range rows {
		d.db.Create(&row)
	}
	return len(rows), nil
}

func (d *Database) SyncAllDepositWithdraw() (int, error) {
	start := time.Unix(0, 0)
	end := time.Now()
	return d.SyncDepositWithdraw(start, end)
}

func (d *Database) SyncRecentDepositWithdraw() (int, error) {
	start := d.GetLastSyncDepositWithdrawTime()
	end := time.Now()
	return d.SyncDepositWithdraw(start, end)
}
func (d *Database) GetLastSyncDepositWithdrawTime() time.Time {
	var last PoloniexDepositWithdrawRow
	d.db.Order("timestamp desc").First(&last)
	if last.ID == 0 {
		return time.Unix(0, 0)
	}
	return last.Timestamp
}
