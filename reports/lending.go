package reports

import "github.com/if1live/kanae/histories/lendings"

type LendingReport struct {
	rows []lendings.Lending
}

func NewLendingReport(rows []lendings.Lending) LendingReport {
	return LendingReport{
		rows: rows,
	}
}
