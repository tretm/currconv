package cc

import "fmt"

const (
	BUY            string = "BUY"
	SALE           string = "SALE"
	TYPESTORAGEMAP string = "TYPESTORAGEMAP"
	DEFALUTCOUNTRY string = "RU"
)

// Debug go signalcoucher(exitch)
var (
	Debug bool
)

// Println customized fmt.Println
func Println(a ...interface{}) {
	if Debug {
		fmt.Println(a...)
	}
}

// ResultStorage designed to return the result of courses by country
type ResultStorage interface{}

// ParcerFunc data type representing a generic function for parsing data for a specific source
type ParcerFunc func(c interface{})

// MapCurrencyPrice currency buying and selling rate
type MapCurrencyPrice struct {
	Buy  float64
	Sale float64
}

// Storage store interface for implementing data stores in various ways
type Storage interface {
	//New
	GetCurrencyPrice(country, currancy, saleOrBy string) (float64, error)
	GetCurrencyCountry(country string) (ResultStorage, error)
	Set(parcfunc ParcerFunc, country, currencyname string) error
	Update(dat interface{}) error
	Stop()
}
