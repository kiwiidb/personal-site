package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

const (
	DEFAULT_MIN = int(1e3)
	DEFAULT_MAX = int(1e8)
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
	MIN := DEFAULT_MIN
	MAX := DEFAULT_MAX
	minimum := os.Getenv("MIN_SENDABLE")
	var err error
	if minimum != "" {
		MIN, err = strconv.Atoi(minimum)
		if err != nil {
			MIN = DEFAULT_MIN
		}
	}
	maximum := os.Getenv("MAX_SENDABLE")
	if maximum != "" {
		MAX, err = strconv.Atoi(maximum)
		if err != nil {
			MAX = DEFAULT_MAX
		}
	}
	response := &LNURLPayBody{
		Callback:       fmt.Sprintf("https://%s/api/callback?user=%s", r.Host, user),
		CommentAllowed: 512,
		MaxSendable:    int32(MAX),
		Metadata:       createLnurlMetadata(user, r.Host),
		MinSendable:    int32(MIN),
		Tag:            "payRequest",
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func createLnurlMetadata(username, host string) (result string) {
	return fmt.Sprintf("[[\"text/plain\", \"Payment to %s@%s\"]]", username, host)

}
