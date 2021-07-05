package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Index(t *testing.T) {
	s := New(":80")

	request := httptest.NewRequest("GET", "/", nil)
	responseWriter := httptest.NewRecorder()
	s.Handle(responseWriter, request)

	response := responseWriter.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("GET / status must be 200, got %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Error("error while parsing response body")
	}

	wantBody := "hello world"
	if wantBody != string(body) {
		t.Errorf(`want body = "%s", got = "%s"`, wantBody, body)
	}
}
