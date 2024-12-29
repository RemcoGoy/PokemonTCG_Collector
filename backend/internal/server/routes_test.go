package server

import (
	"backend/internal/test"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOkHandler(t *testing.T) {
	s := &Server{}
	server := httptest.NewServer(http.HandlerFunc(s.okHandler))
	defer server.Close()

	status_code, body := test.DoTestCall(t, server, "GET", "", nil)

	if status_code != http.StatusOK {
		t.Errorf("expected status OK; got %v", status_code)
	}

	if ok, exists := body["ok"]; !exists || ok != true {
		t.Error("expected response to contain ok key with value true")
	}
}
