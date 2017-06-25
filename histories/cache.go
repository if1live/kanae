package histories

import (
	"time"

	"errors"

	"github.com/thrasher-/gocryptotrader/exchanges/poloniex"
)

type TickerCache struct {
	// key example : BTC_LTC
	tickers   map[string]poloniex.PoloniexTicker
	UpdatedAt time.Time

	api *poloniex.Poloniex
}

func NewTickerCache(api *poloniex.Poloniex) *TickerCache {
	return &TickerCache{
		tickers:   make(map[string]poloniex.PoloniexTicker),
		UpdatedAt: time.Unix(0, 0),
		api:       api,
	}
}

func (c *TickerCache) Refresh() error {
	tickers, err := c.api.GetTicker()
	if err != nil {
		return err
	}

	c.UpdatedAt = time.Now()
	c.tickers = tickers
	return nil
}

func (c *TickerCache) Get(asset, currency string) (poloniex.PoloniexTicker, error) {
	key := currency + "_" + asset
	ticker, ok := c.tickers[key]
	if !ok {
		return poloniex.PoloniexTicker{}, errors.New("invalid asset-currency: " + key)
	}
	return ticker, nil
}
