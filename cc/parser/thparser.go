package parser

import (
	"currconv/cc"
	"strings"
)

// THCurr data type for Thai central bank parsing
type THCurr struct {
	Title          string  `xml:"title"`
	Description    string  `xml:"description"`
	TargetCurrency string  `xml:"targetCurrency"`
	Value          float64 `xml:"value"`
}

// THCurrency data type for Thai central bank parsing
type THCurrency struct {
	Url   string
	Items []THCurr `xml:"item"`
}

// decode Thai exchange rate parser
func (th *THCurrency) decode(currencyMap interface{}) {
	switch mcc := currencyMap.(type) {
	case map[string]cc.MapCurrencyPrice:
		for _, item := range th.Items {
			if strings.Contains(item.Title, "Thailand Selling Rate") || strings.Contains(item.Title, "Thailand Average Selling Rate") {
				if curm, ok := mcc[item.TargetCurrency]; ok {
					curm.Sale = item.Value
					mcc[item.TargetCurrency] = curm
				} else {
					curm.Sale = item.Value
					mcc[item.TargetCurrency] = curm
				}
			} else if strings.Contains(item.Title, "Thailand Buying Rate") || strings.Contains(item.Title, "Thailand Average Buying Transfe") {
				if curm, ok := mcc[item.TargetCurrency]; ok {
					curm.Buy = item.Value
					mcc[item.TargetCurrency] = curm
				} else {
					curm.Buy = item.Value
					mcc[item.TargetCurrency] = curm
				}
			}
		}
	}
}
