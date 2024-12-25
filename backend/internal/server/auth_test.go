package server

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"backend/internal/test"

	"github.com/supabase-community/gotrue-go/types"
)

func TestSignupHandler(t *testing.T) {
	user_email := "remco.goy@hotmail.com"
	user_password := "testpwd"

	s := &Server{
		SupabaseFactory: test.NewMockSupabaseFactory(),
	}
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

type FailedSignupAuth struct {
	test.MockAuth
}

func (m *FailedSignupAuth) Signup(req types.SignupRequest) (*types.SignupResponse, error) {
	return nil, errors.New("error")
}

func TestSignupFailedHandler(t *testing.T) {
	user_email := "remco.goy@hotmail.com"
	user_password := "testpwd"

	auth := &FailedSignupAuth{}
	factory := test.NewMockSupabaseFactory()
	factory.Auth = auth

	s := &Server{
		SupabaseFactory: factory,
	}
	server := httptest.NewServer(http.HandlerFunc(s.Signup))
	defer server.Close()

	resp, err := http.Post(server.URL, "application/json", strings.NewReader(fmt.Sprintf("{\"email\":\"%s\",\"password\":\"%s\"}", user_email, user_password)))
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status BadRequest; got %v", resp.Status)
	}

	expected := fmt.Sprintf("{\"error\":\"%s\"}", "error")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}
