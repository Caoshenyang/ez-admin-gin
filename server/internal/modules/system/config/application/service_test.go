package application

import (
	"context"
	"testing"
	"time"

	configdomain "ez-admin-gin/server/internal/modules/system/config/domain"
	"ez-admin-gin/server/internal/platform/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type configTestTransactor struct{}

func (configTestTransactor) WithinTransaction(_ context.Context, fn func(tx *gorm.DB) error) error {
	return fn(nil)
}

type configTestRepo struct {
	findEnabledCalls int
	item             model.SystemConfig
}

func (r *configTestRepo) List(query configdomain.ListQuery, page int, pageSize int) ([]model.SystemConfig, int64, error) {
	return nil, 0, nil
}
func (r *configTestRepo) FindByID(db *gorm.DB, configID uint) (model.SystemConfig, error) {
	return r.item, nil
}
func (r *configTestRepo) FindEnabledByKey(key string) (model.SystemConfig, error) {
	r.findEnabledCalls++
	return r.item, nil
}
func (r *configTestRepo) KeyExists(db *gorm.DB, key string) (bool, error)    { return false, nil }
func (r *configTestRepo) Create(db *gorm.DB, item *model.SystemConfig) error { return nil }
func (r *configTestRepo) UpdateBase(db *gorm.DB, item *model.SystemConfig, req configdomain.UpdateRequest) error {
	return nil
}
func (r *configTestRepo) UpdateStatus(db *gorm.DB, item *model.SystemConfig, status model.SystemConfigStatus) error {
	item.Status = status
	r.item = *item
	return nil
}

type configTestCache struct {
	value       string
	found       bool
	getCalls    int
	deleteCalls []string
}

func (c *configTestCache) Get(ctx context.Context, key string) (string, bool, error) {
	c.getCalls++
	return c.value, c.found, nil
}
func (c *configTestCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	c.value = value
	c.found = true
	return nil
}
func (c *configTestCache) Delete(ctx context.Context, key string) error {
	c.deleteCalls = append(c.deleteCalls, key)
	return nil
}

func TestValuePrefersCache(t *testing.T) {
	repo := &configTestRepo{}
	cache := &configTestCache{value: "cached-value", found: true}
	service := NewService(configTestTransactor{}, repo, cache, zap.NewNop())

	result, err := service.Value(context.Background(), "feature:flag")
	if err != nil {
		t.Fatalf("Value returned error: %v", err)
	}
	if result.Source != "cache" || result.Value != "cached-value" {
		t.Fatalf("expected cache result, got %+v", result)
	}
	if repo.findEnabledCalls != 0 {
		t.Fatalf("expected repo not to be called on cache hit")
	}
}

func TestUpdateStatusDeletesCacheForDisabledConfig(t *testing.T) {
	repo := &configTestRepo{
		item: model.SystemConfig{
			ConfigKey: "feature:flag",
			Status:    model.SystemConfigStatusEnabled,
		},
	}
	cache := &configTestCache{}
	service := NewService(configTestTransactor{}, repo, cache, zap.NewNop())

	if err := service.UpdateStatus(context.Background(), 1, model.SystemConfigStatusDisabled); err != nil {
		t.Fatalf("UpdateStatus returned error: %v", err)
	}
	if len(cache.deleteCalls) != 1 || cache.deleteCalls[0] != "sys_config:feature:flag" {
		t.Fatalf("expected cache delete for config key, got %v", cache.deleteCalls)
	}
}
