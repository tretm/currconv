package storage

import (
	"errors"
	"fmt"
	"sync"

	"currconv/cc"
)

// MapCurrencyList structure for storing available courses
type MapCurrencyList struct {
	curlist map[string]map[string]cc.MapCurrencyPrice
	mut     *sync.Mutex
}

func New() *MapCurrencyList {
	var mcl MapCurrencyList = MapCurrencyList{}
	mcl.curlist = make(map[string]map[string]cc.MapCurrencyPrice)
	mcl.mut = &sync.Mutex{}
	return &mcl
}

// GetCurrencyPrice returns the exchange rate of a specific currency for a specific country
func (cl *MapCurrencyList) GetCurrencyPrice(country, currancy, saleOrBuy string) (float64, error) {
	var result float64
	res, err := cl.GetCurrencyCountry(country)
	if err != nil {
		return 0, err
	}
	if c, ok := res.(map[string]cc.MapCurrencyPrice)[currancy]; ok {
		if saleOrBuy == cc.BUY {
			result = c.Buy
		} else if saleOrBuy == cc.SALE {
			result = c.Sale
		}
		return result, nil
	} else {
		return 0, errors.New(fmt.Sprintf("Currency price for %s not found ", currancy))
	}

}

// GetCurrencyCountry gets the exchange rates of a specific country
func (cl *MapCurrencyList) GetCurrencyCountry(country string) (cc.ResultStorage, error) {
	cl.mut.Lock()
	defer cl.mut.Unlock()
	var res cc.ResultStorage
	if c, ok := cl.curlist[country]; ok {
		res = c
		return res, nil
	} else {
		return nil, errors.New("Country not found")
	}
}

// Set sets the result in the store obtained after parsing the data
func (cl *MapCurrencyList) Set(parcfunc cc.ParcerFunc, country, currencyname string) error {
	currencymap := make(map[string]cc.MapCurrencyPrice)
	parcfunc(currencymap)
	currencymap[currencyname] = cc.MapCurrencyPrice{}
	cl.mut.Lock()
	defer cl.mut.Unlock()
	if len(currencymap) > 1 {
		cl.curlist[country] = currencymap
		return nil
	}
	return errors.New("error set data")
}

// Update updating the result in the storage obtained after parsing the data
func (cl *MapCurrencyList) Update(dat interface{}) error {
	return nil
}

// Stop destructor
func (cl *MapCurrencyList) Stop() {
	cl.mut.Lock()
	defer cl.mut.Unlock()
	for key, val := range cl.curlist {
		for innerkey := range val {
			delete(val, innerkey)
		}
		delete(cl.curlist, key)
	}
}
