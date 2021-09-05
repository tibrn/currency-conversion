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

	mtx := sync.Mutex{}

	updateError := func(newErr error) {
		mtx.Lock()
		defer mtx.Unlock()

		if newErr == nil {
			return
		}

		if err == nil {
			err = newErr
		} else {
			err = fmt.Errorf("%v&%v", err, newErr)
		}
	}

	updateCurrency := func(base string, symbols []string) error {

		rates, err := converter.Rates(base, symbols)

		if err != nil {
			return err
		}

		for symbol, value := range rates {
			err := store.Set(
				fmt.Sprintf(formatExchange, base, symbol),
				fmt.Sprintf("%f", value),
				time.Hour*2,
			)

			log.Printf("Symbol:%s rate:%f", fmt.Sprintf(formatExchange, base, symbol), value)

			if err != nil {
				return err
			}

		}

		return nil
	}

	var (
		group      = sync.WaitGroup{}
		count      = 0
		maxPending = 25
	)

	for base, symbols := range currencies {

		log.Printf("%v , %v\n", base, symbols)
		//Increase
		group.Add(1)
		count++

		go func(base string, symbols []string) {
			defer group.Done()

			err := updateCurrency(base, symbols)

			log.Printf("Error updateCurrency: %v\n", err)

			updateError(err)

		}(base, symbols)

		//Wait to finish if we've reached the max limit of requests
		if count == maxPending {
			group.Wait()

			if err != nil {
				return
			}

			count = 0
		}

	}

	group.Wait()

	log.Printf("Job update done!\n")

	return err

}
