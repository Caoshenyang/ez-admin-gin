package contract

import (
	"encoding/json"
	"os"
	"sort"
	"strings"
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

// TestResponseSchemaEnvelope verifies that every 200 response references the
// httpx.Body envelope, ensuring consistent API response structure.
func TestResponseSchemaEnvelope(t *testing.T) {
	spec := mustParseSwagger(t)

	for path, methods := range spec.Paths {
		for method, op := range methods {
			// setup/init uses a generic schema, not httpx.Body.
			if path == "/setup/init" {
				continue
			}
			opMap, ok := op.(map[string]interface{})
			if !ok {
				continue
			}
			responses, _ := opMap["responses"].(map[string]interface{})
			resp200, _ := responses["200"].(map[string]interface{})
			if resp200 == nil {
				t.Errorf("%s %s: missing 200 response", strings.ToUpper(method), path)
				continue
			}
			schema, _ := resp200["schema"].(map[string]interface{})
			if schema == nil {
				t.Errorf("%s %s: 200 response has no schema", strings.ToUpper(method), path)
				continue
			}
			if !schemaRefsBody(schema) {
				t.Errorf("%s %s: 200 response schema does not reference httpx.Body", strings.ToUpper(method), path)
			}
		}
	}
}

// schemaRefsBody checks if a schema directly or via allOf references httpx.Body.
func schemaRefsBody(schema map[string]interface{}) bool {
	if ref, ok := schema["$ref"].(string); ok && strings.HasSuffix(ref, "httpx.Body") {
		return true
	}
	allOf, _ := schema["allOf"].([]interface{})
	for _, item := range allOf {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if ref, ok := itemMap["$ref"].(string); ok && strings.HasSuffix(ref, "httpx.Body") {
			return true
		}
	}
	return false
}

// TestDefinitionsReachable verifies that every $ref in paths and definitions
// points to an existing definition in the spec.
func TestDefinitionsReachable(t *testing.T) {
	spec := mustParseSwagger(t)

	defs := make(map[string]bool)
	for name := range spec.Definitions {
		defs["#/definitions/"+name] = true
	}

	var missing []string
	visitRefs(spec.Paths, func(ref string) {
		if strings.HasPrefix(ref, "#/definitions/") && !defs[ref] {
			missing = append(missing, ref)
		}
	})
	visitRefs(spec.Definitions, func(ref string) {
		if strings.HasPrefix(ref, "#/definitions/") && !defs[ref] {
			missing = append(missing, ref)
		}
	})

	if len(missing) > 0 {
		unique := make(map[string]bool)
		for _, m := range missing {
			unique[m] = true
		}
		sorted := make([]string, 0, len(unique))
		for m := range unique {
			sorted = append(sorted, m)
		}
		sort.Strings(sorted)
		t.Errorf("unresolved $ref references: %v", sorted)
	}
}

// TestKeyEndpointDataSchemas verifies that important read endpoints have typed
// data schemas in their 200 response (not just a bare httpx.Body reference).
func TestKeyEndpointDataSchemas(t *testing.T) {
	spec := mustParseSwagger(t)

	endpoints := []struct {
		path       string
		method     string
		wantDataRef string
	}{
		{"/auth/login", "post", "domain.LoginResponse"},
		{"/auth/me", "get", "domain.MeResponse"},
		{"/auth/account", "get", "domain.AccountProfileResponse"},
	}

	for _, ep := range endpoints {
		methods, ok := spec.Paths[ep.path]
		if !ok {
			t.Errorf("path %s not found", ep.path)
			continue
		}
		op, ok := methods[ep.method]
		if !ok {
			t.Errorf("method %s not found for path %s", ep.method, ep.path)
			continue
		}
		opMap, _ := op.(map[string]interface{})
		responses, _ := opMap["responses"].(map[string]interface{})
		resp200, _ := responses["200"].(map[string]interface{})
		schema, _ := resp200["schema"].(map[string]interface{})
		if schema == nil {
			t.Errorf("%s %s: 200 response has no schema", ep.method, ep.path)
			continue
		}

		dataRef := extractDataRef(schema)
		if dataRef == "" {
			t.Errorf("%s %s: 200 response has no typed data schema (expected a data.$ref)", ep.method, ep.path)
		}
	}
}

// --- helpers ---

type swaggerSpec struct {
	Paths       map[string]map[string]interface{} `json:"paths"`
	Definitions map[string]interface{}             `json:"definitions"`
}

func mustParseSwagger(t *testing.T) swaggerSpec {
	t.Helper()
	data, err := os.ReadFile(swaggerPath)
	if err != nil {
		t.Fatalf("read swagger.json: %v", err)
	}
	var spec swaggerSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		t.Fatalf("parse swagger.json: %v", err)
	}
	return spec
}

// extractDataRef extracts the $ref from a "data" property inside an allOf schema.
func extractDataRef(schema map[string]interface{}) string {
	allOf, _ := schema["allOf"].([]interface{})
	for _, item := range allOf {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		props, _ := itemMap["properties"].(map[string]interface{})
		data, _ := props["data"].(map[string]interface{})
		if ref, ok := data["$ref"].(string); ok {
			return ref
		}
	}
	return ""
}

// visitRefs walks a JSON structure and calls fn for every $ref string found.
func visitRefs(v interface{}, fn func(string)) {
	switch val := v.(type) {
	case map[string]interface{}:
		if ref, ok := val["$ref"].(string); ok {
			fn(ref)
		}
		for _, child := range val {
			visitRefs(child, fn)
		}
	case []interface{}:
		for _, child := range val {
			visitRefs(child, fn)
		}
	}
}
