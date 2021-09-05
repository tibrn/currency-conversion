package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func Test_convert(t *testing.T) {
	req := require.New(t)

	isCalled := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		req.Equal("test", r.Header.Get("Authorization"))
		isCalled = true

		data := map[string]interface{}{}

		body, err := ioutil.ReadAll(r.Body)

		req.NoError(err)

		err = json.Unmarshal(body, &data)
		req.NoError(err)

		req.Contains(data, "symbol")
		req.Equal("EUR/USD", data["symbol"])

		req.Contains(data, "value")
		req.Equal(5.0, data["value"])

		w.Write([]byte(fmt.Sprintf("%f", 5.0*2.5)))

	}))

	apiKey = "test"

	command := convert(func() (string, bool) {
		return srv.URL, true
	})

	value = "5.0"
	symbol = "EUR/USD"

	command(&cobra.Command{}, []string{})

	req.Equal(srv.URL, viper.GetString(viperHost))
	req.Equal(true, isCalled)
}
