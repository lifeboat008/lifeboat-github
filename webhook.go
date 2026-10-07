package github

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	protocol "github.com/lifeboat008/lifeboat-protocol"
)

type ProjectBinding struct {
	ProjectID      string `json:"project_id"`
	InstallationID int64  `json:"installation_id"`
}

type Adapter struct {
	Secret     []byte
	APIToken   string
	APIURL     string
	Projects   map[string]ProjectBinding
	HTTPClient *http.Client
}

func (a *Adapter) Validate() error {
	if len(a.Secret) < 24 || len(a.APIToken) < 24 || len(a.Projects) == 0 {
		return errors.New("webhook secret, API token, and project bindings are required")
	}
	u, err := url.Parse(a.APIURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return errors.New("valid API URL is required")
	}
	for repo, binding := range a.Projects {
		if !strings.Contains(repo, "/") || binding.ProjectID == "" || binding.InstallationID <= 0 {
			return errors.New("invalid project binding")
		}
	}
	return nil
}

type event struct {
	Action       string `json:"action"`
	Installation struct {
		ID int64 `json:"id"`
	} `json:"installation"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
	PullRequest struct {
		Merged  bool   `json:"merged"`
		HTMLURL string `json:"html_url"`
	} `json:"pull_request"`
	Release struct {
		HTMLURL string `json:"html_url"`
	} `json:"release"`
	Issue struct {
		HTMLURL string `json:"html_url"`
	} `json:"issue"`
}

func (a *Adapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}
	if err := a.Validate(); err != nil {
		http.Error(w, "adapter is not configured", 500)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" && !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json;") {
		http.Error(w, "JSON required", 415)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid body", 400)
		return
	}
	signature := r.Header.Get("X-Hub-Signature-256")
	if !strings.HasPrefix(signature, "sha256=") {
		http.Error(w, "invalid signature", 401)
		return
	}
	actual, err := hex.DecodeString(strings.TrimPrefix(signature, "sha256="))
	if err != nil {
		http.Error(w, "invalid signature", 401)
		return
	}
	mac := hmac.New(sha256.New, a.Secret)
	_, _ = mac.Write(body)
	if !hmac.Equal(mac.Sum(nil), actual) {
		http.Error(w, "invalid signature", 401)
		return
	}
	delivery := r.Header.Get("X-GitHub-Delivery")
	if len(delivery) < 8 || len(delivery) > 128 {
		http.Error(w, "invalid delivery ID", 400)
		return
	}
	var payload event
	if json.Unmarshal(body, &payload) != nil {
		http.Error(w, "invalid event", 400)
		return
	}
	binding, ok := a.Projects[payload.Repository.FullName]
	if !ok || binding.InstallationID != payload.Installation.ID {
		http.Error(w, "repository installation is not enrolled", 403)
		return
	}
	var kind, link string
	switch r.Header.Get("X-GitHub-Event") {
	case "pull_request":
		if payload.Action == "closed" && payload.PullRequest.Merged {
			kind = "merged_pr"
			link = payload.PullRequest.HTMLURL
		}
	case "release":
		if payload.Action == "published" {
			kind = "release"
			link = payload.Release.HTMLURL
		}
	case "issues":
		if payload.Action == "closed" {
			kind = "closed_issue"
			link = payload.Issue.HTMLURL
		}
	}
	if kind == "" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if !strings.HasPrefix(link, "https://github.com/"+payload.Repository.FullName+"/") {
		http.Error(w, "invalid evidence URL", 400)
		return
	}
	hash := sha256.Sum256([]byte(delivery))
	evidence := protocol.Evidence{ID: "ev_" + hex.EncodeToString(hash[:16]), ProjectID: binding.ProjectID, InstallationID: binding.InstallationID, DeliveryID: delivery, Kind: kind, URL: link, ObservedAt: time.Now().UTC()}
	data, _ := json.Marshal(evidence)
	client := a.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}
	req, err := http.NewRequestWithContext(r.Context(), "POST", strings.TrimRight(a.APIURL, "/")+"/v1/evidence", bytes.NewReader(data))
	if err != nil {
		http.Error(w, "API request failed", 502)
		return
	}
	req.Header.Set("Authorization", "Bearer "+a.APIToken)
	req.Header.Set("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		http.Error(w, "API unavailable", 502)
		return
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusConflict {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"duplicate"}`))
		return
	}
	if response.StatusCode != http.StatusCreated {
		http.Error(w, "API rejected evidence", 502)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(`{"status":"recorded"}`))
}
