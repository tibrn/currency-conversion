package cmd

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_makeRequest(t *testing.T) {

	req := require.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req.Equal("/test", r.URL.Path)
		req.Equal(http.MethodPost, r.Method)

		data, err := ioutil.ReadAll(r.Body)
		req.NoError(err)

		byts, err := json.Marshal("test")
		req.NoError(err)
		req.Equal(string(byts), string(data))

		w.Write([]byte("test2"))
	}))

	resp, err := makeRequest(srv.URL, http.MethodPost, "test", "test")

	req.NoError(err)

	req.Equal("test2", resp)
}
