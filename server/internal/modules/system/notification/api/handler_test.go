package api

import (
	"testing"

	"ez-admin-gin/server/internal/platform/authn"
	"ez-admin-gin/server/internal/platform/config"
)

func TestNewWSHandlerStoresTokenManager(t *testing.T) {
	t.Parallel()

	tokenManager, err := authn.NewManager(config.AuthConfig{
		JWTSecret:      "12345678901234567890123456789012",
		Issuer:         "ez-admin-test",
		AccessTokenTTL: 3600,
	})
	if err != nil {
		t.Fatalf("NewManager() error = %v", err)
	}

	handler := NewWSHandler(nil, nil, tokenManager, nil)
	if handler.token != tokenManager {
		t.Fatal("expected token manager to be stored on WSHandler")
	}
}
