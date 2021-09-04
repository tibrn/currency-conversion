package converter

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/peterhellberg/fixer"
	"github.com/stretchr/testify/require"
)

func TestNewFixer(t *testing.T) {

	req := require.New(t)

	fixer := NewFixer("test")

	req.NotNil(fixer)
}

func TestFixer_Rates(t *testing.T) {

	const (
		accessKey = "test"
	)

	t.Run("Without error", func(t *testing.T) {

		req := require.New(t)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			query := r.URL.Query()

			req.Equal(accessKey, query.Get("access_key"))
			req.Equal("EUR", query.Get("base"))
			req.Equal("RON,USD", query.Get("symbols"))

			byts, err := json.Marshal(map[string]interface{}{
				"base": "EUR",
				"date": "2021-02-13",
				"rates": fixer.Rates{
					"USD": 1.2,
					"RON": 5,
				},
			})

			req.NoError(err)

			w.Write(byts)
		}))

		fixerCli := NewFixer(accessKey)

		client := fixer.NewClient(fixer.AccessKey(accessKey), fixer.BaseURL(server.URL))
		fixerCli.client = client

		data, err := fixerCli.Rates("EUR", []string{"USD", "RON"})

		req.NoError(err)

		req.Contains(data, "USD")
		req.Equal(data["USD"], 1.2)

		req.Contains(data, "RON")
		req.Equal(data["RON"], 5.0)
	})

	t.Run("With error", func(t *testing.T) {
		req := require.New(t)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			query := r.URL.Query()

			req.Equal(accessKey, query.Get("access_key"))
			req.Equal("EUR", query.Get("base"))
			req.Equal("RON,USD", query.Get("symbols"))

			w.WriteHeader(http.StatusInternalServerError)
		}))

		fixerCli := NewFixer(accessKey)

		client := fixer.NewClient(fixer.AccessKey(accessKey), fixer.BaseURL(server.URL))
		fixerCli.client = client

		data, err := fixerCli.Rates("EUR", []string{"USD", "RON"})

		req.Error(err)
		req.Nil(data)

	})
}
