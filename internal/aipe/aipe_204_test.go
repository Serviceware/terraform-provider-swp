package aipe

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// AIPE responds to successful PATCH updates with 204 No Content (same as it
// already does for DELETE). UpdateObject must accept that, not just 200.
func TestUpdateObject_Accepts204(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := &AIPEClient{HTTPClient: server.Client(), URL: server.URL}

	if err := client.UpdateObject(context.Background(), "B2XNZPP", map[string]string{"name": "User Management"}); err != nil {
		t.Fatalf("UpdateObject returned error, want nil: %v", err)
	}
}

// Same for the links PATCH endpoint (add/remove), used for both updating and
// "deleting" links.
func TestUpdateDataObjectLinks_Accepts204(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := &AIPEClient{HTTPClient: server.Client(), URL: server.URL}

	err := client.UpdateDataObjectLinks(context.Background(), "QDNNTS6", "link", "owns", nil, []string{"X"})
	if err != nil {
		t.Fatalf("UpdateDataObjectLinks returned error, want nil: %v", err)
	}
}
