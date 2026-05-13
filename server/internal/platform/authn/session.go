package authn

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

const (
	sessionTokenBytes  = 32
	sessionTokenPrefix = "session:refresh:"
	sessionUserPrefix  = "session:user:"
)

// SessionStore manages refresh token sessions.
type SessionStore interface {
	Create(ctx context.Context, userID uint) (token string, err error)
	Validate(ctx context.Context, token string) (uint, error)
	Revoke(ctx context.Context, token string) error
	RevokeAllForUser(ctx context.Context, userID uint) error
}

// RedisSessionStore implements SessionStore using Redis.
type RedisSessionStore struct {
	rdb    *goredis.Client
	ttl    time.Duration
	now    func() time.Time
}

// NewRedisSessionStore creates a Redis-backed session store.
func NewRedisSessionStore(rdb *goredis.Client, ttl time.Duration) *RedisSessionStore {
	return &RedisSessionStore{rdb: rdb, ttl: ttl, now: time.Now}
}

func (s *RedisSessionStore) Create(ctx context.Context, userID uint) (string, error) {
	raw, err := generateRandomToken(sessionTokenBytes)
	if err != nil {
		return "", fmt.Errorf("generate refresh token: %w", err)
	}

	hash := sha256Hex(raw)
	sessionKey := sessionTokenPrefix + hash
	userKey := sessionUserPrefix + strconv.FormatUint(uint64(userID), 10)

	pipe := s.rdb.Pipeline()
	pipe.Set(ctx, sessionKey, strconv.FormatUint(uint64(userID), 10), s.ttl)
	pipe.SAdd(ctx, userKey, hash)
	pipe.Expire(ctx, userKey, s.ttl+time.Hour)
	if _, err := pipe.Exec(ctx); err != nil {
		return "", fmt.Errorf("store refresh token: %w", err)
	}

	return raw, nil
}

func (s *RedisSessionStore) Validate(ctx context.Context, token string) (uint, error) {
	hash := sha256Hex(token)
	val, err := s.rdb.Get(ctx, sessionTokenPrefix+hash).Result()
	if err != nil {
		return 0, fmt.Errorf("refresh token not found or expired")
	}
	userID, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid session data")
	}
	return uint(userID), nil
}

func (s *RedisSessionStore) Revoke(ctx context.Context, token string) error {
	hash := sha256Hex(token)
	s.rdb.Del(ctx, sessionTokenPrefix+hash)
	return nil
}

func (s *RedisSessionStore) RevokeAllForUser(ctx context.Context, userID uint) error {
	userKey := sessionUserPrefix + strconv.FormatUint(uint64(userID), 10)
	hashes, err := s.rdb.SMembers(ctx, userKey).Result()
	if err != nil {
		return err
	}
	if len(hashes) == 0 {
		return nil
	}
	keys := make([]string, len(hashes))
	for i, h := range hashes {
		keys[i] = sessionTokenPrefix + h
	}
	s.rdb.Del(ctx, keys...)
	s.rdb.Del(ctx, userKey)
	return nil
}

func generateRandomToken(size int) (string, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func sha256Hex(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}
