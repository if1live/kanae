package histories

import (
	"github.com/jinzhu/gorm"
)

type LendingView struct {
	db *gorm.DB
}

func NewLendingView(db *gorm.DB) *LendingView {
	return &LendingView{
		db: db,
	}
}

func (q *LendingView) GetAll(currency string) []LendingRow {
	var rows []LendingRow
	q.db.Where(&LendingRow{Currency: currency}).Order("close desc").Find(&rows)
	return rows
}
