package config

import (
	"context"
	"time"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"

	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	cachePrefix = "sys_config:"
	cacheTTL    = time.Hour
)

// Service 负责系统配置的业务规则、事务边界和缓存同步。
type Service struct {
	db    *gorm.DB
	repo  *Repository
	redis *goredis.Client
	log   *zap.Logger
}

// NewService 创建配置服务。
func NewService(db *gorm.DB, repo *Repository, redis *goredis.Client, log *zap.Logger) *Service {
	return &Service{
		db:    db,
		repo:  repo,
		redis: redis,
		log:   log,
	}
}

// List 返回系统配置分页结果。
func (s *Service) List(query ListQuery) (ListResponse, error) {
	page, pageSize := NormalizePage(query.Page, query.PageSize)

	items, total, err := s.repo.List(query, page, pageSize)
	if err != nil {
		return ListResponse{}, err
	}

	result := make([]Response, 0, len(items))
	for _, item := range items {
		result = append(result, BuildResponse(item))
	}

	return ListResponse{
		Items:    result,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// Create 创建系统配置。
func (s *Service) Create(ctx context.Context, req CreateRequest) (Response, error) {
	req, err := NormalizeCreateRequest(req)
	if err != nil {
		return Response{}, err
	}

	created := Entity{
		GroupCode: req.GroupCode,
		ConfigKey: req.Key,
		Name:      req.Name,
		Value:     req.Value,
		Sort:      req.Sort,
		Status:    req.Status,
		Remark:    req.Remark,
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		exists, err := s.repo.KeyExists(tx, req.Key)
		if err != nil {
			return err
		}
		if exists {
			return apperror.BadRequest("配置键已存在")
		}

		return s.repo.Create(tx, &created)
	})
	if err != nil {
		return Response{}, err
	}

	s.syncCache(ctx, created)
	return BuildResponse(created), nil
}

// Update 编辑系统配置。
func (s *Service) Update(ctx context.Context, configID uint, req UpdateRequest) (Response, error) {
	req, err := NormalizeUpdateRequest(req)
	if err != nil {
		return Response{}, err
	}

	var updated Entity
	err = s.db.Transaction(func(tx *gorm.DB) error {
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
		return Response{}, err
	}

	s.syncCache(ctx, updated)
	return BuildResponse(updated), nil
}

// UpdateStatus 单独修改配置状态。
func (s *Service) UpdateStatus(ctx context.Context, configID uint, status model.SystemConfigStatus) error {
	if !ValidStatus(status) {
		return apperror.BadRequest("配置状态不正确")
	}

	var updated Entity
	err := s.db.Transaction(func(tx *gorm.DB) error {
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

// Value 按 key 读取启用中的配置值，并优先命中 Redis。
func (s *Service) Value(ctx context.Context, key string) (ValueResponse, error) {
	key, err := NormalizeKey(key)
	if err != nil {
		return ValueResponse{}, err
	}

	if s.redis != nil {
		value, err := s.redis.Get(ctx, s.cacheKey(key)).Result()
		if err == nil {
			return ValueResponse{
				Key:    key,
				Value:  value,
				Source: "cache",
			}, nil
		}
		if err != nil && err != goredis.Nil {
			s.log.Warn("get system config cache failed", zap.String("key", key), zap.Error(err))
		}
	}

	item, err := s.repo.FindEnabledByKey(key)
	if err != nil {
		return ValueResponse{}, err
	}

	s.writeCache(ctx, item)
	return ValueResponse{
		Key:    item.ConfigKey,
		Value:  item.Value,
		Source: "db",
	}, nil
}

func (s *Service) cacheKey(key string) string {
	return cachePrefix + key
}

func (s *Service) writeCache(ctx context.Context, item model.SystemConfig) {
	if s.redis == nil {
		return
	}

	if err := s.redis.Set(ctx, s.cacheKey(item.ConfigKey), item.Value, cacheTTL).Err(); err != nil {
		s.log.Warn("set system config cache failed", zap.String("key", item.ConfigKey), zap.Error(err))
	}
}

func (s *Service) deleteCache(ctx context.Context, key string) {
	if s.redis == nil {
		return
	}

	if err := s.redis.Del(ctx, s.cacheKey(key)).Err(); err != nil {
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
