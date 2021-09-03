package job

import (
	"currency-conversion/converter"
	"currency-conversion/store"
)

func Start() {
	startJob(store.Get(), converter.Get())
}
