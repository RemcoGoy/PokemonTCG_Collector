package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"backend/internal/test"
)

const (
	USER_EMAIL = "test@test.com"
	USER_PWD   = "testpwd"
	USERNAME   = "test"
	TEST_TOKEN = "eyJhbGciOiJIUzI1NiIsImtpZCI6Ii9uKzhHVVJjVDRxcGZCUmMiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL2N5ZWhkbnV6b2x2eXByYWJ1eG11LnN1cGFiYXNlLmNvL2F1dGgvdjEiLCJzdWIiOiIwMDAwMDAwMC0wMDAwLTAwMDAtMDAwMC0wMDAwMDAwMDAwMDAiLCJhdWQiOiJhdXRoZW50aWNhdGVkIiwiZXhwIjoxLCJpYXQiOjEsImVtYWlsIjoidGVzdEB0ZXN0LmNvbSIsInBob25lIjoiIiwiYXBwX21ldGFkYXRhIjp7InByb3ZpZGVyIjoiZW1haWwiLCJwcm92aWRlcnMiOlsiZW1haWwiXX0sInVzZXJfbWV0YWRhdGEiOnsiZW1haWwiOiJ0ZXN0QHRlc3QuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsInBob25lX3ZlcmlmaWVkIjpmYWxzZSwic3ViIjoiMDAwMDAwMDAtMDAwMC0wMDAwLTAwMDAtMDAwMDAwMDAwMDAwIn0sInJvbGUiOiJhdXRoZW50aWNhdGVkIiwiYWFsIjoieCIsImFtciI6W3sibWV0aG9kIjoicGFzc3dvcmQiLCJ0aW1lc3RhbXAiOjF9XSwic2Vzc2lvbl9pZCI6IjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMSIsImlzX2Fub255bW91cyI6ZmFsc2V9.OJk7JzP6XBJh-C9AcNqo-ErqlUYhPd95RUdvb194xW8"
)

func TestLoginHandler(t *testing.T) {
	s := &Server{
		SupabaseFactory: test.NewMockSupabaseFactory(&test.MockAuth{}),
	}
	server := httptest.NewServer(http.HandlerFunc(s.Login))
	defer server.Close()

	status_code, body := test.DoTestCall(t, server, "POST", "", strings.NewReader(fmt.Sprintf("{\"email\":\"%s\",\"password\":\"%s\"}", USER_EMAIL, USER_PWD)))

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

	status_code, body := test.DoTestCall(t, server, "POST", "", strings.NewReader(fmt.Sprintf("{\"email\":\"%s\",\"password\":\"%s\"}", USER_EMAIL, USER_PWD)))

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

	status_code, body := test.DoTestCall(t, server, "POST", "", strings.NewReader(fmt.Sprintf("{\"email\":\"%s\",\"password\":\"%s\",\"username\":\"%s\"}", USER_EMAIL, USER_PWD, USERNAME)))

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

	status_code, body := test.DoTestCall(t, server, "POST", "", strings.NewReader(fmt.Sprintf("{\"email\":\"%s\",\"password\":\"%s\",\"username\":\"%s\"}", USER_EMAIL, USER_PWD, USERNAME)))

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

	status_code, _ := test.DoTestCall(t, server, "POST", TEST_TOKEN, nil)

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

	status_code, _ := test.DoTestCall(t, server, "POST", "", nil)

	if status_code != http.StatusUnauthorized {
		t.Errorf("expected status Unauthorized; got %v", status_code)
	}
}
