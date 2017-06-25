package balances

import "github.com/jinzhu/gorm"

type View struct {
	db *gorm.DB
}

func NewView(db *gorm.DB) *View {
	return &View{
		db: db,
	}
}

func (v *View) CurrencyRows(currency string) []Transaction {
	var rows []Transaction
	v.db.Where(&Transaction{Currency: currency}).Order("timestamp desc").Find(&rows)
	return rows
}
