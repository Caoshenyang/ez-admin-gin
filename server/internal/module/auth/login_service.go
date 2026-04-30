package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
	authnPlatform "ez-admin-gin/server/internal/platform/authn"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// LoginService 负责登录、签发 Token 和登录日志写入。
type LoginService struct {
	repo  *Repository
	token *authnPlatform.Manager
	log   *zap.Logger
}

// NewLoginService 创建登录服务。
func NewLoginService(repo *Repository, token *authnPlatform.Manager, log *zap.Logger) *LoginService {
	return &LoginService{
		repo:  repo,
		token: token,
		log:   log,
	}
}

// Login 执行用户名密码校验并签发 Token。
func (s *LoginService) Login(ctx context.Context, req LoginRequest, ip string, userAgent string) (LoginResponse, error) {
	req, err := NormalizeLoginRequest(req)
	if err != nil {
		s.recordLogin(ctx, 0, "", model.LoginLogStatusFailed, "用户名和密码不能为空", ip, userAgent)
		return LoginResponse{}, err
	}

	user, err := s.repo.FindUserByUsername(req.Username)
	if err != nil {
		if s.repo.IsNotFound(err) {
			s.recordLogin(ctx, 0, req.Username, model.LoginLogStatusFailed, "用户名或密码错误", ip, userAgent)
			return LoginResponse{}, apperror.Unauthorized("用户名或密码错误")
		}

		s.recordLogin(ctx, 0, req.Username, model.LoginLogStatusFailed, "登录失败", ip, userAgent)
		return LoginResponse{}, apperror.Internal("登录失败", err)
	}

	if user.Status != model.UserStatusEnabled {
		s.recordLogin(ctx, user.ID, user.Username, model.LoginLogStatusFailed, "用户已被禁用", ip, userAgent)
		return LoginResponse{}, apperror.Forbidden("用户已被禁用")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		s.recordLogin(ctx, user.ID, user.Username, model.LoginLogStatusFailed, "用户名或密码错误", ip, userAgent)
		return LoginResponse{}, apperror.Unauthorized("用户名或密码错误")
	}

	accessToken, expiresAt, err := s.token.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		s.recordLogin(ctx, user.ID, user.Username, model.LoginLogStatusFailed, "登录失败", ip, userAgent)
		return LoginResponse{}, apperror.Internal("登录失败", err)
	}

	s.recordLogin(ctx, user.ID, user.Username, model.LoginLogStatusSuccess, "登录成功", ip, userAgent)
	return LoginResponse{
		UserID:      user.ID,
		Username:    user.Username,
		Nickname:    user.Nickname,
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresAt:   expiresAt.UTC().Format(time.RFC3339),
	}, nil
}

func (s *LoginService) recordLogin(
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
		var appErr *apperror.Error
		if !errors.As(err, &appErr) {
			s.log.Warn("create login log failed", zap.Error(err))
			return
		}
		s.log.Warn("create login log failed", zap.String("message", appErr.Message))
	}
}
