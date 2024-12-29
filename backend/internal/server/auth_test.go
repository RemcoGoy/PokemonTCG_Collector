package server

import (
	"encoding/json"
	"fmt"
	"io"
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

	resp, err := http.Post(server.URL, "application/json", strings.NewReader(fmt.Sprintf("{\"email\":\"%s\",\"password\":\"%s\"}", USER_EMAIL, USER_PWD)))
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatalf("error unmarshaling response body: %v", err)
	}

	if _, ok := response["token"]; !ok {
		t.Error("expected response to contain token key")
	}
}

func TestLoginFailed(t *testing.T) {
	s := &Server{
		SupabaseFactory: test.NewMockSupabaseFactory(&test.MockFailedAuth{}),
	}
	server := httptest.NewServer(http.HandlerFunc(s.Login))
	defer server.Close()

	resp, err := http.Post(server.URL, "application/json", strings.NewReader(fmt.Sprintf("{\"email\":\"%s\",\"password\":\"%s\"}", USER_EMAIL, USER_PWD)))
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status BadRequest; got %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatalf("error unmarshaling response body: %v", err)
	}

	if _, ok := response["error"]; !ok {
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

	resp, err := http.Post(server.URL, "application/json", strings.NewReader(fmt.Sprintf("{\"email\":\"%s\",\"password\":\"%s\",\"username\":\"%s\"}", USER_EMAIL, USER_PWD, USERNAME)))
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatalf("error unmarshaling response body: %v", err)
	}

	if _, ok := response["email"]; !ok {
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

	resp, err := http.Post(server.URL, "application/json", strings.NewReader(fmt.Sprintf("{\"email\":\"%s\",\"password\":\"%s\",\"username\":\"%s\"}", USER_EMAIL, USER_PWD, USERNAME)))
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status BadRequest; got %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatalf("error unmarshaling response body: %v", err)
	}

	if _, ok := response["error"]; !ok {
		t.Error("expected response to contain error key")
	}
}

func TestLogoutHandler(t *testing.T) {
	s := &Server{
		SupabaseFactory: test.NewMockSupabaseFactory(&test.MockAuth{}),
	}
	server := httptest.NewServer(http.HandlerFunc(s.Logout))
	defer server.Close()

	req, err := http.NewRequest("POST", server.URL, nil)
	if err != nil {
		t.Fatalf("error creating request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+TEST_TOKEN)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
}

func TestLogoutWithoutToken(t *testing.T) {
	s := &Server{
		SupabaseFactory: test.NewMockSupabaseFactory(&test.MockAuth{}),
	}
	server := httptest.NewServer(http.HandlerFunc(s.Logout))
	defer server.Close()

	resp, err := http.Post(server.URL, "application/json", nil)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status Unauthorized; got %v", resp.Status)
	}
}
