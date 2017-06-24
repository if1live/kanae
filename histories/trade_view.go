package histories

import "github.com/jinzhu/gorm"

type TradeView struct {
	db *gorm.DB
}

func NewTradeView(db *gorm.DB) *TradeView {
	return &TradeView{
		db: db,
	}
}

func (q *TradeView) GetAll(asset, currency string) []TradeRow {
	var rows []TradeRow
	q.db.Where("asset = ? and currency = ?", asset, currency).Order("date desc").Find(&rows)
	return rows
}
