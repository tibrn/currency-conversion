package job

import (
	"currency-conversion/converter"
	"currency-conversion/store"
)

func Start() {

	currencies := Currencies{
		"EUR": []string{"USD"},
		"USD": []string{"EUR"},
	}

	startJob(currencies, store.Get(), converter.Get())
}
