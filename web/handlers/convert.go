package handlers

import (
	"currency-conversion/store"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Conversion struct {
	Value  float64 `json:"value"`
	Symbol string  `json:"symbol"`
}

func HandlerConvert(w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("HandlerConvert, error reading body: %v ", err)
		return
	}

	conv := &Conversion{}

	err = json.Unmarshal(data, conv)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("HandlerConvert, error unmarshal: %v ", err)
		return
	}

	value, isSymbol := store.Get().Get(conv.Symbol)

	if !isSymbol {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("HandlerConvert, symbol not found: %s ", conv.Symbol)
		return
	}

	rate, err := strconv.ParseFloat(value, strconv.IntSize)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("HandlerConvert, error conversion: %v ", err)
		return
	}

	symbolConverted := fmt.Sprintf("%f", rate*conv.Value)

	w.Write([]byte(symbolConverted))
}
