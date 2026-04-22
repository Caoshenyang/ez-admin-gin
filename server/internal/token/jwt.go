package token

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"ez-admin-gin/server/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// ErrInvalidToken 表示 Token 无效、过期或签名不正确。
	ErrInvalidToken = errors.New("invalid token")
)

// Claims 是写入 access_token 的业务载荷。
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Manager 负责生成和解析 access_token。
type Manager struct {
	secret         []byte
	issuer         string
	accessTokenTTL time.Duration
	now            func() time.Time
}

// NewManager 根据配置创建 Token 管理器。
func NewManager(cfg config.AuthConfig) (*Manager, error) {
	secret := strings.TrimSpace(cfg.JWTSecret)
	if len(secret) < 32 {
		return nil, fmt.Errorf("jwt secret must be at least 32 characters")
	}

	if cfg.AccessTokenTTL <= 0 {
		return nil, fmt.Errorf("access token ttl must be greater than 0")
	}

	issuer := strings.TrimSpace(cfg.Issuer)
	if issuer == "" {
		return nil, fmt.Errorf("jwt issuer cannot be empty")
	}

	return &Manager{
		secret:         []byte(secret),
		issuer:         issuer,
		accessTokenTTL: time.Duration(cfg.AccessTokenTTL) * time.Second,
		now:            time.Now,
	}, nil
}

// GenerateAccessToken 生成访问令牌，并返回令牌过期时间。
func (m *Manager) GenerateAccessToken(userID uint, username string) (string, time.Time, error) {
	now := m.now()
	expiresAt := now.Add(m.accessTokenTTL)

	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   fmt.Sprintf("%d", userID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign access token: %w", err)
	}

	return tokenString, expiresAt, nil
}

// ParseAccessToken 解析并校验访问令牌。
func (m *Manager) ParseAccessToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	parsedToken, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (any, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, ErrInvalidToken
			}

			return m.secret, nil
		},
		jwt.WithIssuer(m.issuer),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if !parsedToken.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
