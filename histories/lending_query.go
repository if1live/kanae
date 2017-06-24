package histories

import (
	"github.com/jinzhu/gorm"
)

type LendingQuery struct {
	db *gorm.DB
}

func NewLendingQuery(db *gorm.DB) *LendingQuery {
	return &LendingQuery{
		db: db,
	}
}

func (q *LendingQuery) GetAll(currency string) []LendingRow {
	var rows []LendingRow
	q.db.Where(&LendingRow{Currency: currency}).Order("close desc").Find(&rows)
	return rows
}
