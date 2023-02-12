package parser

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"

	"currconv/cc"

	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
)

const (
	JSON string = "JSON"
	XML  string = "XML"
)

// sendRequest requesting data from a resource containing data for parsing
func sendRequest(url string) (*[]byte, string) {
	resp, err := http.Get(url)
	if err != nil {
		cc.Println("Error fetching data:", err)
		return nil, ""
	}
	defer resp.Body.Close()
	s := checkRespContentType(resp)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		cc.Println("Error reading data:", err)
		return nil, ""
	}
	return &body, s
}

// parceReqToStruct parses the received data
func parceReqToStruct(url string, any interface{}) error {
	body, treq := sendRequest(url)
	bb := bytes.Replace(*body, []byte(","), []byte("."), -1)
	reader := bytes.NewReader(bb)
	if treq == XML {
		decoder := xml.NewDecoder(reader)
		decoder.CharsetReader = charset.NewReaderLabel
		err := decoder.Decode(any)
		if err != nil {
			cc.Println("!!!!", err)
			return err
		}
	} else if treq == JSON {
		decoder := json.NewDecoder(reader)
		err := decoder.Decode(any)
		if err != nil {
			cc.Println("!!!!", err)
			return err
		}
	}
	return nil
}

// checkRespContentType checks the format of the response from the server
func checkRespContentType(r *http.Response) string {
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		return JSON
	} else if strings.Contains(contentType, "application/xml") || strings.Contains(contentType, "text/xml") {
		return XML
	}
	return ""
}

// Decoder decrypts the received data and data in the store about exchange rates
func Decoder(parcf cc.ParcerFunc, s cc.Storage, country, currencyname, url string, any interface{}) error {
	err := parceReqToStruct(url, any)
	if err != nil {
		cc.Println("!!!!", err)
		return err
	}
	s.Set(parcf, country, currencyname)
	return nil
}

// Run parser
func Run(storage cc.Storage, exitch chan struct{}) {
	requestCurrency(storage)
	go getCurrencyByTime(storage, 30, exitch)

}

// requestCurrency request rates from various banks
func requestCurrency(storage cc.Storage) {
	ru := &RUCurrency{Url: "http://www.cbr.ru/scripts/XML_daily.asp"}
	Decoder(ru.decode, storage, "RU", "RUB", ru.Url, &ru)
	th := THCurrency{Url: "https://www.bot.or.th/App/RSS/fxrate-all.xml"}
	Decoder(th.decode, storage, "TH", "THB", th.Url, &th)
	cc.Println(storage)
}

// GetCurrencyByTime request rates from various banks on time
func getCurrencyByTime(cur cc.Storage, delay time.Duration, exit chan struct{}) {
GETCUR:
	for {
		select {
		case <-time.After(delay * time.Second):
			requestCurrency(cur)
		case <-exit:
			break GETCUR
		}
	}
	cc.Println("Exit from currency requst")
}
