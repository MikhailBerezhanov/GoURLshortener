package http_server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"url_shortener/url"
)

var server http.Server

func finishWithError(w http.ResponseWriter, msg string, status int) {
	log.Println(msg)
	http.Error(w, msg, status)
}

// Wrapper to add recovery for request handlers
func recoverHandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				finishWithError(w, fmt.Sprintf("Handler panic with error: %v", err), http.StatusInternalServerError)
			}
		}()

		handler.ServeHTTP(w, request)
	}
}

func getURLfromRequestBody(r *http.Request) (string, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	var rec url.Record
	err = json.Unmarshal(body, &rec)
	if err != nil {
		return "", err
	}

	if len(rec.URL) == 0 {
		return "", fmt.Errorf("URL field is missing or empty")
	}

	return rec.URL, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	reqURL, err := getURLfromRequestBody(r)
	if err != nil {
		finishWithError(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusBadRequest)
		return
	}

	rec := url.NewRecord(reqURL)
	data, err := json.Marshal(rec)
	if err != nil {
		log.Println(err)
		finishWithError(w, fmt.Sprintf("Failed to create response json: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, string(data))
}

func Start(port uint16) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", recoverHandler(handler))

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
