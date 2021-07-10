package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_Index(t *testing.T) {
	s := New(":80", "../../templates")

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

func Test_Login(t *testing.T) {
	s := New(":80", "../../templates")

	request := httptest.NewRequest("GET", "/login/", nil)
	responseWriter := httptest.NewRecorder()
	s.Handle(responseWriter, request)

	response := responseWriter.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("GET /login status must be 200, got %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Error("error while parsing response body")
	}

	wantContent := "<title>Login</title>"
	if !strings.Contains(string(body), wantContent) {
		t.Errorf(`want body to contain %s, got = "%s"`, wantContent, body)
	}
}

func Test_Register(t *testing.T) {
	s := New(":80", "../../templates")

	request := httptest.NewRequest("GET", "/register/", nil)
	responseWriter := httptest.NewRecorder()
	s.Handle(responseWriter, request)

	response := responseWriter.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("GET /login status must be 200, got %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Error("error while parsing response body")
	}

	wantContent := "<title>Register</title>"
	if !strings.Contains(string(body), wantContent) {
		t.Errorf(`want body to contain %s, got = "%s"`, wantContent, body)
	}
}

func Test_Users(t *testing.T) {
	s := New(":80", "../../templates")

	request := httptest.NewRequest("GET", "/users/", nil)
	responseWriter := httptest.NewRecorder()
	s.Handle(responseWriter, request)

	response := responseWriter.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("GET /login status must be 200, got %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Error("error while parsing response body")
	}

	wantBody := "users"
	if wantBody != string(body) {
		t.Errorf(`want body = "%s", got = "%s"`, wantBody, body)
	}
}

func Test_Friends(t *testing.T) {
	s := New(":80", "../../templates")

	request := httptest.NewRequest("GET", "/friends/", nil)
	responseWriter := httptest.NewRecorder()
	s.Handle(responseWriter, request)

	response := responseWriter.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("GET /login status must be 200, got %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Error("error while parsing response body")
	}

	wantBody := "friends"
	if wantBody != string(body) {
		t.Errorf(`want body = "%s", got = "%s"`, wantBody, body)
	}
}

func Test_User(t *testing.T) {
	s := New(":80", "../../templates")

	request := httptest.NewRequest("GET", "/users/1", nil)
	responseWriter := httptest.NewRecorder()
	s.Handle(responseWriter, request)

	response := responseWriter.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("GET /login status must be 200, got %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Error("error while parsing response body")
	}

	wantBody := "user"
	if wantBody != string(body) {
		t.Errorf(`want body = "%s", got = "%s"`, wantBody, body)
	}
}
