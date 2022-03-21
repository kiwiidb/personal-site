package handler

import (
	"net/http"
)

func LnUrlPHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(r.URL.Query().Get("user")))
	w.WriteHeader(http.StatusOK)
}
