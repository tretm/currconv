package parser

import (
	"currconv/cc"
	"math"
)

// RUCurr data type for parsing the Russian central bank
type RUCurr struct {
	CharCode string  `xml:"CharCode"`
	Value    float64 `xml:"Value"`
	Nominal  float64 `xml:"Nominal"`
}

// RUCurrency data type for parsing the Russian central bank
type RUCurrency struct {
	Url    string
	Valute []RUCurr `xml:"Valute"`
}

// decode Russian exchange rate parser
func (ru *RUCurrency) decode(currencyMap interface{}) {
	switch curM := currencyMap.(type) {
	case map[string]cc.MapCurrencyPrice:
		for _, item := range ru.Valute {
			b := math.Round((item.Value/item.Nominal)*10000) / 10000
			curM[item.CharCode] = cc.MapCurrencyPrice{Buy: b, Sale: b}
		}
	}

}
