package github

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignedMergedPRRecordedOnce(t *testing.T) {
	calls := 0
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/evidence" || r.Header.Get("Authorization") != "Bearer a-strong-api-token-123456789" {
			t.Errorf("incorrect API request")
			w.WriteHeader(400)
			return
		}
		calls++
		if calls == 1 {
			w.WriteHeader(201)
		} else {
			w.WriteHeader(409)
		}
	}))
	defer api.Close()
	adapter := &Adapter{Secret: []byte("a-strong-webhook-secret-123456789"), APIToken: "a-strong-api-token-123456789", APIURL: api.URL, Projects: map[string]ProjectBinding{"owner/repo": {ProjectID: "project1", InstallationID: 42}}}
	body := []byte(`{"action":"closed","installation":{"id":42},"repository":{"full_name":"owner/repo"},"pull_request":{"merged":true,"html_url":"https://github.com/owner/repo/pull/2"}}`)
	call := func(valid bool) int {
		req := httptest.NewRequest("POST", "/webhook", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-GitHub-Event", "pull_request")
		req.Header.Set("X-GitHub-Delivery", "delivery-123")
		mac := hmac.New(sha256.New, adapter.Secret)
		_, _ = mac.Write(body)
		sig := hex.EncodeToString(mac.Sum(nil))
		if !valid {
			sig = "00"
		}
		req.Header.Set("X-Hub-Signature-256", "sha256="+sig)
		response := httptest.NewRecorder()
		adapter.ServeHTTP(response, req)
		return response.Code
	}
	if got := call(false); got != 401 {
		t.Fatalf("invalid signature: %d", got)
	}
	if got := call(true); got != 201 {
		t.Fatalf("first delivery: %d", got)
	}
	if got := call(true); got != 200 {
		t.Fatalf("duplicate delivery: %d", got)
	}
	if calls != 2 {
		t.Fatalf("API calls: %d", calls)
	}
}
