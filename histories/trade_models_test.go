package histories

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"strconv"

	"github.com/if1live/kanae/kanaelib"
)

func Test_TradeRow_buyFeeAmount(t *testing.T) {
	// buy example
	// amount : 13.00373802
	// fee : 0.03250935 SYS (0.25%)
	// 0.03250935 SYS = 13.00373802 SYS * (0.01) * (0.25)
	cases := []struct {
		amountStr    string
		feeStr       string
		feeAmountStr string
	}{
		{"13.00373802", "0.25", "0.03250935"},
		{"20.75377718", "0.25", "0.05188444"},
	}
	for _, c := range cases {
		amount, _ := strconv.ParseFloat(c.amountStr, 64)
		fee, _ := strconv.ParseFloat(c.feeStr, 64)
		r := TradeRow{
			Amount: amount,
			Fee:    fee * float64(0.01),
		}
		v := r.buyFeeAmount()
		assert.Equal(t, c.feeAmountStr, kanaelib.ToFloatStr(v))
	}
}

func Test_TradeRow_sellFeeAmount(t *testing.T) {
	cases := []struct {
		totalStr     string
		feeStr       string
		feeAmountStr string
	}{
		{"0.01085732", "0.15", "0.00001629"},
	}
	for _, c := range cases {
		total, _ := strconv.ParseFloat(c.totalStr, 64)
		fee, _ := strconv.ParseFloat(c.feeStr, 64)
		r := TradeRow{
			Total: total,
			Fee:   fee * float64(0.01),
		}
		v := r.sellFeeAmount()
		assert.Equal(t, c.feeAmountStr, kanaelib.ToFloatStr(v))
	}
}
