//go:build integration

package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"ez-admin-gin/server/tests/testutil"
)

func TestHealthzReturnsOK(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	req, _ := http.NewRequest(http.MethodGet, app.URL("/healthz"), nil)
	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	var result responseBody
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if result.Code != 0 {
		t.Errorf("code = %d, want 0", result.Code)
	}

	var data struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(result.Data, &data); err != nil {
		t.Fatalf("unmarshal data: %v", err)
	}
	if data.Status != "ok" {
		t.Errorf("status = %q, want %q", data.Status, "ok")
	}
}

func TestReadyzReturnsOK(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	req, _ := http.NewRequest(http.MethodGet, app.URL("/readyz"), nil)
	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	var result responseBody
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if result.Code != 0 {
		t.Errorf("code = %d, want 0", result.Code)
	}

	var data struct {
		Database string `json:"database"`
		Redis    string `json:"redis"`
	}
	if err := json.Unmarshal(result.Data, &data); err != nil {
		t.Fatalf("unmarshal data: %v", err)
	}
	if data.Database != "ok" {
		t.Errorf("database = %q, want %q", data.Database, "ok")
	}
	if data.Redis != "ok" {
		t.Errorf("redis = %q, want %q", data.Redis, "ok")
	}
}

func TestMetricsEndpoint(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	// Generate some traffic so metrics are non-empty.
	app.SeedAdmin(t, "metricsadmin", "admin123", "Metrics Admin")

	req, _ := http.NewRequest(http.MethodGet, app.URL("/metrics"), nil)
	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}

	text := string(body)
	if !strings.Contains(text, "http_requests_total") {
		t.Error("metrics response should contain http_requests_total")
	}
	if !strings.Contains(text, "http_request_duration_seconds") {
		t.Error("metrics response should contain http_request_duration_seconds")
	}
}
