package application

import (
	"context"

	notidomain "ez-admin-gin/server/internal/modules/system/notification/domain"
	"ez-admin-gin/server/internal/pkg/paging"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

// Service 封装通知的业务逻辑。
type Service struct {
	tx   NotificationTransactor
	repo NotificationRepository
}

func NewService(tx NotificationTransactor, repo NotificationRepository) *Service {
	return &Service{tx: tx, repo: repo}
}

// List 分页查询当前用户通知。
func (s *Service) List(userID uint, query notidomain.ListQuery) (notidomain.ListResponse, error) {
	page, pageSize := paging.NormalizePage(query.Page, query.PageSize)
	notiType, err := notidomain.NormalizeTypeFilter(query.Type)
	if err != nil {
		return notidomain.ListResponse{}, err
	}
	isRead, err := notidomain.NormalizeReadFilter(query.IsRead)
	if err != nil {
		return notidomain.ListResponse{}, err
	}

	items, total, err := s.repo.List(userID, page, pageSize, notiType, isRead)
	if err != nil {
		return notidomain.ListResponse{}, err
	}

	result := make([]notidomain.Response, 0, len(items))
	for _, item := range items {
		result = append(result, notidomain.BuildResponse(item))
	}

	return notidomain.ListResponse{Items: result, Total: total, Page: page, PageSize: pageSize}, nil
}

// UnreadCount 获取当前用户未读通知数。
func (s *Service) UnreadCount(userID uint) (int64, error) {
	return s.repo.UnreadCount(userID)
}

// MarkRead 标记指定通知已读。
func (s *Service) MarkRead(userID uint, ids []uint64) error {
	if len(ids) == 0 {
		return nil
	}
	return s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		return s.repo.MarkRead(tx, userID, ids)
	})
}

// MarkAllRead 标记全部通知已读。
func (s *Service) MarkAllRead(userID uint) error {
	return s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		return s.repo.MarkAllRead(tx, userID)
	})
}

// CreateForUser 创建一条通知（供系统内部调用）。
func (s *Service) CreateForUser(userID uint, title, content string, notiType model.NotificationType, extra model.JSONMap) (notidomain.Response, error) {
	title, content, notiType, err := notidomain.NormalizeCreateRequest(title, content, notiType)
	if err != nil {
		return notidomain.Response{}, err
	}

	item := notidomain.Entity{
		UserID:  userID,
		Type:    notiType,
		Title:   title,
		Content: content,
		Extra:   extra,
	}

	if err := s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		return s.repo.Create(tx, &item)
	}); err != nil {
		return notidomain.Response{}, err
	}

	return notidomain.BuildResponse(item), nil
}

// FindByID 按 ID 查找通知。
func (s *Service) FindByID(db *gorm.DB, id uint64) (notidomain.Entity, error) {
	return s.repo.FindByID(db, id)
}

// BuildNotificationMessage 构建推送消息（供 Hub 使用）。
func BuildNotificationMessage(resp notidomain.Response) map[string]any {
	return map[string]any{
		"type": "notification",
		"data": resp,
	}
}

// BuildUnreadCountMessage 构建未读数消息。
func BuildUnreadCountMessage(count int64) map[string]any {
	return map[string]any{
		"type": "unread_count",
		"data": map[string]any{"count": count},
	}
}
