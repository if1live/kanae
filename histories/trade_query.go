package histories

import "github.com/jinzhu/gorm"

type TradeQuery struct {
	db *gorm.DB
}

func NewTradeQuery(db *gorm.DB) *TradeQuery {
	return &TradeQuery{
		db: db,
	}
}

func (q *TradeQuery) GetAll(asset, currency string) []TradeRow {
	var rows []TradeRow
	q.db.Where("asset = ? and currency = ?", asset, currency).Order("date desc").Find(&rows)
	return rows
}
