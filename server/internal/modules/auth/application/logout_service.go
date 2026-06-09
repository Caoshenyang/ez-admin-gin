package application

import (
	"context"

	authnPlatform "ez-admin-gin/server/internal/platform/authn"
)

// LogoutService handles token revocation on logout.
type LogoutService struct {
	refreshStore *authnPlatform.RefreshTokenStore
	tokenManager *authnPlatform.Manager
}

// NewLogoutService creates a new LogoutService.
func NewLogoutService(refreshStore *authnPlatform.RefreshTokenStore, tokenManager *authnPlatform.Manager) *LogoutService {
	return &LogoutService{refreshStore: refreshStore, tokenManager: tokenManager}
}

// Logout revokes the refresh token and blacklists the access token.
func (s *LogoutService) Logout(ctx context.Context, refreshToken string, accessToken string) error {
	if refreshToken != "" && s.refreshStore != nil {
		_ = s.refreshStore.Revoke(ctx, refreshToken)
	}

	if accessToken != "" && s.refreshStore != nil {
		ttl := s.tokenManager.AccessTokenTTL()
		_ = s.refreshStore.BlacklistAccessToken(ctx, accessToken, ttl)
	}

	return nil
}
