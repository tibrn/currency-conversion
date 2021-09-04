package job

import (
	"currency-conversion/config"
	"currency-conversion/converter"
	"currency-conversion/store"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/avast/retry-go"
)

type Currencies map[string][]string

var (
	currencies = Currencies{}
	cfg        = config.Get()
)

const formatExchange = "%s/%s"

func startJob(store store.Store, converter converter.Converter) {
	for {

		err := retry.Do(func() error {
			return updateCurrencies(store, converter)
		})

		if err != nil {
			log.Printf("Error updateCurrencies:%v\n", err)
		}

		time.Sleep(time.Hour)
	}
}

func updateCurrencies(store store.Store, converter converter.Converter) (err error) {

	defer func() {
		if errRec := recover(); errRec != nil {
			err = fmt.Errorf("%v", errRec)
		}
	}()

	var (
		group      = sync.WaitGroup{}
		count      = 0
		maxPending = 25
	)

	updateCurrency := func(base string, symbols []string) {
		defer group.Done()

		rates, errConv := converter.Rates(base, symbols)

		if errConv != nil {
			//Set return error
			err = errConv
			return
		}

		for symbol, value := range rates {
			store.Set(fmt.Sprintf(formatExchange, base, symbol), fmt.Sprintf("%v", value), cfg.ExpirationProject)
		}
	}

	for base, symbols := range currencies {

		//Increase
		group.Add(1)
		count++

		go updateCurrency(base, symbols)

		//Wait to finish if we've reached the max limit of requests
		if count == maxPending {
			group.Wait()
		}

		if err != nil {
			return
		}
	}

	group.Wait()

	return nil

}
