package cmd

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func Test_create(t *testing.T) {

	req := require.New(t)

	isCalled := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isCalled = true

		w.Write([]byte("test"))
	}))

	command := create(func() (string, bool) {
		return srv.URL, true
	})

	command(&cobra.Command{}, []string{})

	req.Equal(srv.URL, viper.GetString(viperHost))
	req.Equal("test", viper.GetString(viperApiKey))
	req.Equal(true, isCalled)
}
