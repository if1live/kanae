package histories

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/thrasher-/gocryptotrader/exchanges/poloniex"
)

type Database struct {
	db *gorm.DB
}

func NewDatabase(filepath string) (Database, error) {
	db, err := gorm.Open("sqlite3", filepath)
	if err != nil {
		return Database{}, err
	}

	db.AutoMigrate(&TradeRow{})
	db.AutoMigrate(&LendingRow{})
	db.AutoMigrate(&BalanceRow{})

	return Database{
		db: db,
	}, nil
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) MakeTradeSync(api *poloniex.Poloniex) *TradeSync {
	return NewTradeSync(d.db, api)
}

func (d *Database) MakeLendingSync(api *poloniex.Poloniex) *LendingSync {
	return NewLendingSync(d.db, api)
}
func (d *Database) MakeBalanceSync(api *poloniex.Poloniex) *BalanceSync {
	return NewBalanceSync(d.db, api)
}

func (d *Database) MakeTradeView() *TradeView {
	return NewTradeView(d.db)
}

func (d *Database) MakeLendingView() *LendingView {
	return NewLendingView(d.db)
}
