package pastee

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"appengine"
	"appengine/datastore"
)

type PastesPostReq struct {
	Content string
	Mac     string
	Expiry  string
}

type PastesPostResp struct {
	Id string `json:"id"`
}

type Paste struct {
	Content string
	Mac     string
	Expiry  time.Time
}

func init() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/pastes/", pastesGetHandler)
	http.HandleFunc("/pastes", pastesPostHandler)
}

func indexHandler(w http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(w, "Index handler")
}

// Handles GET requests to /pastes/{id}.
func pastesGetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Extract id from URL.
	id := strings.Replace(r.URL.Path, "/pastes/", "", -1)

	fmt.Fprintf(w, "%v", id)
}

func pastesPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	const kMaxBodyLength = 256 * 1024
	if r.ContentLength > kMaxBodyLength {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}
	if r.ContentLength < 0 {
		w.WriteHeader(http.StatusLengthRequired)
		return
	}

	postData := make([]byte, r.ContentLength)
	_, err := r.Body.Read(postData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "POST body required")
		return
	}

	var request PastesPostReq
	err = json.Unmarshal(postData, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%v\n", err)
		return
	}

	ctx := appengine.NewContext(r)
	code, response, err := pastesPostRPC(&ctx, &request)
	w.WriteHeader(code)
	if err != nil {
		fmt.Fprintf(w, "error: %+v\n", err)
	}

	responseBytes, err := json.Marshal(response)
	fmt.Fprintf(w, "request: %+v\n", request)
	fmt.Fprintf(w, "response: %+v\n", string(responseBytes))
}

func pastesPostRPC(ctx *appengine.Context, request *PastesPostReq) (int, PastesPostResp, error) {
	// TODO(ms): These should be configurable.
	const kMaxContentLength = 256 * 1024 // 256 KiB
	const kMaxMacLength = 128
	const kMaxLifetime = 7 * 24 * time.Hour

	// Length checking.
	if request.Content == "" {
		return http.StatusBadRequest, PastesPostResp{}, errors.New("content is required")
	} else if len(request.Content) > kMaxContentLength {
		return http.StatusBadRequest, PastesPostResp{}, errors.New(
			fmt.Sprintf("max content is %d bytes", kMaxContentLength))
	} else if len(request.Mac) > kMaxMacLength {
		return http.StatusBadRequest, PastesPostResp{}, errors.New(
			fmt.Sprintf("max mac length is %d bytes", kMaxMacLength))
	}

	// Parse and validate expiration date.
	var expiry time.Time
	if request.Expiry != "" {
		var err error
		expiry, err = time.Parse(time.RFC3339, request.Expiry)
		if err != nil {
			return http.StatusBadRequest, PastesPostResp{}, errors.New("bad time format")
		}
	}
	now := time.Now()
	if expiry.After(now.Add(kMaxLifetime)) {
		return http.StatusBadRequest, PastesPostResp{}, errors.New(
			fmt.Sprintf("maximum lifetime is %v", kMaxLifetime))
	}

	// Construct Paste entity for datastore.
	var paste Paste
	paste.Content = request.Content
	paste.Mac = request.Mac
	paste.Expiry = expiry

	fmt.Fprintf(os.Stderr, "Paste: %+v\n", paste)

	// Insert Paste.
	key, err := datastore.Put(
		*ctx, datastore.NewIncompleteKey(*ctx, "paste", nil), &paste)
	if err != nil {
		return http.StatusInternalServerError, PastesPostResp{}, err
	}

	// Paste created successfully.
	var response PastesPostResp
	response.Id = MBase31{Value: key.IntID()}.ToString()
	return http.StatusCreated, response, nil
}
