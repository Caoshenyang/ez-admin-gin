package application

import (
	"context"
	"time"

	authdomain "ez-admin-gin/server/internal/modules/auth/domain"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	authnPlatform "ez-admin-gin/server/internal/platform/authn"
)

// RefreshService handles refresh token rotation.
type RefreshService struct {
	refreshStore *authnPlatform.RefreshTokenStore
	token        TokenIssuer
}

// NewRefreshService creates a new RefreshService.
func NewRefreshService(refreshStore *authnPlatform.RefreshTokenStore, token TokenIssuer) *RefreshService {
	return &RefreshService{refreshStore: refreshStore, token: token}
}

// Refresh rotates the refresh token: verify old → revoke old → issue new pair.
func (s *RefreshService) Refresh(ctx context.Context, oldRefreshToken string) (authdomain.LoginResponse, string, error) {
	session, err := s.refreshStore.Verify(ctx, oldRefreshToken)
	if err != nil {
		return authdomain.LoginResponse{}, "", errorsx.Unauthorized("refresh token 无效或已过期")
	}

	// Revoke the old refresh token (rotation).
	_ = s.refreshStore.Revoke(ctx, oldRefreshToken)

	// Issue new access token.
	accessToken, expiresAt, err := s.token.GenerateAccessToken(session.UserID, session.Username)
	if err != nil {
		return authdomain.LoginResponse{}, "", err
	}

	// Issue new refresh token.
	newRefreshToken, err := s.token.GenerateRefreshToken(ctx, session.UserID, session.Username)
	if err != nil {
		return authdomain.LoginResponse{}, "", err
	}

	return authdomain.LoginResponse{
		UserID:      session.UserID,
		Username:    session.Username,
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresAt:   expiresAt.UTC().Format(time.RFC3339),
	}, newRefreshToken, nil
}
