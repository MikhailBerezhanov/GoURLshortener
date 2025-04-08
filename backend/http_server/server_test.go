package http_server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"url_shortener/datastore"
)

func TestPostHandlerShouldCreateNewRecord(t *testing.T) {
	body := `{"url": "http://test/url.org"}`
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockDB := datastore.NewMemDB()
	handler := withRecoveryAndContext(postHandler, mockDB)
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201 Created, got %v", resp.Status)
	}

	// buf := new(strings.Builder)
	// _, _ = buf.ReadFrom(resp.Body)
	// expected := "Received: Hello from test"
	// if buf.String() != expected {
	// 	t.Errorf("expected response body %q, got %q", expected, buf.String())
	// }
}
