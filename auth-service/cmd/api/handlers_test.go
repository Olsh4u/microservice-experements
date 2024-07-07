package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

func Test_Authenticate(t *testing.T) {
	jsonToReturn := `
		{
			"error": false,
			"message": "You are authenticated"
		}
		`

	client := NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusAccepted,
			Body:       io.NopCloser(bytes.NewBufferString(jsonToReturn)),
			Header:     make(http.Header),
		}
	})

	testApp.Client = client

	postBody := map[string]interface{}{
		"email":    "test@example.com",
		"password": "password",
	}

	body, _ := json.Marshal(postBody)

	req, _ := http.NewRequest("POST", "/authenticate", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(testApp.Authenticate)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusAccepted)
	}
}
