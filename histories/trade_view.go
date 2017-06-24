package histories

import (
	"sort"

	"github.com/jinzhu/gorm"
)

type TradeView struct {
	db *gorm.DB
}

func NewTradeView(db *gorm.DB) *TradeView {
	return &TradeView{
		db: db,
	}
}

func (q *TradeView) All(asset, currency string) []TradeRow {
	var rows []TradeRow
	q.db.Where("asset = ? and currency = ?", asset, currency).Order("date desc").Find(&rows)
	return rows
}

func (v *TradeView) UsedAssets(currency string) []string {
	var rows []TradeRow
	v.db.Where("currency = ?", currency).Select("DISTINCT(asset)").Find(&rows)

	assets := []string{}
	for _, r := range rows {
		assets = append(assets, r.Asset)
	}
	sort.Strings(assets)
	return assets
}
