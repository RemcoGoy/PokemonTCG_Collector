package server

import (
	"backend/internal/middleware"
	"backend/internal/test"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListCollectionsHandler(t *testing.T) {
	s := &Server{
		SupabaseFactory: test.NewMockSupabaseFactory(&test.MockAuth{}),
		DbConnector:     test.NewMockDbConnector(),
	}
	handler := middleware.CheckJwtToken(http.HandlerFunc(s.ListCollections))
	server := httptest.NewServer(handler)
	defer server.Close()

	t.Setenv("JWT_SECRET", "testsecret")

	status_code, body := test.DoTestCall(t, server, "GET", test.TEST_TOKEN, nil, "")

	if status_code != http.StatusOK {
		t.Errorf("expected status OK; got %v", status_code)
	}

	if _, ok := body["data"]; !ok {
		t.Error("expected response to contain data key")
	}

	if _, ok := body["count"]; !ok {
		t.Error("expected response to contain count key")
	}
}

func TestGetCollectionHandler(t *testing.T) {
	s := &Server{
		SupabaseFactory: test.NewMockSupabaseFactory(&test.MockAuth{}),
		DbConnector:     test.NewMockDbConnector(),
	}
	handler := middleware.CheckJwtToken(s.CollectionCtx(http.HandlerFunc(s.GetCollection)))
	server := httptest.NewServer(handler)
	defer server.Close()

	t.Setenv("JWT_SECRET", "testsecret")
	status_code, body := test.DoTestCall(t, server, "GET", test.TEST_TOKEN, nil, "/"+test.COLLECTION_ID)

	if status_code != http.StatusOK {
		t.Errorf("expected status OK; got %v", status_code)
	}

	if name, ok := body["name"]; !ok || name != test.COLLECTION_NAME {
		t.Errorf("expected response to contain name key with value %s; got %v", test.COLLECTION_NAME, name)
	}
}

func TestCreateCollectionHandler(t *testing.T) {
	s := &Server{
		SupabaseFactory: test.NewMockSupabaseFactory(&test.MockAuth{}),
		DbConnector:     test.NewMockDbConnector(),
	}
	handler := middleware.CheckJwtToken(http.HandlerFunc(s.CreateCollection))
	server := httptest.NewServer(handler)
	defer server.Close()

	t.Setenv("JWT_SECRET", "testsecret")

	status_code, body := test.DoTestCall(t, server, "POST", test.TEST_TOKEN, strings.NewReader(fmt.Sprintf("{\"name\":\"%s\"}", test.COLLECTION_NAME)), "")

	if status_code != http.StatusOK {
		t.Errorf("expected status OK; got %v", status_code)
	}

	if name, ok := body["name"]; !ok || name != test.COLLECTION_NAME {
		t.Errorf("expected response to contain name key with value %s; got %v", test.COLLECTION_NAME, name)
	}

	if id, ok := body["id"]; !ok || id != test.COLLECTION_ID {
		t.Errorf("expected response to contain id key with value %s; got %v", test.COLLECTION_ID, id)
	}
}
