package handlers

import "net/http"

func hasError(w http.ResponseWriter, err error) bool {
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return true
	}
	return false
}
