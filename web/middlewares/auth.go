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

func refreshAuthorization(createdAt, authorization string) {

	createdDate, err := time.Parse(time.RFC3339Nano, createdAt)

	if err != nil {
		log.Printf("Error parsing createdAt:%v", err)
		return
	}

	//If the request is before 24hours from last update/creation don't update
	if time.Now().Before(createdDate.Add(time.Hour * 24)) {
		return
	}

	//Update authoriation createdAt and ttl
	err = store.Get().Set(
		authorization,
		time.Now().
			Format(time.RFC3339Nano),
		cfg.ExpirationProject,
	)

	if err != nil {
		log.Printf("Error update authorization expire time:%v", err)
	}

}

//Authorize.. authorize request
func Authorize(store store.Store, next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get(authHeader)

		if createdAt, isAuth := store.Get(authorization); !isAuth {

			w.WriteHeader(http.StatusUnauthorized)
			w.Write(invalidAuthResponse)
			return
		} else {
			go refreshAuthorization(createdAt, authorization)
		}

		// next
		next(w, r)
	}

}
