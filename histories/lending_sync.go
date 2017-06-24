package histories

import (
	"strconv"
	"time"

	"github.com/deckarep/golang-set"
	"github.com/jinzhu/gorm"
	"github.com/thrasher-/gocryptotrader/exchanges/poloniex"
)

type LendingSync struct {
	db  *gorm.DB
	api *poloniex.Poloniex
}

func NewLendingSync(db *gorm.DB, api *poloniex.Poloniex) *LendingSync {
	return &LendingSync{
		db:  db,
		api: api,
	}
}

func (sync *LendingSync) Sync(start, end time.Time) (int, error) {
	startTime := strconv.FormatInt(start.Unix(), 10)
	endTime := strconv.FormatInt(end.Unix(), 10)
	retval, err := GetLendingHistory(sync.api, startTime, endTime)
	if err != nil {
		return -1, err
	}

	var existRows []LendingRow
	sync.db.Select("lending_id").Find(&existRows)
	idSet := mapset.NewSet()
	for _, r := range existRows {
		idSet.Add(r.LendingID)
	}

	rows := []LendingRow{}
	for _, history := range retval {
		if idSet.Contains(history.ID) {
			continue
		}
		row := NewLendingRow(history)
		rows = append(rows, row)
	}
	for _, row := range rows {
		sync.db.Create(&row)
	}
	return len(rows), nil
}

func (sync *LendingSync) SyncAll() (int, error) {
	start := time.Unix(0, 0)
	end := time.Now()
	return sync.Sync(start, end)
}

func (sync *LendingSync) SyncRecent() (int, error) {
	start := sync.GetLastTime()
	end := time.Now()
	return sync.Sync(start, end)
}

func (sync *LendingSync) GetLastTime() time.Time {
	var last LendingRow
	sync.db.Order("open desc").First(&last)
	if last.ID == 0 {
		return time.Unix(0, 0)
	}
	return last.Open
}
