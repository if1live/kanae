package reports

import (
	"github.com/if1live/kanae/histories"
)

type BalanceReport struct {
	Currency string
	rows     []histories.BalanceRow
}

func NewBalanceReport(currency string, rows []histories.BalanceRow) BalanceReport {
	return BalanceReport{
		Currency: currency,
		rows:     rows,
	}
}

func (r *BalanceReport) Balance() float64 {
	var total float64
	total += r.Deposits()
	total -= r.Withdrawals()
	return total
}

func (r *BalanceReport) Deposits() float64 {
	var total float64
	for _, row := range r.rows {
		if row.Type == histories.BalanceTypeDeposit {
			total += row.Amount
		}
	}
	return total
}

func (r *BalanceReport) Withdrawals() float64 {
	var total float64
	for _, row := range r.rows {
		if row.Type == histories.BalanceTypeWithdrawal {
			total += row.Amount
		}
	}
	return total
}
