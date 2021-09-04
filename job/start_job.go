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
	cfg = config.Get()
)

const formatExchange = "%s/%s"

func startJob(currencies Currencies, store store.Store, converter converter.Converter) {
	for {

		err := retry.Do(func() error {
			return updateCurrencies(currencies, store, converter)
		})

		if err != nil {
			log.Printf("Error updateCurrencies:%v\n", err)
		}

		time.Sleep(time.Hour)
	}
}

func updateCurrencies(currencies Currencies, store store.Store, converter converter.Converter) (err error) {

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
			errSet := store.Set(
				fmt.Sprintf(formatExchange, base, symbol),
				fmt.Sprintf("%v", value),
				time.Hour*2,
			)

			if errSet != nil {
				//Set return error
				err = errSet
				return
			}
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
