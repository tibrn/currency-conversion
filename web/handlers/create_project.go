package handlers

import (
	"currency-conversion/helpers"
	"net/http"
)

func HandlerCreateProject(w http.ResponseWriter, r *http.Request) {

	authorization, err := helpers.GenerateRandomString(32)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(authorization))
}
