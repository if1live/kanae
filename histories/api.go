package histories

import (
	"fmt"
	"time"

	"strconv"

	"github.com/if1live/kanae/exchanges"
	"github.com/if1live/kanae/kanaelib"
	"github.com/if1live/kanae/lendings"
	"github.com/jinzhu/gorm"
)

// poloniex api server clone
type APIServer struct {
	db *gorm.DB
}

func NewAPIServer(db *gorm.DB) *APIServer {
	return &APIServer{
		db: db,
	}
}

const (
	tradeHistoryTypeSellsOnly    = 0
	tradeHistoryTypeBuysOnly     = 1
	tradeHistoryTypeBuysAndSells = 2
	tradeHistoryTypeLoadEarnings = 3
)

// https://poloniex.com/private.php
// command=returnPaginatedTradeHistory
// start=0
// end=1917895987
// page=1
// tradesPerPage=50
// type=0
func (s *APIServer) PaginateTradeHistory(start, end time.Time, page, tadePerPage int, apiType int) []TradeHistory {
	if apiType == tradeHistoryTypeLoadEarnings {
		var rows []lendings.Lending
		q := s.db.Where("close between ? and ?", start, end)
		q.Order("close desc").Find(&rows)

		histories := make([]TradeHistory, len(rows))
		for i, r := range rows {
			histories[i] = NewTradeHistoryFromLending(&r)
		}
		return histories

	} else {
		var rows []exchanges.Exchange

		q := s.db.Where("date between ? and ?", start, end)
		switch apiType {
		case tradeHistoryTypeBuysOnly:
			q = q.Where("type = ?", exchanges.ExchangeBuy)
		case tradeHistoryTypeSellsOnly:
			q = q.Where("type = ?", exchanges.ExchangeSell)
		case tradeHistoryTypeBuysAndSells:
		}
		q.Order("date desc").Find(&rows)

		histories := make([]TradeHistory, len(rows))
		for i, r := range rows {
			histories[i] = NewTradeHistoryFromExchange(&r)
		}
		return histories
	}
}

// https://poloniex.com/private.php
// command=returnPersonalTradeHistory
// start=0
// end=1917895987
// retval: key=BTC, BTC_AMP, value=array
func (s *APIServer) PersonalTradeHistory(start, end time.Time) map[string][]PersonalTradeHistory {
	var exchangerows []exchanges.Exchange
	s.db.Where("date between ? and ?", start, end).Order("date desc").Find(&exchangerows)
	exchangehistories := make([]PersonalTradeHistory, len(exchangerows))
	for i, r := range exchangerows {
		h := NewPersonalTradeHistoryFromExchange(&r)
		exchangehistories[i] = h
	}

	// tradeId - currencyPair
	tradeIDMap := map[string]string{}
	for _, r := range exchangerows {
		key := strconv.FormatInt(r.TradeID, 10)
		tradeIDMap[key] = r.CurrencyPair()
	}

	var lendingrows []lendings.Lending
	s.db.Where("close between ? and ?", start, end).Order("close desc").Find(&lendingrows)
	lendinghistories := make([]PersonalTradeHistory, len(lendingrows))
	for i, r := range lendingrows {
		h := NewPersonalTradeHistoryFromLending(&r)
		lendinghistories[i] = h
	}

	// lending id - currency
	lendingIDMap := map[string]string{}
	for _, r := range lendingrows {
		key := "s" + strconv.FormatInt(r.LendingID, 10)
		lendingIDMap[key] = r.Currency
	}

	retval := map[string][]PersonalTradeHistory{}
	for _, h := range exchangehistories {
		currencyPair := tradeIDMap[h.TradeID]
		list, ok := retval[currencyPair]
		if ok {
			retval[currencyPair] = append(list, h)
		} else {
			retval[currencyPair] = []PersonalTradeHistory{h}
		}
	}

	for _, h := range lendinghistories {
		currency := lendingIDMap[h.TradeID]
		list, ok := retval[currency]
		if ok {
			retval[currency] = append(list, h)
		} else {
			retval[currency] = []PersonalTradeHistory{h}
		}
	}
	return retval
}

type TradeHistory struct {
	CurrencyPair string `json:"currencyPair"`
	Date         string `json:"date"`
	Type         string `json:"type"`
	Category     string `json:"category"`
	Rate         string `json:"rate"`
	Amount       string `json:"amount"`
	Total        string `json:"total"`
	Fee          string `json:"fee"`
}

func NewTradeHistoryFromExchange(r *exchanges.Exchange) TradeHistory {
	feeAmountStr := kanaelib.ToFloatStr(r.FeeAmount())
	feeUnit := ""
	if r.Type == exchanges.ExchangeBuy {
		feeUnit = r.Asset
	} else if r.Type == exchanges.ExchangeSell {
		feeUnit = r.Currency
	}
	feePercent := r.Fee * 100
	fee := fmt.Sprintf("%s %s (%.2f%%)", feeAmountStr, feeUnit, feePercent)

	return TradeHistory{
		CurrencyPair: r.CurrencyPair(),
		Type:         r.Type,
		Category:     r.Category,
		Amount:       kanaelib.ToFloatStr(r.Amount),
		Rate:         kanaelib.ToFloatStr(r.Rate),
		Date:         r.DateStr(),
		Total:        kanaelib.ToFloatStr(r.MyTotal()) + " " + r.Currency,
		Fee:          fee,
	}
}
func NewTradeHistoryFromLending(r *lendings.Lending) TradeHistory {
	return TradeHistory{
		CurrencyPair: r.Currency,
		Category:     "lendingEarning",
		Type:         "1",
		Date:         r.Close.Format("2006-01-02 15:04:05"),
		Amount:       kanaelib.ToFloatStr(r.Amount),
		Fee:          fmt.Sprintf("%.0f%%", r.FeeRate()*100),
		Rate:         fmt.Sprintf("%.4f%%", r.Rate*100),
		Total:        kanaelib.ToFloatStr(r.Interest),
	}
}

type PersonalTradeHistory struct {
	Amount        string `json:"amount"`
	Category      string `json:"category"`
	Date          string `json:"date"`
	Fee           string `json:"fee"`
	GlobalTradeID string `json:"globalTradeID,int"`
	OrderNumber   string `json:"orderNumber"`
	Rate          string `json:"rate"`
	Total         string `json:"total"`
	TradeID       string `json:"tradeID"`
	Type          string `json:"type"`
}

func NewPersonalTradeHistoryFromExchange(r *exchanges.Exchange) PersonalTradeHistory {
	return PersonalTradeHistory{
		Amount:        kanaelib.ToFloatStr(r.Amount),
		Category:      r.Category,
		Date:          r.DateStr(),
		Fee:           kanaelib.ToFloatStr(r.Fee),
		GlobalTradeID: strconv.FormatInt(r.GlobalTradeID, 10),
		OrderNumber:   strconv.FormatInt(r.OrderNumber, 10),
		Rate:          kanaelib.ToFloatStr(r.Rate),
		Total:         kanaelib.ToFloatStr(r.Total),
		TradeID:       strconv.FormatInt(r.TradeID, 10),
		Type:          r.Type,
	}
}

func NewPersonalTradeHistoryFromLending(r *lendings.Lending) PersonalTradeHistory {
	return PersonalTradeHistory{
		Amount:        kanaelib.ToFloatStr(r.Amount),
		Category:      "lendingEarning",
		Date:          r.Close.Format("2006-01-02 15:04:05"),
		Fee:           kanaelib.ToFloatStr(r.FeeRate()),
		GlobalTradeID: "",
		OrderNumber:   "",
		Rate:          kanaelib.ToFloatStr(r.Rate),
		Total:         kanaelib.ToFloatStr(r.Interest),
		TradeID:       "s" + strconv.FormatInt(r.LendingID, 10),
		Type:          "buy",
	}
}
