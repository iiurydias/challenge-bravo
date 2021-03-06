package currency

import (
	customErrors "github.com/iiurydias/challenge-bravo/currency-rate-updater/service/errors"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
)

type currency struct {
	host string
}

func New(host string) Currency {
	return &currency{host: host}
}

// IT GETS SPECIFIC CURRENCY FROM THIRD PARTY APPLICATION AND HYDRATE IT
func (c *currency) GetCurrencyRate(code string) (float64, error) {
	result, err := c.getAllCurrenciesRate()
	if err != nil {
		return 0, err
	}
	value, ok := result[code]
	if !ok {
		return 0, customErrors.ErrInvalidCurrency
	}
	return value.Float(), nil
}

// IT GETS ALL CURRENCIES FROM THIRD PARTY APPLICATION AND HYDRATE IT
func (c *currency) GetAllCurrenciesRate() (map[string]float64, error) {
	currencies := make(map[string]float64)
	result, err := c.getAllCurrenciesRate()
	if err != nil {
		return nil, err
	}
	for key, value := range result {
		currencies[key] = value.Float()
	}
	return currencies, nil
}

// IT GETS ALL CURRENCIES FROM THIRD PARTY APPLICATION
func (c *currency) getAllCurrenciesRate() (map[string]gjson.Result, error) {
	response, err := http.Get(c.host)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all currencies rate")
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read body reader")
	}
	result := gjson.Get(string(body), "rates")
	return result.Map(), nil
}
