package reports

import (
	"github.com/if1live/kanae/histories"
)

type LendingReport struct {
	rows []histories.LendingRow
}

func NewLendingReport(rows []histories.LendingRow) LendingReport {
	return LendingReport{
		rows: rows,
	}
}
