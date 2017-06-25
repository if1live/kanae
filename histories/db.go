package histories

import (
	"github.com/if1live/kanae/histories/balances"
	"github.com/if1live/kanae/histories/exchanges"
	"github.com/if1live/kanae/histories/lendings"
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

	db.AutoMigrate(&exchanges.Exchange{})
	db.AutoMigrate(&lendings.Lending{})
	db.AutoMigrate(&balances.Transaction{})

	return Database{
		db: db,
	}, nil
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) MakeExchangeSync(api *poloniex.Poloniex) *exchanges.Sync {
	return exchanges.NewSync(d.db, api)
}

func (d *Database) MakeLendingSync(api *poloniex.Poloniex) *lendings.Sync {
	return lendings.NewSync(d.db, api)
}
func (d *Database) MakeBalanceSync(api *poloniex.Poloniex) *balances.Sync {
	return balances.NewSync(d.db, api)
}

func (d *Database) MakeExchangeView() *exchanges.View {
	return exchanges.NewView(d.db)
}

func (d *Database) MakeLendingView() *lendings.View {
	return lendings.NewView(d.db)
}

func (d *Database) MakeBalanceView() *balances.View {
	return balances.NewView(d.db)
}
