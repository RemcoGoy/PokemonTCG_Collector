package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"backend/internal/test"
)

func TestLoginHandler(t *testing.T) {
	s := &Server{
		SupabaseFactory: test.NewMockSupabaseFactory(&test.MockAuth{}),
	}
	server := httptest.NewServer(http.HandlerFunc(s.Login))
	defer server.Close()

	status_code, body := test.DoTestCall(t, server, "POST", "", strings.NewReader(fmt.Sprintf("{\"email\":\"%s\",\"password\":\"%s\"}", test.USER_EMAIL, test.USER_PWD)), "")

	if status_code != http.StatusOK {
		t.Errorf("expected status OK; got %v", status_code)
	}

	if _, ok := body["token"]; !ok {
		t.Error("expected response to contain token key")
	}
}

func TestLoginFailed(t *testing.T) {
	s := &Server{
		SupabaseFactory: test.NewMockSupabaseFactory(&test.MockFailedAuth{}),
	}
	server := httptest.NewServer(http.HandlerFunc(s.Login))
	defer server.Close()

	status_code, body := test.DoTestCall(t, server, "POST", "", strings.NewReader(fmt.Sprintf("{\"email\":\"%s\",\"password\":\"%s\"}", test.USER_EMAIL, test.USER_PWD)), "")

	if status_code != http.StatusBadRequest {
		t.Errorf("expected status BadRequest; got %v", status_code)
	}

	if _, ok := body["error"]; !ok {
		t.Error("expected response to contain error key")
	}
}

func TestSignupHandler(t *testing.T) {
	s := &Server{
		SupabaseFactory: test.NewMockSupabaseFactory(&test.MockAuth{}),
		DbConnector:     test.NewMockDbConnector(),
	}
	server := httptest.NewServer(http.HandlerFunc(s.Signup))
	defer server.Close()

	status_code, body := test.DoTestCall(t, server, "POST", "", strings.NewReader(fmt.Sprintf("{\"email\":\"%s\",\"password\":\"%s\",\"username\":\"%s\"}", test.USER_EMAIL, test.USER_PWD, test.USERNAME)), "")

	if status_code != http.StatusOK {
		t.Errorf("expected status OK; got %v", status_code)
	}

	if _, ok := body["email"]; !ok {
		t.Error("expected response to contain email key")
	}
}

func TestSignupFailed(t *testing.T) {
	s := &Server{
		SupabaseFactory: test.NewMockSupabaseFactory(&test.MockFailedAuth{}),
		DbConnector:     test.NewMockDbConnector(),
	}
	server := httptest.NewServer(http.HandlerFunc(s.Signup))
	defer server.Close()

	status_code, body := test.DoTestCall(t, server, "POST", "", strings.NewReader(fmt.Sprintf("{\"email\":\"%s\",\"password\":\"%s\",\"username\":\"%s\"}", test.USER_EMAIL, test.USER_PWD, test.USERNAME)), "")

	if status_code != http.StatusBadRequest {
		t.Errorf("expected status BadRequest; got %v", status_code)
	}

	if _, ok := body["error"]; !ok {
		t.Error("expected response to contain error key")
	}
}

func TestLogoutHandler(t *testing.T) {
	s := &Server{
		SupabaseFactory: test.NewMockSupabaseFactory(&test.MockAuth{}),
	}
	server := httptest.NewServer(http.HandlerFunc(s.Logout))
	defer server.Close()

	status_code, _ := test.DoTestCall(t, server, "POST", test.TEST_TOKEN, nil, "")

	if status_code != http.StatusOK {
		t.Errorf("expected status OK; got %v", status_code)
	}
}

func TestLogoutWithoutToken(t *testing.T) {
	s := &Server{
		SupabaseFactory: test.NewMockSupabaseFactory(&test.MockAuth{}),
	}
	server := httptest.NewServer(http.HandlerFunc(s.Logout))
	defer server.Close()

	status_code, _ := test.DoTestCall(t, server, "POST", "", nil, "")

	if status_code != http.StatusUnauthorized {
		t.Errorf("expected status Unauthorized; got %v", status_code)
	}
}
