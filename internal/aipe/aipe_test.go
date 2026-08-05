package aipe

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Reproduces the panic seen in production: AIPE returns a non-string value
// (a nested object) for a property, e.g. "external-system-reference": {}.
// GetObject must not crash and must still return the other, valid string
// properties.
func TestGetObject_NonStringProperty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"dataObject": {
				"system": {"id": "B2XNZPP"},
				"name": "User Management",
				"comment": null,
				"external-system-reference": {}
			}
		}`))
	}))
	defer server.Close()

	client := &AIPEClient{HTTPClient: server.Client(), URL: server.URL}

	object, err := client.GetObject(context.Background(), "B2XNZPP")
	if err != nil {
		t.Fatalf("GetObject returned error, want nil: %v", err)
	}

	if got, want := object["name"], "User Management"; got != want {
		t.Errorf("object[\"name\"] = %q, want %q", got, want)
	}

	if _, present := object["external-system-reference"]; present {
		t.Errorf("object[\"external-system-reference\"] should have been dropped, got present")
	}

	if _, present := object["system"]; present {
		t.Errorf("object[\"system\"] should have been removed, got present")
	}
}
