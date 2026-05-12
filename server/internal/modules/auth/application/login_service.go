// Package application 实现 auth 模块的业务逻辑：登录、菜单查询、账户管理和仪表盘。
package application

import (
	"context"
	"errors"
	"strings"
	"time"

	authdomain "ez-admin-gin/server/internal/modules/auth/domain"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/model"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// AccountLockChecker checks if an account is locked and records/clears attempts.
type AccountLockChecker interface {
	IsLocked(ctx context.Context, username string) bool
	RecordFailure(ctx context.Context, username string)
	ClearAttempts(ctx context.Context, username string)
}

// LoginService 处理用户登录认证与登录日志记录。
type LoginService struct {
	repo       LoginRepository
	token      TokenIssuer
	lockChecker AccountLockChecker
	log        *zap.Logger
}

func NewLoginService(repo LoginRepository, token TokenIssuer, lockChecker AccountLockChecker, log *zap.Logger) *LoginService {
	return &LoginService{repo: repo, token: token, lockChecker: lockChecker, log: log}
}

// Login 校验用户名密码，成功后签发 Access Token 并记录登录日志。
func (s *LoginService) Login(
	ctx context.Context,
	req authdomain.LoginRequest,
	ip string,
	userAgent string,
) (authdomain.LoginResponse, string, error) {
	req, err := authdomain.NormalizeLoginRequest(req)
	if err != nil {
		s.RecordLogin(ctx, 0, "", model.LoginLogStatusFailed, "用户名和密码不能为空", ip, userAgent)
		return authdomain.LoginResponse{}, "", err
	}

	if s.lockChecker != nil && s.lockChecker.IsLocked(ctx, req.Username) {
		s.RecordLogin(ctx, 0, req.Username, model.LoginLogStatusFailed, "账户已锁定，请稍后再试", ip, userAgent)
		return authdomain.LoginResponse{}, "", errorsx.TooManyRequests("账户已锁定，请稍后再试")
	}

	user, err := s.repo.FindUserByUsername(req.Username)
	if err != nil {
		if s.repo.IsNotFound(err) {
			s.RecordLogin(ctx, 0, req.Username, model.LoginLogStatusFailed, "用户名或密码错误", ip, userAgent)
			return authdomain.LoginResponse{}, "", errorsx.Unauthorized("用户名或密码错误")
		}

		s.RecordLogin(ctx, 0, req.Username, model.LoginLogStatusFailed, "登录失败", ip, userAgent)
		return authdomain.LoginResponse{}, "", errorsx.Internal("登录失败", err)
	}

	if user.Status != model.UserStatusEnabled {
		s.RecordLogin(ctx, user.ID, user.Username, model.LoginLogStatusFailed, "用户已被禁用", ip, userAgent)
		return authdomain.LoginResponse{}, "", errorsx.Forbidden("用户已被禁用")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		s.RecordLogin(ctx, user.ID, user.Username, model.LoginLogStatusFailed, "用户名或密码错误", ip, userAgent)
		if s.lockChecker != nil {
			s.lockChecker.RecordFailure(ctx, req.Username)
		}
		return authdomain.LoginResponse{}, "", errorsx.Unauthorized("用户名或密码错误")
	}

	accessToken, expiresAt, err := s.token.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		s.RecordLogin(ctx, user.ID, user.Username, model.LoginLogStatusFailed, "登录失败", ip, userAgent)
		return authdomain.LoginResponse{}, "", errorsx.Internal("登录失败", err)
	}

	refreshToken, err := s.token.GenerateRefreshToken(ctx, user.ID, user.Username)
	if err != nil {
		s.RecordLogin(ctx, user.ID, user.Username, model.LoginLogStatusFailed, "登录失败", ip, userAgent)
		return authdomain.LoginResponse{}, "", errorsx.Internal("登录失败", err)
	}

	if s.lockChecker != nil {
		s.lockChecker.ClearAttempts(ctx, req.Username)
	}
	s.RecordLogin(ctx, user.ID, user.Username, model.LoginLogStatusSuccess, "登录成功", ip, userAgent)
	return authdomain.LoginResponse{
		UserID:      user.ID,
		Username:    user.Username,
		Nickname:    user.Nickname,
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresAt:   expiresAt.UTC().Format(time.RFC3339),
	}, refreshToken, nil
}

// RecordLogin 写入一条登录日志，失败仅 warn 不中断主流程。
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
