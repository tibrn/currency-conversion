package converter

import "currency-conversion/config"

type Converter interface {
	Rates(base string, symbols []string) (map[string]float64, error)
}

var conv Converter

func Get() Converter {
	return conv
}

func init() {

	cfg := config.Get()

	conv = NewFixer(cfg.FixerAccessKey)
}
