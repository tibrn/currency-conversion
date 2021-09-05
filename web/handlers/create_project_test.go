package handlers

import (
	"currency-conversion/helpers"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bou.ke/monkey"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestHandlerCreateProject(t *testing.T) {

	req := require.New(t)
	const (
		randomString = "test"
	)
	monkey.Patch(time.Now, func() time.Time {
		date, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
		req.NoError(err)

		return date.UTC()
	})

	monkey.Patch(helpers.GenerateRandomString, func(n int) (string, error) {
		return randomString, nil
	})

	t.Run("Without error", func(t *testing.T) {
		req := require.New(t)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := NewMockStore(ctrl)

		store.
			EXPECT().
			Set("test", time.Now().UTC().Format(time.RFC3339), cfg.ExpirationProject).
			Return(nil)

		createProject := HandlerCreateProject(store)

		rec := httptest.NewRecorder()

		r, err := http.NewRequest(http.MethodPost, "/create", nil)
		req.NoError(err)
		createProject(rec, r)

		req.Equal(http.StatusOK, rec.Result().StatusCode)

		body, err := ioutil.ReadAll(rec.Body)
		req.NoError(err)
		req.Equal(randomString, string(body))
	})

	t.Run("With error", func(t *testing.T) {
		req := require.New(t)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := NewMockStore(ctrl)

		store.
			EXPECT().
			Set("test", time.Now().UTC().Format(time.RFC3339), cfg.ExpirationProject).
			Return(errors.New("fake"))

		createProject := HandlerCreateProject(store)

		rec := httptest.NewRecorder()

		r, err := http.NewRequest(http.MethodPost, "/create", nil)
		req.NoError(err)
		createProject(rec, r)

		req.Equal(http.StatusInternalServerError, rec.Result().StatusCode)

	})
}
