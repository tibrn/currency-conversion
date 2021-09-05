package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestHandlerConvert(t *testing.T) {
	t.Run("Without error", func(t *testing.T) {
		req := require.New(t)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := NewMockStore(ctrl)

		store.
			EXPECT().
			Get("test").
			Return("0.5", true)

		conv := Conversion{
			Symbol: "test",
			Value:  5,
		}

		byts, err := json.Marshal(conv)
		req.NoError(err)

		convertSymbol := HandlerConvert(store)

		rec := httptest.NewRecorder()

		r, err := http.NewRequest(http.MethodPost, "/convert", bytes.NewBuffer(byts))
		req.NoError(err)
		convertSymbol(rec, r)

		req.Equal(http.StatusOK, rec.Result().StatusCode)

		body, err := ioutil.ReadAll(rec.Body)
		req.NoError(err)
		req.Equal(fmt.Sprintf("%f", 5*0.5), string(body))
	})

	t.Run("Symbol not found", func(t *testing.T) {
		req := require.New(t)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := NewMockStore(ctrl)

		store.
			EXPECT().
			Get("test").
			Return("", false)

		conv := Conversion{
			Symbol: "test",
		}

		byts, err := json.Marshal(conv)
		req.NoError(err)

		convertSymbol := HandlerConvert(store)

		rec := httptest.NewRecorder()

		r, err := http.NewRequest(http.MethodPost, "/convert", bytes.NewBuffer(byts))
		req.NoError(err)
		convertSymbol(rec, r)

		req.Equal(http.StatusNotFound, rec.Result().StatusCode)

	})

	t.Run("Invalid input", func(t *testing.T) {
		req := require.New(t)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := NewMockStore(ctrl)

		convertSymbol := HandlerConvert(store)

		rec := httptest.NewRecorder()

		r, err := http.NewRequest(http.MethodPost, "/convert", bytes.NewBuffer([]byte{}))
		req.NoError(err)
		convertSymbol(rec, r)

		req.Equal(http.StatusInternalServerError, rec.Result().StatusCode)

	})

	t.Run("Invalid rate conversion for symbol", func(t *testing.T) {
		req := require.New(t)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := NewMockStore(ctrl)

		store.
			EXPECT().
			Get("test").
			Return("abasd", true)

		conv := Conversion{
			Symbol: "test",
		}

		byts, err := json.Marshal(conv)
		req.NoError(err)

		convertSymbol := HandlerConvert(store)

		rec := httptest.NewRecorder()

		r, err := http.NewRequest(http.MethodPost, "/convert", bytes.NewBuffer(byts))
		req.NoError(err)
		convertSymbol(rec, r)

		req.Equal(http.StatusInternalServerError, rec.Result().StatusCode)

	})
}
