package aipe

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// AIPE validates certain properties (e.g. "uses-postgres") as JSON booleans.
// The Terraform-facing schema only supports string properties, so HCL bools
// arrive here as the strings "true"/"false". The wire payload must send
// those as real JSON booleans, not quoted strings, while leaving other
// string values untouched.
func TestUpdateObject_SendsBooleanLikePropertiesAsJSONBool(t *testing.T) {
	var capturedBody []byte

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := &AIPEClient{HTTPClient: server.Client(), URL: server.URL}

	err := client.UpdateObject(context.Background(), "X", map[string]string{
		"name":           "production-acme-certs",
		"uses-postgres":  "true",
		"uses-something": "false",
	})
	if err != nil {
		t.Fatalf("UpdateObject returned error, want nil: %v", err)
	}

	var payload struct {
		DataObject map[string]interface{} `json:"dataObject"`
	}
	if err := json.Unmarshal(capturedBody, &payload); err != nil {
		t.Fatalf("could not parse request body %s: %v", capturedBody, err)
	}

	if got, want := payload.DataObject["name"], "production-acme-certs"; got != want {
		t.Errorf("name = %v (%T), want %v (string)", got, got, want)
	}

	if got, ok := payload.DataObject["uses-postgres"].(bool); !ok || got != true {
		t.Errorf("uses-postgres = %v (%T), want true (bool)", payload.DataObject["uses-postgres"], payload.DataObject["uses-postgres"])
	}

	if got, ok := payload.DataObject["uses-something"].(bool); !ok || got != false {
		t.Errorf("uses-something = %v (%T), want false (bool)", payload.DataObject["uses-something"], payload.DataObject["uses-something"])
	}
}

func TestCreateObject_SendsBooleanLikePropertiesAsJSONBool(t *testing.T) {
	var capturedBody []byte

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"dataObjectId": "NEWID"}`))
	}))
	defer server.Close()

	client := &AIPEClient{HTTPClient: server.Client(), URL: server.URL}

	_, err := client.CreateObject(context.Background(), "sce-service-instance", map[string]string{
		"name":          "production-acme-certs",
		"uses-postgres": "true",
	})
	if err != nil {
		t.Fatalf("CreateObject returned error, want nil: %v", err)
	}

	var payload struct {
		DataObject map[string]interface{} `json:"dataObject"`
	}
	if err := json.Unmarshal(capturedBody, &payload); err != nil {
		t.Fatalf("could not parse request body %s: %v", capturedBody, err)
	}

	if got, ok := payload.DataObject["uses-postgres"].(bool); !ok || got != true {
		t.Errorf("uses-postgres = %v (%T), want true (bool)", payload.DataObject["uses-postgres"], payload.DataObject["uses-postgres"])
	}
}
