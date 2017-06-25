package lendings

type Report struct {
	rows []Lending
}

func NewReport(rows []Lending) Report {
	return Report{
		rows: rows,
	}
}
