package balances

type Report struct {
	Currency string
	rows     []Transaction
}

func NewReport(currency string, rows []Transaction) Report {
	return Report{
		Currency: currency,
		rows:     rows,
	}
}

func (r *Report) Balance() float64 {
	var total float64
	total += r.Deposits()
	total -= r.Withdrawals()
	return total
}

func (r *Report) Deposits() float64 {
	var total float64
	for _, row := range r.rows {
		if row.Type == TypeDeposit {
			total += row.Amount
		}
	}
	return total
}

func (r *Report) Withdrawals() float64 {
	var total float64
	for _, row := range r.rows {
		if row.Type == TypeWithdrawal {
			total += row.Amount
		}
	}
	return total
}
