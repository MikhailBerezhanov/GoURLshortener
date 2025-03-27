package http_server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
	"url_shortener/url"
)

var server http.Server

func handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("handler recovered error: %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}()

	rec := url.NewRecord(r.URL.String()) // TODO : take field from POST body
	data, err := json.Marshal(rec)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to create record json", http.StatusInternalServerError)
	}
	fmt.Fprintln(w, string(data))
}

func Start(port uint16) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /shorten", handler)

	addr := fmt.Sprintf("localhost:%d", port)

	server = http.Server{Addr: addr, Handler: mux}

	log.Printf("Server running on %q\n", addr)

	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("Server closed")
	} else if err != nil {
		log.Fatal(err)
	}
}

func Stop() {
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown error: %v\n", err)
	}

	log.Println("Server stopped")
}
