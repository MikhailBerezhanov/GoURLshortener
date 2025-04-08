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

type contextKey string

const recordStoreKey contextKey = "recordStore"

var server http.Server

func finishWithError(w http.ResponseWriter, msg string, status int) {
	log.Println(msg)
	http.Error(w, msg, status)
}

// Wrapper to add recovery and context for request handlers
func withRecoveryAndContext(handler http.HandlerFunc, recordStore url.RecordStore) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				finishWithError(w, fmt.Sprintf("Handler panic with error: %v", err), http.StatusInternalServerError)
			}
		}()

		ctx := context.WithValue(request.Context(), recordStoreKey, recordStore)

		handler.ServeHTTP(w, request.WithContext(ctx))
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

func getRecordStoreFromRequestContext(r *http.Request) url.RecordStore {
	return r.Context().Value(recordStoreKey).(url.RecordStore)
}

func sendJsonResponse(w http.ResponseWriter, rec *url.Record) {
	data, err := json.Marshal(rec)
	if err != nil {
		finishWithError(w, fmt.Sprintf("Failed to create response json: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, string(data))
}

func postHandler(w http.ResponseWriter, r *http.Request) {
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
	recordStore := getRecordStoreFromRequestContext(r)
	if err := recordStore.InsertRecord(rec); err != nil {
		finishWithError(w, fmt.Sprintf("Failed to insert record to store: %v", err), http.StatusInternalServerError)
		return
	}

	sendJsonResponse(w, rec)
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

	recordStore := getRecordStoreFromRequestContext(r)
	rec, err := recordStore.SelectRecord(shortCode)
	if err != nil {
		if errors.Is(err, url.ErrRecordNotExist) {
			finishWithError(w, fmt.Sprintf("No data for requested shortCode: %q", shortCode), http.StatusNotFound)
		} else {
			finishWithError(w, fmt.Sprintf("Failed to select record from store: %q", shortCode), http.StatusInternalServerError)
		}

		return
	}

	sendJsonResponse(w, &rec)
}

func Start(port uint16, recordStore url.RecordStore) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", withRecoveryAndContext(postHandler, recordStore))
	mux.HandleFunc("GET /shorten/", withRecoveryAndContext(getHandler, recordStore))

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
