package middlewares

import (
	"currency-conversion/config"
	"currency-conversion/store"
	"log"
	"net/http"
	"time"
)

const authHeader = "Authorization"

var (
	cfg                 = config.Get()
	invalidAuthResponse = []byte("Not Authorized!")
)

func refreshAuthorization(expiresAt, value string) {

	expireDate, err := time.Parse(time.RFC3339Nano, expiresAt)

	if err != nil {
		log.Printf("Error parsing expiresAt:%v", err)
		return
	}

	if time.Now().After(expireDate.Add(time.Hour * -24)) {
		go store.Get().Set(
			value,
			time.Now().
				Add(cfg.ExpirationProject).
				Format(time.RFC3339Nano),
			cfg.ExpirationProject,
		)
	}
}

//Authorize.. authorize request
func Authorize(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get(authHeader)

		if expiresAt, isAuth := store.Get().Get(auth); !isAuth {

			w.WriteHeader(http.StatusUnauthorized)
			w.Write(invalidAuthResponse)
			return
		} else {
			refreshAuthorization(expiresAt, auth)
		}

		// next
		next(w, r)
	}

}
