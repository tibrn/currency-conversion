package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bou.ke/monkey"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestAuthorize(t *testing.T) {

	t.Run("Missing authorization", func(t *testing.T) {

		req := require.New(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := NewMockStore(ctrl)

		isCalled := false
		authorizer := Authorize(store, func(w http.ResponseWriter, r *http.Request) {
			isCalled = true
		})

		rec := httptest.NewRecorder()

		r, err := http.NewRequest(http.MethodPost, "/test", nil)
		req.NoError(err)

		authorizer(rec, r)
		req.Equal(false, isCalled)
	})

	t.Run("Invalid authorization", func(t *testing.T) {

		req := require.New(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := NewMockStore(ctrl)

		store.
			EXPECT().
			Get("invalid_auth").
			Return(time.Now().Format(time.RFC3339), false)

		isCalled := false
		authorizer := Authorize(store, func(w http.ResponseWriter, r *http.Request) {
			isCalled = true
		})

		rec := httptest.NewRecorder()

		r, err := http.NewRequest(http.MethodPost, "/test", nil)
		r.Header.Add("Authorization", "invalid_auth")
		req.NoError(err)

		authorizer(rec, r)
		req.Equal(false, isCalled)
	})

	t.Run("Next called", func(t *testing.T) {

		req := require.New(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := NewMockStore(ctrl)

		now := time.Now().UTC().Format(time.RFC3339)
		store.
			EXPECT().
			Get("next_called").
			Return(now, true)

		isCalled := false
		authorizer := Authorize(store, func(rw http.ResponseWriter, r *http.Request) {
			isCalled = true
		})

		rec := httptest.NewRecorder()

		r, err := http.NewRequest(http.MethodPost, "/test", nil)
		r.Header.Add("Authorization", "next_called")
		req.NoError(err)

		authorizer(rec, r)
		req.Equal(true, isCalled)
	})

	t.Run("Update authorization ttl", func(t *testing.T) {
		req := require.New(t)

		monkey.Patch(time.Now, func() time.Time {
			date, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
			req.NoError(err)

			return date
		})

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := NewMockStore(ctrl)

		store.
			EXPECT().
			Get("update_auth").
			Return(time.Now().UTC().Add(time.Hour*-25).Format(time.RFC3339), true)

		store.
			EXPECT().
			Set("update_auth", time.Now().UTC().Format(time.RFC3339), cfg.ExpirationProject)

		isCalled := false
		authorizer := Authorize(store, func(rw http.ResponseWriter, r *http.Request) {
			isCalled = true
		})

		rec := httptest.NewRecorder()

		r, err := http.NewRequest(http.MethodPost, "/test", nil)
		r.Header.Add("Authorization", "update_auth")
		req.NoError(err)

		authorizer(rec, r)

		req.Equal(true, isCalled)

		time.Sleep(time.Millisecond * 100)

	})

}

func Test_refreshAuthorization(t *testing.T) {

	req := require.New(t)

	monkey.Patch(time.Now, func() time.Time {
		date, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
		req.NoError(err)

		return date
	})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := NewMockStore(ctrl)

	store.
		EXPECT().
		Set("auth_refresh", time.Now().UTC().Format(time.RFC3339), cfg.ExpirationProject)

	refreshAuthorization(store, "test", "auth_refresh")

	refreshAuthorization(store, time.Now().UTC().Add(time.Hour*-24+time.Minute).Format(time.RFC3339), "auth_refresh")

	refreshAuthorization(store, time.Now().UTC().Add(time.Hour*-24-time.Second).Format(time.RFC3339), "auth_refresh")

}
