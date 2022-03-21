package handler

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
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
	payload, err := createInvoicePayload(user, comment, intAmt)
	if err != nil {
		json.NewEncoder(w).Encode(&Error{
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

func createInvoicePayload(user, comment string, amount int) (result *CallBack, err error) {
	descriptionToHash := createLnurlMetadata(user)
	hasher := sha256.New()
	_, err = hasher.Write([]byte(descriptionToHash))
	if err != nil {
		return nil, err
	}
	descriptionHash := hex.EncodeToString(hasher.Sum(nil))
	payload := &bytes.Buffer{}
	err = json.NewEncoder(payload).Encode(&LNDHubRequest{
		Amount:          amount,
		Memo:            "comment",
		DescriptionHash: descriptionHash,
	})
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(fmt.Sprintf("https://%s/invoice/%s", LNDHUB_HOST, LNDHUB_LOGIN), "application/json", payload)
	if err != nil {
		return nil, err
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
