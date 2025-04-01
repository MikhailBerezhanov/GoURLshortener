package http_server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"url_shortener/url"
)

var server http.Server

var recordStore url.RecordStore // TODO: usa handlers context instead

func finishWithError(w http.ResponseWriter, msg string, status int) {
	log.Println(msg)
	http.Error(w, msg, status)
}

// Wrapper to add recovery for request handlers
func recoveryHandler(handler http.HandlerFunc) http.HandlerFunc {
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

func shortURLcreationHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		finishWithError(w, fmt.Sprintf("Unsupported Content-Type: %s", contentType), http.StatusUnsupportedMediaType)
		return
	}

	reqURL, err := getURLfromRequestBody(r)
	if err != nil {
		finishWithError(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusBadRequest)
		return
	}

	rec := url.NewRecord(reqURL)
	id, err := recordStore.InsertRecord(*rec)
	if err != nil {
		finishWithError(w, fmt.Sprintf("Failed to insert record to store: %v", err), http.StatusInternalServerError)
		return
	}
	rec.Id = id

	data, err := json.Marshal(rec)
	if err != nil {
		finishWithError(w, fmt.Sprintf("Failed to create response json: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, string(data))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	var shortCode string
	if len(parts) > 1 {
		shortCode = parts[len(parts)-1]
	}
	if len(shortCode) == 0 {
		finishWithError(w, fmt.Sprintf("Invalid URL path format: %s", path), http.StatusBadRequest)
		return
	}

	rec, err := recordStore.SelectRecord(shortCode)
	if err != nil {
		finishWithError(w, fmt.Sprintf("Failed to select record from store: %q", shortCode), http.StatusInternalServerError)
		return
	}

	// TODO: separate func
	data, err := json.Marshal(rec)
	if err != nil {
		finishWithError(w, fmt.Sprintf("Failed to create response json: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, string(data))
}

func Start(port uint16, dataStore url.RecordStore) {
	recordStore = dataStore

	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", recoveryHandler(shortURLcreationHandler))
	mux.HandleFunc("GET /shorten/", recoveryHandler(getHandler))

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
