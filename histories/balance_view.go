package histories

import "github.com/jinzhu/gorm"

type BalanceView struct {
	db *gorm.DB
}

func NewBalanceView(db *gorm.DB) *BalanceView {
	return &BalanceView{
		db: db,
	}
}

func (v *BalanceView) CurrencyRows(currency string) []BalanceRow {
	var rows []BalanceRow
	v.db.Where(&BalanceRow{Currency: currency}).Order("timestamp desc").Find(&rows)
	return rows
}
