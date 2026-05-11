package contract

import (
	"encoding/json"
	"os"
	"testing"
)

const swaggerPath = "../../docs/swagger.json"

// TestSwaggerFileExists verifies the OpenAPI spec file is present.
func TestSwaggerFileExists(t *testing.T) {
	if _, err := os.Stat(swaggerPath); os.IsNotExist(err) {
		t.Fatalf("swagger.json not found at %s — run `swag init` to generate it", swaggerPath)
	}
}

// TestSwaggerParsable verifies the OpenAPI spec is valid JSON with expected top-level keys.
func TestSwaggerParsable(t *testing.T) {
	data, err := os.ReadFile(swaggerPath)
	if err != nil {
		t.Fatalf("read swagger.json: %v", err)
	}

	var spec map[string]json.RawMessage
	if err := json.Unmarshal(data, &spec); err != nil {
		t.Fatalf("swagger.json is not valid JSON: %v", err)
	}

	for _, key := range []string{"paths", "info", "swagger"} {
		if _, ok := spec[key]; !ok {
			t.Errorf("swagger.json missing top-level key: %s", key)
		}
	}
}

// TestKeyEndpointsDeclared verifies critical API endpoints exist in the OpenAPI spec.
func TestKeyEndpointsDeclared(t *testing.T) {
	data, err := os.ReadFile(swaggerPath)
	if err != nil {
		t.Fatalf("read swagger.json: %v", err)
	}

	var spec struct {
		Paths map[string]map[string]json.RawMessage `json:"paths"`
	}
	if err := json.Unmarshal(data, &spec); err != nil {
		t.Fatalf("parse swagger.json: %v", err)
	}

	requiredEndpoints := []struct {
		path   string
		method string
	}{
		{"/auth/login", "post"},
		{"/auth/me", "get"},
		{"/system/users", "get"},
		{"/system/roles", "get"},
		{"/system/menus", "get"},
	}

	for _, ep := range requiredEndpoints {
		methods, ok := spec.Paths[ep.path]
		if !ok {
			t.Errorf("path %s not found in OpenAPI spec", ep.path)
			continue
		}
		if _, ok := methods[ep.method]; !ok {
			t.Errorf("method %s not declared for path %s in OpenAPI spec", ep.method, ep.path)
		}
	}
}

// TestSwaggerInfo verifies the API metadata is correctly set.
func TestSwaggerInfo(t *testing.T) {
	data, err := os.ReadFile(swaggerPath)
	if err != nil {
		t.Fatalf("read swagger.json: %v", err)
	}

	var spec struct {
		Info struct {
			Title       string `json:"title"`
			Version     string `json:"version"`
			Description string `json:"description"`
		} `json:"info"`
		BasePath string `json:"basePath"`
	}
	if err := json.Unmarshal(data, &spec); err != nil {
		t.Fatalf("parse swagger.json: %v", err)
	}

	if spec.Info.Title == "" {
		t.Error("swagger info.title is empty")
	}
	if spec.Info.Version == "" {
		t.Error("swagger info.version is empty")
	}
	if spec.BasePath != "/api/v1" {
		t.Errorf("swagger basePath = %q, want /api/v1", spec.BasePath)
	}
}

// TODO: TestResponseSchemaConsistency — verify real HTTP responses match OpenAPI schema.
// This requires running the server against a test database and comparing the response
// JSON structure against the schema definitions in swagger.json. Deferred to next phase.
