package server

import (
	"backend/internal/middleware"
	"backend/internal/test"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUserHandler(t *testing.T) {
	s := &Server{
		SupabaseFactory: test.NewMockSupabaseFactory(&test.MockAuth{}),
		DbConnector:     test.NewMockDbConnector(),
	}
	handler := middleware.CheckJwtToken(http.HandlerFunc(s.GetUser))
	server := httptest.NewServer(handler)
	defer server.Close()

	t.Setenv("JWT_SECRET", "testsecret")

	status_code, body := test.DoTestCall(t, server, "GET", test.TEST_TOKEN, nil, "")

	if status_code != http.StatusOK {
		t.Errorf("expected status OK; got %v", status_code)
	}

	if _, ok := body["profile"]; !ok {
		t.Error("expected response to contain profile key")
	}
}
