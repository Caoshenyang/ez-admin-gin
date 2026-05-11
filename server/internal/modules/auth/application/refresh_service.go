package application

import (
	"context"

	authdomain "ez-admin-gin/server/internal/modules/auth/domain"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
)

// RefreshService handles refresh token validation, rotation, and new access token issuance.
type RefreshService struct {
	session SessionStore
	token   TokenIssuer
	repo    LoginRepository
}

func NewRefreshService(session SessionStore, token TokenIssuer, repo LoginRepository) *RefreshService {
	return &RefreshService{session: session, token: token, repo: repo}
}

// Refresh validates the refresh token, rotates it, and issues a new access token.
func (s *RefreshService) Refresh(ctx context.Context, refreshToken string) (LoginResult, error) {
	userID, err := s.session.Validate(ctx, refreshToken)
	if err != nil {
		return LoginResult{}, errorsx.Unauthorized("refresh token 无效或已过期")
	}

	user, err := s.repo.FindUserByIDSimple(userID)
	if err != nil {
		return LoginResult{}, errorsx.Unauthorized("用户不存在")
	}

	if user.Status != 1 {
		return LoginResult{}, errorsx.Forbidden("用户已被禁用")
	}

	// Rotate: revoke old token and create new one.
	_ = s.session.Revoke(ctx, refreshToken)

	newRefreshToken, err := s.session.Create(ctx, user.ID)
	if err != nil {
		return LoginResult{}, errorsx.Internal("创建会话失败", err)
	}

	accessToken, expiresAt, err := s.token.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return LoginResult{}, errorsx.Internal("签发令牌失败", err)
	}

	return LoginResult{
		Response: authdomain.LoginResponse{
			UserID:      user.ID,
			Username:    user.Username,
			Nickname:    user.Nickname,
			AccessToken: accessToken,
			TokenType:   "Bearer",
			ExpiresAt:   expiresAt.UTC().Format("2006-01-02T15:04:05Z"),
		},
		RefreshToken: newRefreshToken,
	}, nil
}

// Logout revokes the refresh token session.
func (s *RefreshService) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return nil
	}
	return s.session.Revoke(ctx, refreshToken)
}
