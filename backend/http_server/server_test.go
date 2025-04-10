package http_server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"url_shortener/datastore"
	"url_shortener/url"
)

func createMockRequest(method, url, body string) (req *http.Request, w *httptest.ResponseRecorder) {
	req = httptest.NewRequest(method, url, bytes.NewBufferString(body))
	w = httptest.NewRecorder()
	return
}

func TestPostHandlerShouldReturn400OnBadRequest(t *testing.T) {
	body := `{"invalidFieldName": "http://test/url.org"}`
	req, w := createMockRequest(http.MethodPost, "/shorten", body)
	req.Header.Set("Content-Type", "application/json")

	mockDB := datastore.NewMemDB()
	handler := withRecoveryAndContext(postHandler, mockDB)
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 BadRequest, got %v", resp.Status)
	}
}

func TestPostHandlerShouldReturn415OnUnsupportedContentType(t *testing.T) {
	body := `{"url": "http://test/url.org"}`
	req, w := createMockRequest(http.MethodPost, "/shorten", body)
	req.Header.Set("Content-Type", "text/plain")

	mockDB := datastore.NewMemDB()
	handler := withRecoveryAndContext(postHandler, mockDB)
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnsupportedMediaType {
		t.Errorf("expected status 415 UnsupportedMediaType, got %v", resp.Status)
	}
}

func TestPostHandlerShouldCreateNewRecord(t *testing.T) {
	body := `{"url": "http://test/url.org"}`
	req, w := createMockRequest(http.MethodPost, "/shorten", body)
	req.Header.Set("Content-Type", "application/json")

	mockDB := datastore.NewMemDB()
	handler := withRecoveryAndContext(postHandler, mockDB)
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201 Created, got %v", resp.Status)
	}

	var rec url.Record
	if err := json.NewDecoder(resp.Body).Decode(&rec); err != nil {
		t.Errorf("expected record decoded from json response body, actual failed with err: %v", err)
	}

	recFromDB, err := mockDB.SelectRecord(rec.ShortCode)
	if err != nil {
		t.Errorf("expected record %s was inserted to db, actual select failed with err: %v", &rec, err)
	}

	// if rec != recFromDB {
	// 	t.Errorf("expected record selected from db and received from response will be eqaul, actual rec: %s, recFromDB: %s", &rec, &recFromDB)
	// }
}
