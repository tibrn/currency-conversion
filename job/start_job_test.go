package job

import (
	"currency-conversion/converter"
	"currency-conversion/store"
	"errors"
	"fmt"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func expected(t *testing.T, isErr bool) (Currencies, store.Store, converter.Converter, func()) {

	ctrl := gomock.NewController(t)

	store := NewMockStore(ctrl)
	converter := NewMockConverter(ctrl)

	currencies := Currencies{
		"EUR": []string{"USD", "RON"},
		"USD": []string{"EUR"},
	}

	var err error

	if isErr {
		err = errors.New("fake")
	}

	converter.
		EXPECT().
		Rates("EUR", []string{"USD", "RON"}).
		Return(map[string]float64{
			"USD": 1.2,
			"RON": 5.0,
		}, nil)

	converter.
		EXPECT().
		Rates("USD", []string{"EUR"}).
		Return(map[string]float64{
			"EUR": 0.8,
		}, nil)

	store.
		EXPECT().
		Set(fmt.Sprintf(formatExchange, "EUR", "USD"), fmt.Sprintf("%f", 1.2), time.Hour*2).
		Return(err)

	if !isErr {
		store.
			EXPECT().
			Set(fmt.Sprintf(formatExchange, "EUR", "RON"), fmt.Sprintf("%f", 5.0), time.Hour*2).
			Return(err)
	}

	store.
		EXPECT().
		Set(fmt.Sprintf(formatExchange, "USD", "EUR"), fmt.Sprintf("%f", 0.8), time.Hour*2).
		Return(nil)

	return currencies, store, converter, func() {
		ctrl.Finish()
	}
}

func Test_updateCurrencies(t *testing.T) {

	t.Run("Without error", func(t *testing.T) {
		req := require.New(t)
		currencies, store, converter, finish := expected(t, false)
		defer finish()
		err := updateCurrencies(currencies, store, converter)

		req.NoError(err)
	})

	t.Run("With error", func(t *testing.T) {
		req := require.New(t)
		currencies, store, converter, finish := expected(t, true)
		defer finish()
		err := updateCurrencies(currencies, store, converter)

		req.Error(err)
	})
}

func Test_startJob(t *testing.T) {

	currencies, store, converter, finish := expected(t, false)
	defer finish()

	go startJob(currencies, store, converter)

	time.Sleep(time.Millisecond * 50)

}
