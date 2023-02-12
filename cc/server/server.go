package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"currconv/cc"

	"github.com/gorilla/mux"
)

// HttpServer structure for working with the server and transferring data for processing
type HttpServer struct {
	St  cc.Storage
	srv *http.Server
}

// RequestForConvert data type intended for an incoming request aimed at obtaining the amount of the converted currency
type RequestForConvert struct {
	Form    string  `json:"from"`
	To      string  `json:"to"`
	Country string  `jsont:"country"`
	Amount  float64 `json:"amount"`
}

// ResponseAfterConvert data type intended to respond to a currency conversion request
type ResponseAfterConvert struct {
	Amount float64 `json:"amout"`
}

type ResponseError struct {
	Error string `json:"error"`
}

// respError writer bed response
func respError(errTxt string) []byte {
	err := ResponseError{Error: errTxt}
	errB, _ := json.Marshal(err)
	return errB
}

// respOk writer good response
func respOk(w http.ResponseWriter, data interface{}) {
	b, err := json.Marshal(data)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(respError("error parcing to json"))
	} else {
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(b)
		if err != nil {
			cc.Println(err)
		}
	}
}

// respErr returns an error to the client
func respErr(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// listCurrency returns a list of exchange rates for the specified central bank
func (h *HttpServer) listCurrency(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	country := strings.ToUpper(vars["country"])
	res, err := h.St.GetCurrencyCountry(country)
	if err != nil {
		respErr(w, respError("Wrong country index"))
	} else {
		respOk(w, res)
	}
}

// currencyConverter performs currency conversionperforms currency conversion.
func (h *HttpServer) currencyConverter(w http.ResponseWriter, r *http.Request) {
	if err := checkValidRequest(r); err != nil {
		respErr(w, respError(err.Error()))
		return
	}
	var req RequestForConvert
	cc.Println(r.Body)
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respErr(w, respError(err.Error()))
		cc.Println(err)
	} else {
		val, err := h.convert(strings.ToUpper(countrySelector(req.Country)), strings.ToUpper(strings.TrimSpace(req.Form)), strings.ToUpper(strings.TrimSpace(req.To)))
		if err != nil {
			respErr(w, respError(err.Error()))
		} else {
			res := countAmount(val, req.Amount)
			respOk(w, res)
		}
	}

}

// checkvalidrequest
func checkValidRequest(r *http.Request) error {
	return nil
}

// HttpApi starts the web server
func (h *HttpServer) HttpApi() {
	r := mux.NewRouter()
	h.srv = &http.Server{Addr: ":8080", Handler: r}
	r.HandleFunc("/listcurrency/{country}", h.listCurrency).Methods("GET")
	r.HandleFunc("/currencyconverter", h.currencyConverter).Methods("POST")
	cc.Println("Starting server on :8080")
	go func() {
		if err := h.srv.ListenAndServe(); err != nil {
			cc.Println("listen:", err)
		}
	}()

}

// HttpApiStop stop webserver
func (h *HttpServer) HttpApiStop() {
	if err := h.srv.Shutdown(context.Background()); err != nil {
		fmt.Printf("HTTP server Shutdown: %v\n", err)
	}

	cc.Println("http server stoped")
}
