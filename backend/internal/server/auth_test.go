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

func TestLoginHandler(t *testing.T) {
	user_email := "remco.goy@hotmail.com"
	user_password := "testpwd"

	s := &Server{
		SupabaseFactory: test.NewMockSupabaseFactory(&test.MockAuth{}),
	}
	server := httptest.NewServer(http.HandlerFunc(s.Login))
	defer server.Close()

	resp, err := http.Post(server.URL, "application/json", strings.NewReader(fmt.Sprintf("{\"email\":\"%s\",\"password\":\"%s\"}", user_email, user_password)))
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
	user_email := "remco.goy@hotmail.com"
	user_password := "testpwd"

	s := &Server{
		SupabaseFactory: test.NewMockSupabaseFactory(&test.MockFailedAuth{}),
	}
	server := httptest.NewServer(http.HandlerFunc(s.Login))
	defer server.Close()

	resp, err := http.Post(server.URL, "application/json", strings.NewReader(fmt.Sprintf("{\"email\":\"%s\",\"password\":\"%s\"}", user_email, user_password)))
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
