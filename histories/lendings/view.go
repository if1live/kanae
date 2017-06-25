package lendings

import (
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

func (q *View) All(currency string) []Lending {
	var rows []Lending
	q.db.Where(&Lending{Currency: currency}).Order("close desc").Find(&rows)
	return rows
}
