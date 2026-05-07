package application

import (
	"context"
	"errors"
	"strings"
	"time"

	authdomain "ez-admin-gin/server/internal/modules/auth/domain"
	authinfra "ez-admin-gin/server/internal/modules/auth/infra"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	authnPlatform "ez-admin-gin/server/internal/platform/authn"
	"ez-admin-gin/server/internal/platform/model"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type LoginService struct {
	repo  *authinfra.Repository
	token *authnPlatform.Manager
	log   *zap.Logger
}

func NewLoginService(repo *authinfra.Repository, token *authnPlatform.Manager, log *zap.Logger) *LoginService {
	return &LoginService{repo: repo, token: token, log: log}
}

func (s *LoginService) Login(
	ctx context.Context,
	req authdomain.LoginRequest,
	ip string,
	userAgent string,
) (authdomain.LoginResponse, error) {
	req, err := authdomain.NormalizeLoginRequest(req)
	if err != nil {
		s.RecordLogin(ctx, 0, "", model.LoginLogStatusFailed, "用户名和密码不能为空", ip, userAgent)
		return authdomain.LoginResponse{}, err
	}

	user, err := s.repo.FindUserByUsername(req.Username)
	if err != nil {
		if s.repo.IsNotFound(err) {
			s.RecordLogin(ctx, 0, req.Username, model.LoginLogStatusFailed, "用户名或密码错误", ip, userAgent)
			return authdomain.LoginResponse{}, errorsx.Unauthorized("用户名或密码错误")
		}

		s.RecordLogin(ctx, 0, req.Username, model.LoginLogStatusFailed, "登录失败", ip, userAgent)
		return authdomain.LoginResponse{}, errorsx.Internal("登录失败", err)
	}

	if user.Status != model.UserStatusEnabled {
		s.RecordLogin(ctx, user.ID, user.Username, model.LoginLogStatusFailed, "用户已被禁用", ip, userAgent)
		return authdomain.LoginResponse{}, errorsx.Forbidden("用户已被禁用")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		s.RecordLogin(ctx, user.ID, user.Username, model.LoginLogStatusFailed, "用户名或密码错误", ip, userAgent)
		return authdomain.LoginResponse{}, errorsx.Unauthorized("用户名或密码错误")
	}

	accessToken, expiresAt, err := s.token.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		s.RecordLogin(ctx, user.ID, user.Username, model.LoginLogStatusFailed, "登录失败", ip, userAgent)
		return authdomain.LoginResponse{}, errorsx.Internal("登录失败", err)
	}

	s.RecordLogin(ctx, user.ID, user.Username, model.LoginLogStatusSuccess, "登录成功", ip, userAgent)
	return authdomain.LoginResponse{
		UserID:      user.ID,
		Username:    user.Username,
		Nickname:    user.Nickname,
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresAt:   expiresAt.UTC().Format(time.RFC3339),
	}, nil
}

func (s *LoginService) RecordLogin(
	_ context.Context,
	userID uint,
	username string,
	status model.LoginLogStatus,
	message string,
	ip string,
	userAgent string,
) {
	record := model.LoginLog{
		UserID:    userID,
		Username:  strings.TrimSpace(username),
		Status:    status,
		Message:   message,
		IP:        ip,
		UserAgent: userAgent,
	}

	if err := s.repo.CreateLoginLog(&record); err != nil && s.log != nil {
		var appErr *errorsx.Error
		if !errors.As(err, &appErr) {
			s.log.Warn("create login log failed", zap.Error(err))
			return
		}
		s.log.Warn("create login log failed", zap.String("message", appErr.Message))
	}
}
