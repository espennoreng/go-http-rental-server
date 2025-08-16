package api

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

type errorResponse struct {
	Message string `json:"message"`
}

// AssertStatus checks if the HTTP response recorder has the expected status code.
func AssertStatus(t *testing.T, res *httptest.ResponseRecorder, expected int) {
	t.Helper() // Mark this function as a test helper.

	if res.Code != expected {
		t.Errorf("expected status %d, got %d", expected, res.Code)
	}
}

// AssertErrorBody checks if the response body contains the expected error message.
// It also handles the trailing newline that http.Error adds.
func AssertJSONErrorBody(t *testing.T, res *httptest.ResponseRecorder, expected string) {
	t.Helper() // Mark this function as a test helper.

	var response errorResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if !strings.Contains(response.Message, expected) {
		t.Errorf("expected error message to contain '%s', got '%s'", expected, response.Message)
	}
}

// AssertJSONContentType checks if the response has the correct Content-Type header for JSON.
func AssertJSONContentType(t *testing.T, res *httptest.ResponseRecorder) {
	t.Helper() // Mark this function as a test helper.

	if res.Header().Get(ContentType) != ContentTypeJSON {
		t.Errorf("expected Content-Type application/json, got %s", res.Header().Get(ContentType))
	}
}

