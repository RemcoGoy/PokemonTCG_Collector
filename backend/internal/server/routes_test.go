package server

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/supabase-community/gotrue-go/types"
)

type MockSupabaseClient struct {
	Auth *MockAuth
}

type MockAuth struct {
	SignupFunc func(types.SignupRequest) (types.SignupResponse, error)
}

func (m *MockAuth) Signup(req types.SignupRequest) (types.SignupResponse, error) {
	return m.SignupFunc(req)
}

func TestOkHandler(t *testing.T) {
	s := &Server{}
	server := httptest.NewServer(http.HandlerFunc(s.okHandler))
	defer server.Close()
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()
	// Assertions
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	expected := "{\"ok\":true}"
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}

func TestSignupHandler(t *testing.T) {
	user_email := "remco.goy@hotmail.com"
	user_password := "testpwd"

	s := &Server{}
	server := httptest.NewServer(http.HandlerFunc(s.Signup))
	defer server.Close()
	resp, err := http.Post(server.URL, "application/json", strings.NewReader(fmt.Sprintf("{\"email\":\"%s\",\"password\":\"%s\"}", user_email, user_password)))
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	expected := fmt.Sprintf("{\"email\":\"%s\"}", user_email)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}
