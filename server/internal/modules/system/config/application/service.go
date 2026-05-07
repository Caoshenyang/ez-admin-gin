package application

import (
	"context"
	"time"

	configdomain "ez-admin-gin/server/internal/modules/system/config/domain"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/pkg/paging"
	"ez-admin-gin/server/internal/platform/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	cachePrefix = "sys_config:"
	cacheTTL    = time.Hour
)

type Service struct {
	tx    ConfigTransactor
	repo  ConfigRepository
	cache ConfigCache
	log   *zap.Logger
}

func NewService(tx ConfigTransactor, repo ConfigRepository, cache ConfigCache, log *zap.Logger) *Service {
	return &Service{tx: tx, repo: repo, cache: cache, log: log}
}

func (s *Service) List(query configdomain.ListQuery) (configdomain.ListResponse, error) {
	page, pageSize := paging.NormalizePage(query.Page, query.PageSize)

	items, total, err := s.repo.List(query, page, pageSize)
	if err != nil {
		return configdomain.ListResponse{}, err
	}

	result := make([]configdomain.Response, 0, len(items))
	for _, item := range items {
		result = append(result, configdomain.BuildResponse(item))
	}

	return configdomain.ListResponse{Items: result, Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *Service) Create(ctx context.Context, req configdomain.CreateRequest) (configdomain.Response, error) {
	req, err := configdomain.NormalizeCreateRequest(req)
	if err != nil {
		return configdomain.Response{}, err
	}

	created := configdomain.Entity{
		GroupCode: req.GroupCode,
		ConfigKey: req.Key,
		Name:      req.Name,
		Value:     req.Value,
		Sort:      req.Sort,
		Status:    req.Status,
		Remark:    req.Remark,
	}

	err = s.tx.WithinTransaction(ctx, func(tx *gorm.DB) error {
		exists, err := s.repo.KeyExists(tx, req.Key)
		if err != nil {
			return err
		}
		if exists {
			return errorsx.BadRequest("配置键已存在")
		}

		return s.repo.Create(tx, &created)
	})
	if err != nil {
		return configdomain.Response{}, err
	}

	s.syncCache(ctx, created)
	return configdomain.BuildResponse(created), nil
}

func (s *Service) Update(ctx context.Context, configID uint, req configdomain.UpdateRequest) (configdomain.Response, error) {
	req, err := configdomain.NormalizeUpdateRequest(req)
	if err != nil {
		return configdomain.Response{}, err
	}

	var updated configdomain.Entity
	err = s.tx.WithinTransaction(ctx, func(tx *gorm.DB) error {
		item, err := s.repo.FindByID(tx, configID)
		if err != nil {
			return err
		}
		if err := s.repo.UpdateBase(tx, &item, req); err != nil {
			return err
		}

		updated = item
		return nil
	})
	if err != nil {
		return configdomain.Response{}, err
	}

	s.syncCache(ctx, updated)
	return configdomain.BuildResponse(updated), nil
}

func (s *Service) UpdateStatus(ctx context.Context, configID uint, status model.SystemConfigStatus) error {
	if !configdomain.ValidStatus(status) {
		return errorsx.BadRequest("配置状态不正确")
	}

	var updated configdomain.Entity
	err := s.tx.WithinTransaction(ctx, func(tx *gorm.DB) error {
		item, err := s.repo.FindByID(tx, configID)
		if err != nil {
			return err
		}
		if err := s.repo.UpdateStatus(tx, &item, status); err != nil {
			return err
		}

		updated = item
		return nil
	})
	if err != nil {
		return err
	}

	s.syncCache(ctx, updated)
	return nil
}

func (s *Service) Value(ctx context.Context, key string) (configdomain.ValueResponse, error) {
	key, err := configdomain.NormalizeKey(key)
	if err != nil {
		return configdomain.ValueResponse{}, err
	}

	if s.cache != nil {
		value, found, err := s.cache.Get(ctx, s.cacheKey(key))
		if found {
			return configdomain.ValueResponse{Key: key, Value: value, Source: "cache"}, nil
		}
		if err != nil {
			s.log.Warn("get system config cache failed", zap.String("key", key), zap.Error(err))
		}
	}

	item, err := s.repo.FindEnabledByKey(key)
	if err != nil {
		return configdomain.ValueResponse{}, err
	}

	s.writeCache(ctx, item)
	return configdomain.ValueResponse{Key: item.ConfigKey, Value: item.Value, Source: "db"}, nil
}

func (s *Service) cacheKey(key string) string {
	return cachePrefix + key
}

func (s *Service) writeCache(ctx context.Context, item model.SystemConfig) {
	if s.cache == nil {
		return
	}
	if err := s.cache.Set(ctx, s.cacheKey(item.ConfigKey), item.Value, cacheTTL); err != nil {
		s.log.Warn("set system config cache failed", zap.String("key", item.ConfigKey), zap.Error(err))
	}
}

func (s *Service) deleteCache(ctx context.Context, key string) {
	if s.cache == nil {
		return
	}
	if err := s.cache.Delete(ctx, s.cacheKey(key)); err != nil {
		s.log.Warn("delete system config cache failed", zap.String("key", key), zap.Error(err))
	}
}

func (s *Service) syncCache(ctx context.Context, item model.SystemConfig) {
	if item.Status == model.SystemConfigStatusEnabled {
		s.writeCache(ctx, item)
		return
	}
	s.deleteCache(ctx, item.ConfigKey)
}
