package histories

import (
	"strconv"
	"time"

	"github.com/deckarep/golang-set"
	"github.com/jinzhu/gorm"
	"github.com/thrasher-/gocryptotrader/exchanges/poloniex"
)

type BalanceSync struct {
	db  *gorm.DB
	api *poloniex.Poloniex
}

func NewBalanceSync(db *gorm.DB, api *poloniex.Poloniex) *BalanceSync {
	return &BalanceSync{
		db:  db,
		api: api,
	}
}

func (sync *BalanceSync) Sync(start, end time.Time) (int, error) {
	startTime := strconv.FormatInt(start.Unix(), 10)
	endTime := strconv.FormatInt(end.Unix(), 10)
	retval, err := sync.api.GetDepositsWithdrawals(startTime, endTime)
	if err != nil {
		return -1, err
	}

	var existRows []BalanceRow
	sync.db.Select("transaction_id").Find(&existRows)
	idSet := mapset.NewSet()
	for _, r := range existRows {
		idSet.Add(r.TransactionID)
	}

	rows := []BalanceRow{}
	retvals := NewBalanceRows(retval)
	for _, history := range retvals {
		if idSet.Contains(history.TransactionID) {
			continue
		}
		rows = append(rows, history)
	}
	for _, row := range rows {
		sync.db.Create(&row)
	}
	return len(rows), nil
}

func (sync *BalanceSync) SyncAll() (int, error) {
	start := time.Unix(0, 0)
	end := time.Now()
	return sync.Sync(start, end)
}

func (sync *BalanceSync) SyncRecent() (int, error) {
	start := sync.GetLastTime()
	end := time.Now()
	return sync.Sync(start, end)
}
func (sync *BalanceSync) GetLastTime() time.Time {
	var last BalanceRow
	sync.db.Order("timestamp desc").First(&last)
	if last.ID == 0 {
		return time.Unix(0, 0)
	}
	return last.Timestamp
}
