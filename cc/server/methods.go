package server

import (
	"math"
	"strings"

	"currconv/cc"
)

// count exchange rate calculation
func count(from, to float64) float64 {
	var result float64
	if from == 0 && to != 0 {
		result = 1 / to
	} else if to == 0 {
		result = from
	} else if to != 0 {
		result = from / to
	}
	result = math.Round(result*10000) / 10000
	return result

}

// convert exchange rate conversion depending on the country
func (h *HttpServer) convert(country, from, to string) (float64, error) {
	var f, t, res float64
	f, err := h.St.GetCurrencyPrice(country, from, cc.BUY)
	if err != nil {
		// err = fmt.Errorf("Wrong currency name from")
		return res, err
	}
	t, err = h.St.GetCurrencyPrice(country, to, cc.SALE)
	if err != nil {
		// err = fmt.Errorf("Wrong currency name to")
		return res, err
	}
	res = count(f, t)
	cc.Println(country, from, "/", to, res)
	return res, nil
}

// countrySelector used to select country
func countrySelector(country string) string {
	if country == "" {
		return cc.DEFALUTCOUNTRY
	} else {
		return strings.TrimSpace(country)
	}
}

// countAmount used to calculate the total value to be returned
func countAmount(price, amount float64) *ResponseAfterConvert {
	if amount < 1 {
		amount = 1
	}
	return &ResponseAfterConvert{Amount: math.Round((price*amount)*10000) / 10000}
}
