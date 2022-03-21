package handler

import (
	"net/http"
)

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(r.URL.Query().Get("user")))
	w.WriteHeader(http.StatusOK)
}
