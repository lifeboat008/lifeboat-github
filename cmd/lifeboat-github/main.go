package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	github "github.com/lifeboat008/lifeboat-github"
)

func main() {
	adapter := &github.Adapter{
		Secret:   []byte(os.Getenv("LIFEBOAT_GITHUB_WEBHOOK_SECRET")),
		APIToken: os.Getenv("LIFEBOAT_GITHUB_API_TOKEN"),
		APIURL:   os.Getenv("LIFEBOAT_API_URL"),
	}
	if err := json.Unmarshal([]byte(os.Getenv("LIFEBOAT_GITHUB_PROJECTS_JSON")), &adapter.Projects); err != nil {
		log.Fatal("LIFEBOAT_GITHUB_PROJECTS_JSON must map owner/repo to project and installation IDs")
	}
	if err := adapter.Validate(); err != nil {
		log.Fatal(err)
	}
	address := os.Getenv("LIFEBOAT_GITHUB_LISTEN")
	if address == "" {
		address = "127.0.0.1:8081"
	}
	mux := http.NewServeMux()
	mux.Handle("POST /webhook", adapter)
	server := &http.Server{Addr: address, Handler: mux, ReadHeaderTimeout: 10 * time.Second}
	log.Printf("Lifeboat GitHub adapter listening on %s", address)
	log.Fatal(server.ListenAndServe())
}
