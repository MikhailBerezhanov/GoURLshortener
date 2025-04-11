package http_server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"url_shortener/datastore"
	"url_shortener/url"
)

type FailDBMock struct{}

func (f FailDBMock) InsertRecord(*url.Record) error { return fmt.Errorf("mock insert error") }
func (f FailDBMock) SelectRecord(string) (url.Record, error) {
	return url.Record{}, fmt.Errorf("mock insert error")
}

func createMockRequest(method, url, body string) (req *http.Request, w *httptest.ResponseRecorder) {
	var bodyReader io.Reader
	if len(body) > 0 {
		bodyReader = bytes.NewBufferString(body)
	}

	req = httptest.NewRequest(method, url, bodyReader)
	w = httptest.NewRecorder()
	return
}

func equal(r1, r2 *url.Record) bool {
	return r1.Id == r2.Id && r1.ShortCode == r2.ShortCode && r1.URL == r2.URL
}

func createHandlerUnderTest(h func(w http.ResponseWriter, request *http.Request)) (handler http.HandlerFunc, mockDB *datastore.MemDB) {
	mockDB = datastore.NewMemDB()
	handler = withRecoveryAndContext(h, mockDB)
	return
}

// POST /shorten tests

func TestPostHandlerShouldReturn400OnBadRequest(t *testing.T) {
	body := `{"invalidFieldName": "http://test/url.org"}`
	req, w := createMockRequest(http.MethodPost, "/shorten", body)
	req.Header.Set("Content-Type", "application/json")

	handler, _ := createHandlerUnderTest(postHandler)
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

	handler, _ := createHandlerUnderTest(postHandler)
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnsupportedMediaType {
		t.Errorf("expected status 415 UnsupportedMediaType, got %v", resp.Status)
	}
}

func TestPostHandlerShouldReturn500OnFailedInsertToDB(t *testing.T) {
	body := `{"url": "http://test/url.org"}`
	req, w := createMockRequest(http.MethodPost, "/shorten", body)
	req.Header.Set("Content-Type", "application/json")

	var mockDB FailDBMock
	handler := withRecoveryAndContext(postHandler, &mockDB)
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500 InternalServerError, got %v", resp.Status)
	}
}

func TestPostHandlerShouldCreateNewRecord(t *testing.T) {
	body := `{"url": "http://test/url.org"}`
	req, w := createMockRequest(http.MethodPost, "/shorten", body)
	req.Header.Set("Content-Type", "application/json")

	handler, mockDB := createHandlerUnderTest(postHandler)
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

	if !equal(&rec, &recFromDB) {
		t.Errorf("expected record selected from db and received from response will be eqaul, actual rec: %s, recFromDB: %s", &rec, &recFromDB)
	}
}

// `GET` /shorten/<shortCode> tests

func TestGetHandlerShouldReturn400OnBadRequest(t *testing.T) {
	req, w := createMockRequest(http.MethodGet, "/shorten/", "")

	handler, _ := createHandlerUnderTest(getHandler)
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 BadRequest, got %v", resp.Status)
	}
}

func TestGetHandlerShouldReturn404OnNotFoundShortCode(t *testing.T) {
	req, w := createMockRequest(http.MethodGet, "/shorten/abc123", "")

	handler, _ := createHandlerUnderTest(getHandler)
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404 NotFound, got %v", resp.Status)
	}
}

func TestGetHandlerShouldReturn500OnFailedSelectFromDB(t *testing.T) {
	req, w := createMockRequest(http.MethodGet, "/shorten/abc123", "")

	var mockDB FailDBMock
	handler := withRecoveryAndContext(getHandler, &mockDB)
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500 InternalServerError, got %v", resp.Status)
	}
}

func TestGetHandlerShouldReturnRecordByShortCode(t *testing.T) {
	mockDB := datastore.NewMemDB()
	testRecord := url.NewRecord("test/url")
	if err := mockDB.InsertRecord(testRecord); err != nil {
		t.Errorf("expected mockDB will InsertRecord, actual failed with error %v", err)
	}

	req, w := createMockRequest(http.MethodGet, "/shorten/"+testRecord.ShortCode, "")

	handler := withRecoveryAndContext(getHandler, mockDB)
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %v", resp.Status)
	}

	var receivedRecord url.Record
	if err := json.NewDecoder(resp.Body).Decode(&receivedRecord); err != nil {
		t.Errorf("expected record decoded from json response body, actual failed with err: %v", err)
	}

	if !equal(testRecord, &receivedRecord) {
		t.Errorf("expected testRecord and received from response will be eqaul, actual testRecord: %s, receivedRecord: %s", testRecord, &receivedRecord)
	}
}
