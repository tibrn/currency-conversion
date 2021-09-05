package converter

import (
	"context"
	"net/url"
	"sort"
	"strings"

	"github.com/peterhellberg/fixer"
)

type Fixer struct {
	client *fixer.Client
}

func NewFixer(accessKey string) *Fixer {

	return &Fixer{
		client: fixer.NewClient(fixer.AccessKey(accessKey)),
	}
}

func (f *Fixer) Rates(base string, symbols []string) (map[string]float64, error) {

	sort.Strings(symbols)

	symbolsValue := url.Values{}
	symbolsValue.Add("symbols", strings.Join(symbols, ","))

	resp, err := f.client.Latest(context.Background(), fixer.Base(fixer.Currency(base)), symbolsValue)

	if err != nil {
		return nil, err
	}

	rates := make(map[string]float64)

	for symbol, value := range resp.Rates {
		rates[string(symbol)] = value
	}

	return rates, err
}
