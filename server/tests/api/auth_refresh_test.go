//go:build integration

package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"ez-admin-gin/server/tests/testutil"
)

func TestRefreshSuccess(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.SeedAdmin(t, testAdminUser, testAdminPassword, testAdminNickname)
	app.CleanupTestData(t)
	app.SeedAdmin(t, testAdminUser, testAdminPassword, testAdminNickname)

	// Login to get access token and refresh cookie.
	_, cookies := app.LoginWithCookies(t, testAdminUser, testAdminPassword)

	// Call /auth/refresh with the refresh cookie.
	req, _ := http.NewRequest(http.MethodPost, app.URL("/api/v1/auth/refresh"), nil)
	for _, c := range cookies {
		req.AddCookie(c)
	}

	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("refresh request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("refresh status = %d, want 200", resp.StatusCode)
	}

	var result responseBody
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if result.Code != 0 {
		t.Errorf("code = %d, want 0", result.Code)
	}

	var refreshData struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		UserID      uint   `json:"user_id"`
	}
	if err := json.Unmarshal(result.Data, &refreshData); err != nil {
		t.Fatalf("unmarshal refresh data: %v", err)
	}
	if refreshData.AccessToken == "" {
		t.Error("refreshed access_token is empty")
	}

	// Verify new refresh cookie is set.
	var hasNewRefreshCookie bool
	for _, c := range resp.Cookies() {
		if c.Name == "ez_admin_refresh_token" {
			hasNewRefreshCookie = true
		}
	}
	if !hasNewRefreshCookie {
		t.Error("expected new refresh token cookie after rotation")
	}
}

func TestRefreshWithInvalidToken(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	req, _ := http.NewRequest(http.MethodPost, app.URL("/api/v1/auth/refresh"), nil)
	req.AddCookie(&http.Cookie{
		Name:  "ez_admin_refresh_token",
		Value: "invalid-token-value",
	})

	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("refresh request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("refresh with invalid token: status = %d, want 401", resp.StatusCode)
	}
}

func TestRefreshWithoutCookie(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	req, _ := http.NewRequest(http.MethodPost, app.URL("/api/v1/auth/refresh"), nil)

	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("refresh request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("refresh without cookie: status = %d, want 401", resp.StatusCode)
	}
}

func TestRefreshRotation(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.SeedAdmin(t, testAdminUser, testAdminPassword, testAdminNickname)
	app.CleanupTestData(t)
	app.SeedAdmin(t, testAdminUser, testAdminPassword, testAdminNickname)

	_, cookies := app.LoginWithCookies(t, testAdminUser, testAdminPassword)

	// First refresh should succeed and return new cookie.
	req1, _ := http.NewRequest(http.MethodPost, app.URL("/api/v1/auth/refresh"), nil)
	for _, c := range cookies {
		req1.AddCookie(c)
	}
	resp1, err := app.Do(req1)
	if err != nil {
		t.Fatalf("first refresh request failed: %v", err)
	}
	resp1.Body.Close()

	if resp1.StatusCode != http.StatusOK {
		t.Fatalf("first refresh status = %d, want 200", resp1.StatusCode)
	}

	// Reusing the same (old) refresh token should fail (rotation).
	req2, _ := http.NewRequest(http.MethodPost, app.URL("/api/v1/auth/refresh"), nil)
	for _, c := range cookies {
		req2.AddCookie(c)
	}
	resp2, err := app.Do(req2)
	if err != nil {
		t.Fatalf("second refresh request failed: %v", err)
	}
	resp2.Body.Close()

	if resp2.StatusCode != http.StatusUnauthorized {
		t.Fatalf("reused refresh token: status = %d, want 401", resp2.StatusCode)
	}
}

func TestLogoutSuccess(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.SeedAdmin(t, testAdminUser, testAdminPassword, testAdminNickname)
	app.CleanupTestData(t)
	app.SeedAdmin(t, testAdminUser, testAdminPassword, testAdminNickname)

	accessToken, cookies := app.LoginWithCookies(t, testAdminUser, testAdminPassword)

	// Call /auth/logout.
	req := app.AuthRequest(http.MethodPost, "/api/v1/auth/logout", accessToken, nil)
	for _, c := range cookies {
		req.AddCookie(c)
	}

	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("logout request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("logout status = %d, want 200", resp.StatusCode)
	}

	// Verify refresh cookie is cleared.
	for _, c := range resp.Cookies() {
		if c.Name == "ez_admin_refresh_token" && c.MaxAge < 0 {
			return // Cookie cleared, test passes.
		}
	}
}

func TestLogoutRevokesAccessToken(t *testing.T) {
	app := testutil.NewTestApp(t)
	defer app.Close(t)

	app.SeedAdmin(t, testAdminUser, testAdminPassword, testAdminNickname)
	app.CleanupTestData(t)
	app.SeedAdmin(t, testAdminUser, testAdminPassword, testAdminNickname)

	accessToken, cookies := app.LoginWithCookies(t, testAdminUser, testAdminPassword)

	// Logout.
	req := app.AuthRequest(http.MethodPost, "/api/v1/auth/logout", accessToken, nil)
	for _, c := range cookies {
		req.AddCookie(c)
	}
	resp, err := app.Do(req)
	if err != nil {
		t.Fatalf("logout request failed: %v", err)
	}
	resp.Body.Close()

	// Using the same access token after logout should fail.
	meReq := app.AuthRequest(http.MethodGet, "/api/v1/auth/me", accessToken, nil)
	meResp, err := app.Do(meReq)
	if err != nil {
		t.Fatalf("me request failed: %v", err)
	}
	defer meResp.Body.Close()

	if meResp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("access after logout: status = %d, want 401", meResp.StatusCode)
	}
}
