package handler

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type CallBack struct {
	// payment request
	Pr string `json:"pr,omitempty"`
	// routes, unused
	Routes []interface{} `json:"routes"`
}
type Error struct {
	// reason
	Reason string `json:"reason,omitempty"`
	// status
	Status string `json:"status,omitempty"`
}
type LNDHubResponse struct {
	PayReq         string `json:"pay_req"`
	PaymentRequest string `json:"payment_request"`
	RHash          string `json:"r_hash"`
}
type LNDHubRequest struct {
	Amount          int    `json:"amt"` // amount in Satoshi
	Memo            string `json:"memo"`
	DescriptionHash string `json:"description_hash" validate:"omitempty,hexadecimal,len=64"`
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
	payload, err := createInvoicePayload(user, r.Host, comment, intAmt)
	if err != nil {
		_ = json.NewEncoder(w).Encode(&Error{
			Reason: err.Error(),
			Status: "ERROR",
		})

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(payload)
	if err != nil {
		fmt.Println(err)
	}
	w.WriteHeader(http.StatusOK)
}

func createInvoicePayload(user, host, comment string, amountMsat int) (result *CallBack, err error) {
	descriptionToHash := createLnurlMetadata(user, host)
	hasher := sha256.New()
	_, err = hasher.Write([]byte(descriptionToHash))
	if err != nil {
		return nil, err
	}
	descriptionHash := hasher.Sum(nil)
	payload := &bytes.Buffer{}
	err = json.NewEncoder(payload).Encode(&LNDHubRequest{
		Amount:          amountMsat / 1000,
		Memo:            comment,
		DescriptionHash: hex.EncodeToString([]byte(descriptionHash)),
	})
	if err != nil {
		return nil, err
	}
	lndhubHost := os.Getenv("LNDHUB_HOST")
	lndhubLogin := os.Getenv("LNDHUB_LOGIN")
	if lndhubHost == "" || lndhubLogin == "" {
		return nil, fmt.Errorf("LNURL Pay Server is not configured correctly, contact admin.")
	}

	resp, err := http.Post(fmt.Sprintf("https://%s/invoice/%s", lndhubHost, lndhubLogin), "application/json", payload)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Something went wrong calling lndhub, status code %d", resp.StatusCode)
	}
	invoice := &LNDHubResponse{}
	err = json.NewDecoder(resp.Body).Decode(invoice)
	if err != nil {
		return nil, err
	}
	return &CallBack{
		Pr: invoice.PaymentRequest,
	}, err
}
