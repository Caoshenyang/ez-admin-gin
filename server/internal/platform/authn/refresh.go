package authn

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

const (
	refreshTokenPrefix  = "refresh_token:"
	userSessionsPrefix  = "user_sessions:"
	tokenBlacklistPrefix = "token_blacklist:"
)

// RefreshSession holds the data stored alongside a refresh token.
type RefreshSession struct {
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	CreatedAt int64  `json:"created_at"`
}

// RefreshTokenStore manages refresh tokens in Redis with rotation and revocation support.
type RefreshTokenStore struct {
	rdb *goredis.Client
	ttl time.Duration
}

// NewRefreshTokenStore creates a new Redis-backed refresh token store.
func NewRefreshTokenStore(rdb *goredis.Client, ttl time.Duration) *RefreshTokenStore {
	return &RefreshTokenStore{rdb: rdb, ttl: ttl}
}

// Create generates a new refresh token, stores it in Redis, and returns the plaintext token.
func (s *RefreshTokenStore) Create(ctx context.Context, userID uint, username string) (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("generate refresh token: %w", err)
	}
	token := hex.EncodeToString(raw)
	hash := sha256.Sum256([]byte(token))
	key := refreshTokenPrefix + hex.EncodeToString(hash[:])

	session := RefreshSession{
		UserID:    userID,
		Username:  username,
		CreatedAt: time.Now().Unix(),
	}
	data, err := json.Marshal(session)
	if err != nil {
		return "", fmt.Errorf("marshal refresh session: %w", err)
	}

	pipe := s.rdb.Pipeline()
	pipe.Set(ctx, key, data, s.ttl)
	pipe.SAdd(ctx, userSessionsKey(userID), hex.EncodeToString(hash[:]))
	pipe.Expire(ctx, userSessionsKey(userID), s.ttl)
	if _, err := pipe.Exec(ctx); err != nil {
		return "", fmt.Errorf("store refresh token: %w", err)
	}

	return token, nil
}

// Verify checks a refresh token and returns the associated session data.
func (s *RefreshTokenStore) Verify(ctx context.Context, token string) (RefreshSession, error) {
	hash := sha256.Sum256([]byte(token))
	key := refreshTokenPrefix + hex.EncodeToString(hash[:])

	data, err := s.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return RefreshSession{}, fmt.Errorf("refresh token not found or expired")
	}

	var session RefreshSession
	if err := json.Unmarshal(data, &session); err != nil {
		return RefreshSession{}, fmt.Errorf("corrupt refresh session data")
	}

	return session, nil
}

// Revoke invalidates a single refresh token.
func (s *RefreshTokenStore) Revoke(ctx context.Context, token string) error {
	hash := sha256.Sum256([]byte(token))
	hashStr := hex.EncodeToString(hash[:])
	key := refreshTokenPrefix + hashStr

	session, err := s.verifyByKey(ctx, key)
	if err != nil {
		return nil // already revoked, treat as success
	}

	pipe := s.rdb.Pipeline()
	pipe.Del(ctx, key)
	pipe.SRem(ctx, userSessionsKey(session.UserID), hashStr)
	_, _ = pipe.Exec(ctx)
	return nil
}

// RevokeAllForUser invalidates all refresh tokens for a given user.
func (s *RefreshTokenStore) RevokeAllForUser(ctx context.Context, userID uint) error {
	sKey := userSessionsKey(userID)
	hashes, err := s.rdb.SMembers(ctx, sKey).Result()
	if err != nil {
		return nil
	}

	pipe := s.rdb.Pipeline()
	for _, h := range hashes {
		pipe.Del(ctx, refreshTokenPrefix+h)
	}
	pipe.Del(ctx, sKey)
	_, _ = pipe.Exec(ctx)
	return nil
}

// BlacklistAccessToken adds an access token to a short-lived blacklist in Redis.
func (s *RefreshTokenStore) BlacklistAccessToken(ctx context.Context, tokenString string, remainingTTL time.Duration) error {
	hash := sha256.Sum256([]byte(tokenString))
	key := tokenBlacklistPrefix + hex.EncodeToString(hash[:])
	return s.rdb.Set(ctx, key, "1", remainingTTL).Err()
}

// IsBlacklisted checks if an access token has been blacklisted.
func (s *RefreshTokenStore) IsBlacklisted(ctx context.Context, tokenString string) bool {
	hash := sha256.Sum256([]byte(tokenString))
	key := tokenBlacklistPrefix + hex.EncodeToString(hash[:])
	n, err := s.rdb.Exists(ctx, key).Result()
	return err == nil && n > 0
}

func (s *RefreshTokenStore) verifyByKey(ctx context.Context, key string) (RefreshSession, error) {
	data, err := s.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return RefreshSession{}, fmt.Errorf("not found")
	}
	var session RefreshSession
	if err := json.Unmarshal(data, &session); err != nil {
		return RefreshSession{}, fmt.Errorf("corrupt data")
	}
	return session, nil
}

func userSessionsKey(userID uint) string {
	return fmt.Sprintf("%s%d", userSessionsPrefix, userID)
}
