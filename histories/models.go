package histories

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/thrasher-/gocryptotrader/exchanges/poloniex"
)

const (
	RecordTypeDeposit     = "deposit"
	RecordTypeWithdrawals = "withdrawals"
)

func convertPoloniexDate(val string) time.Time {
	// date example : 2017-06-18 04:31:08
	t, _ := time.Parse("2006-01-02 15:04:05", val)
	return t
}

// PoloniexAuthenticatedTradeHistory
type PoloniexTradeRow struct {
	gorm.Model

	// from zenbot format
	// {exchange_slug}.{asset}-{currency}
	// poloniex.AMP-BTC
	Asset    string
	Currency string

	GlobalTradeID int64 `gorm:"unique"`
	TradeID       int64
	Date          time.Time
	Rate          float64
	Amount        float64
	Total         float64
	Fee           float64
	OrderNumber   int64
	Type          string
	Category      string
}

func NewPoloniexTradeRow(asset, currency string, h poloniex.PoloniexAuthentictedTradeHistory) PoloniexTradeRow {
	return PoloniexTradeRow{
		Asset:    asset,
		Currency: currency,

		GlobalTradeID: h.GlobalTradeID,
		TradeID:       h.TradeID,
		Date:          convertPoloniexDate(h.Date),
		Rate:          h.Rate,
		Amount:        h.Amount,
		Total:         h.Total,
		Fee:           h.Fee,
		OrderNumber:   h.OrderNumber,
		Type:          h.Type,
		Category:      h.Category,
	}
}

type PoloniexLendingRow struct {
	gorm.Model

	LendingID int64 `gorm:"unique"`
	Currency  string
	Rate      float64
	Amount    float64
	Duration  float64
	Interest  float64
	Fee       float64
	Earned    float64
	Open      time.Time
	Close     time.Time
}

func NewPoloniexLendingRow(h PoloniexLendingHistory) PoloniexLendingRow {
	return PoloniexLendingRow{
		LendingID: h.ID,
		Currency:  h.Currency,
		Rate:      h.Rate,
		Amount:    h.Amount,
		Duration:  h.Duration,
		Interest:  h.Interest,
		Fee:       h.Fee,
		Earned:    h.Earned,
		Open:      convertPoloniexDate(h.Open),
		Close:     convertPoloniexDate(h.Close),
	}
}

// merge two struct
// - poloniex.PoloniexDepositsWithdrawals.Deposits
// - poloniex.PoloniexDepositsWithdrawals.Withdrawals
type PoloniexDepositWithdrawRow struct {
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

func NewPoloniexDepositWithdrawRows(h poloniex.PoloniexDepositsWithdrawals) []PoloniexDepositWithdrawRow {
	deposits := []PoloniexDepositWithdrawRow{}
	for _, row := range h.Deposits {
		r := PoloniexDepositWithdrawRow{
			Type:          RecordTypeDeposit,
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

	withdrawals := []PoloniexDepositWithdrawRow{}
	for _, row := range h.Withdrawals {
		r := PoloniexDepositWithdrawRow{
			Type:             RecordTypeWithdrawals,
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

	rows := []PoloniexDepositWithdrawRow{}
	rows = append(rows, deposits...)
	rows = append(rows, withdrawals...)
	return rows
}
