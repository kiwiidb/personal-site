package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	HOST         = "kwintendebacker.com"
	MAX          = 1e6
	MIN          = 1e4
	LNDHUB_HOST  = "clnhub.mainnet.getalby.com"
	LNDHUB_LOGIN = "ZnznCgl3VxV5hMFF9MVT"
)

type LNURLPayBody struct {
	// callback
	Callback string `json:"callback,omitempty"`
	// comment allowed
	CommentAllowed int32 `json:"commentAllowed,omitempty"`
	// max sendable
	MaxSendable int32 `json:"maxSendable,omitempty"`
	// metadata
	Metadata string `json:"metadata,omitempty"`
	// min sendable
	MinSendable int32 `json:"minSendable,omitempty"`
	// tag
	Tag string `json:"tag,omitempty"`
}

func LnUrlPHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := r.URL.Query().Get("user")
	if user == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	response := &LNURLPayBody{
		Callback:       fmt.Sprintf("https://%s/api/callback?user=%s", HOST, user),
		CommentAllowed: 512,
		MaxSendable:    MAX,
		Metadata:       createLnurlMetadata(user),
		MinSendable:    MIN,
		Tag:            "payRequest",
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func createLnurlMetadata(username string) (result string) {
	return fmt.Sprintf("[[\"text/plain\", \"Payment to %s@%s\"]]", username, HOST)

}
