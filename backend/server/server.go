package http_server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {

			// http.Err
			log.Printf("handler recovered error: %v\n", err)
			// w.WriteHeader(http.StatusInternalServerError)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}()

	fmt.Fprintf(w, "Hello, World!")
}

var server http.Server

func Start(port uint16) {
	// http.HandleFunc("/shorten", handler) // each request will call handler function

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
