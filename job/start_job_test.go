package job

import (
	"currency-conversion/converter"
	"currency-conversion/store"
	"fmt"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func expected(t *testing.T) (Currencies, store.Store, converter.Converter, func()) {

	ctrl := gomock.NewController(t)

	store := NewMockStore(ctrl)
	converter := NewMockConverter(ctrl)

	currencies := Currencies{
		"EUR": []string{"USD", "RON"},
		"USD": []string{"EUR"},
	}

	converter.
		EXPECT().
		Rates("EUR", []string{"USD", "RON"}).
		Return(map[string]float64{
			"USD": 1.2,
			"RON": 5.0,
		})

	converter.
		EXPECT().
		Rates("USD", []string{"EUR"}).
		Return(map[string]float64{
			"EUR": 0.8,
		})

	store.
		EXPECT().
		Set(fmt.Sprintf(formatExchange, "EUR", "USD"), fmt.Sprintf("%f", 1.2), time.Hour*2)

	store.
		EXPECT().
		Set(fmt.Sprintf(formatExchange, "EUR", "RON"), fmt.Sprintf("%f", 5.0), time.Hour*2)

	store.
		EXPECT().
		Set(fmt.Sprintf(formatExchange, "USD", "EUR"), fmt.Sprintf("%f", 0.8), time.Hour*2)

	return currencies, store, converter, func() {
		ctrl.Finish()
	}
}
func Test_updateCurrencies(t *testing.T) {
	req := require.New(t)

	currencies, store, converter, finish := expected(t)
	defer finish()
	err := updateCurrencies(currencies, store, converter)

	req.NoError(err)
}

func Test_startJob(t *testing.T) {

	currencies, store, converter, finish := expected(t)
	defer finish()

	go startJob(currencies, store, converter)

	time.Sleep(time.Millisecond * 50)

}
