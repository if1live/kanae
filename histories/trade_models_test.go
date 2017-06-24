package histories

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TradeRow_FeeAmount(t *testing.T) {
	/*
		sample data

		from poloniex trade history
		type: Buy
		price/share: 0.00024092
		amount: 20.75377718
		fee: 0.05188444 AMP (0.25%)
		total: 0.00499999 BTC

		from API
		AMP	BTC
		rate: 0.00024092
		amount: 20.75377718
		total: 0.00499999
		fee: 0.0025
		type: buy
	*/

	cases := []struct {
		rate      float64
		amount    float64
		fee       float64
		feeAmount float64
	}{
		{
			0.00024092,
			20.75377718,
			0.0025,
			0.05188444,
		},
		{
			1.29e-06,
			7751.93798449,
			0.0025,
			19.37984496,
		},
	}

	for _, c := range cases {
		row := TradeRow{
			Rate:   c.rate,
			Amount: c.amount,
			Fee:    c.fee,
		}
		actual := row.FixedFeeAmount()
		assert.Equal(t, c.feeAmount, actual)
	}
}
