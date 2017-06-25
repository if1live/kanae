package exchanges

import (
	"sort"

	"github.com/jinzhu/gorm"
)

type View struct {
	db *gorm.DB
}

func NewView(db *gorm.DB) *View {
	return &View{
		db: db,
	}
}

func (q *View) Get(asset, currency string) []Exchange {
	var rows []Exchange
	q.db.Where("asset = ? and currency = ?", asset, currency).Order("date desc").Find(&rows)
	return rows
}

func (v *View) All() []Exchange {
	var rows []Exchange
	v.db.Order("date desc").Find(&rows)
	return rows
}

func (v *View) UsedAssets(currency string) []string {
	var rows []Exchange
	v.db.Where("currency = ?", currency).Select("DISTINCT(asset)").Find(&rows)

	assets := []string{}
	for _, r := range rows {
		assets = append(assets, r.Asset)
	}
	sort.Strings(assets)
	return assets
}
