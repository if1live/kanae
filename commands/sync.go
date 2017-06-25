package commands

import (
	"fmt"
	"reflect"

	"github.com/if1live/kanae/histories"
	"github.com/if1live/kanae/kanaelib"
)

type Sync struct {
	settings kanaelib.Settings
}

func NewSync(s kanaelib.Settings) *Sync {
	return &Sync{
		settings: s,
	}
}
func (cmd *Sync) Execute() error {
	db, err := histories.NewDatabase(cmd.settings.DatabaseFileName)
	if err != nil {
		return err
	}
	defer db.Close()

	api := cmd.settings.MakePoloniex()
	syncs := []histories.Synchronizer{
		db.MakeExchangeSync(api),
		db.MakeLendingSync(api),
		db.MakeBalanceSync(api),
	}
	for _, sync := range syncs {
		rowcount, err := sync.SyncRecent()
		if err != nil {
			return err
		}

		syncName := reflect.TypeOf(sync).String()
		// TODO use logger
		fmt.Printf("%s : %d exchange rows added\n", syncName, rowcount)
	}

	return nil
}
