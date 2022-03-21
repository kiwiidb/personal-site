package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type CallBack struct {
	// payment request
	Pr string `json:"pr,omitempty"`
	// routes, unused
	Routes []interface{} `json:"routes"`
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := r.URL.Query().Get("user")
	comment := r.URL.Query().Get("comment")
	amount := r.URL.Query().Get("amount")
	if user == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	intAmt, err := strconv.Atoi(amount)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	payload, err := createInvoicePayload(user, comment, intAmt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(payload)
	if err != nil {
		fmt.Println(err)
	}
	w.WriteHeader(http.StatusOK)
}

func createInvoicePayload(user, comments string, amount int) (result *CallBack, err error) {
	return &CallBack{
		Pr: "lnbctest123",
	}, err
}
