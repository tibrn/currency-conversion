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

func refreshAuthorization(store store.Store, createdAt, authorization string) {

	createdDate, err := time.Parse(time.RFC3339, createdAt)

	if err != nil {
		log.Printf("Error parsing createdAt:%v", err)
		return
	}

	//If the request is before 24hours from last update/creation don't update
	if time.Now().UTC().Before(createdDate.Add(time.Hour * 24)) {
		return
	}

	//Update authoriation createdAt and ttl
	err = store.Set(
		authorization,
		time.Now().UTC().
			Format(time.RFC3339),
		cfg.ExpirationProject,
	)

	if err != nil {
		log.Printf("Error update authorization expire time:%v", err)
	}

}

func unauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write(invalidAuthResponse)
}

//Authorize.. authorize request
func Authorize(store store.Store, next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get(authHeader)

		if len(authorization) == 0 {
			unauthorized(w)
			return
		}

		if createdAt, isAuth := store.Get(authorization); !isAuth {
			unauthorized(w)
			return
		} else {
			go refreshAuthorization(store, createdAt, authorization)
		}

		// next
		next(w, r)
	}

}
