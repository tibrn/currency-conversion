package handlers

import (
	"currency-conversion/config"
	"currency-conversion/helpers"
	"currency-conversion/store"
	"net/http"
	"time"
)

var (
	cfg = config.Get()
)

func HandlerCreateProject(store store.Store) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		authorization, err := helpers.GenerateRandomString(32)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = store.Set(
			authorization,
			time.Now().UTC().
				Format(time.RFC3339),
			cfg.ExpirationProject,
		)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write([]byte(authorization))
	}
}
