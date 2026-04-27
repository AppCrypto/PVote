package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	app, err := NewApp()
	if err != nil {
		log.Fatalf("failed to initialize demo app: %v", err)
	}

	addr := os.Getenv("PVOTE_WEB_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	log.Printf("PVote demo listening on http://127.0.0.1%s", addr)
	if err := http.ListenAndServe(addr, app.routes()); err != nil {
		log.Fatalf("demo server stopped: %v", err)
	}
}
