package histories

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/thrasher-/gocryptotrader/exchanges/poloniex"
)

const (
	BalanceTypeDeposit    = "deposit"
	BalanceTypeWithdrawal = "withdrawal"
)

// merge two struct
// - poloniex.PoloniexDepositsWithdrawals.Deposits
// - poloniex.PoloniexDepositsWithdrawals.Withdrawals
type BalanceRow struct {
	gorm.Model

	// deposit / withdrawal
	Type string

	WithdrawalNumber int64 // only withdrawals
	Currency         string
	Address          string
	Amount           float64
	Confirmations    int
	TransactionID    string `gorm:"unique"`
	Timestamp        time.Time
	Status           string
	IPAddress        string // only withdrawals
}

func NewBalanceRows(h poloniex.PoloniexDepositsWithdrawals) []BalanceRow {
	deposits := []BalanceRow{}
	for _, row := range h.Deposits {
		r := BalanceRow{
			Type:          BalanceTypeDeposit,
			Currency:      row.Currency,
			Address:       row.Address,
			Amount:        row.Amount,
			Confirmations: row.Confirmations,
			TransactionID: row.TransactionID,
			Timestamp:     time.Unix(row.Timestamp, 0),
			Status:        row.Status,
		}
		deposits = append(deposits, r)
	}

	withdrawals := []BalanceRow{}
	for _, row := range h.Withdrawals {
		r := BalanceRow{
			Type:             BalanceTypeWithdrawal,
			WithdrawalNumber: row.WithdrawalNumber,
			Currency:         row.Currency,
			Address:          row.Address,
			Amount:           row.Amount,
			Confirmations:    row.Confirmations,
			TransactionID:    row.TransactionID,
			Timestamp:        time.Unix(row.Timestamp, 0),
			Status:           row.Status,
			IPAddress:        row.IPAddress,
		}
		withdrawals = append(withdrawals, r)
	}

	rows := []BalanceRow{}
	rows = append(rows, deposits...)
	rows = append(rows, withdrawals...)
	return rows
}
