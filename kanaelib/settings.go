package kanaelib

import (
	"io/ioutil"

	"github.com/thrasher-/gocryptotrader/config"
	"github.com/thrasher-/gocryptotrader/exchanges/poloniex"

	yaml "gopkg.in/yaml.v2"
)

type Settings struct {
	APIKey    string `yaml:"poloniex_api_key"`
	APISecret string `yaml:"poloniex_api_secret"`
}

func LoadSettings(filepath string) (Settings, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return Settings{}, err
	}

	var s Settings
	err = yaml.Unmarshal(data, &s)
	if err != nil {
		return Settings{}, err
	}

	return s, nil
}

func (s *Settings) PoloniexExchange() *poloniex.Poloniex {
	conf := config.ExchangeConfig{
		Enabled:                 true,
		APIKey:                  s.APIKey,
		APISecret:               s.APISecret,
		AuthenticatedAPISupport: true,
		Verbose:                 true,
	}
	exchange := poloniex.Poloniex{}
	exchange.Setup(conf)
	return &exchange
}