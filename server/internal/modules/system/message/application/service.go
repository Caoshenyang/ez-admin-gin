// Package application 实现消息提醒配置的业务逻辑。
package application

import (
	"context"

	messagedomain "ez-admin-gin/server/internal/modules/system/message/domain"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/pkg/paging"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

// Service 封装消息模板和提醒规则的系统级配置逻辑。
type Service struct {
	tx   MessageTransactor
	repo MessageRepository
}

func NewService(tx MessageTransactor, repo MessageRepository) *Service {
	return &Service{tx: tx, repo: repo}
}

// ListTemplates 分页查询消息模板。
func (s *Service) ListTemplates(query messagedomain.TemplateListQuery) (messagedomain.TemplateListResponse, error) {
	page, pageSize := paging.NormalizePage(query.Page, query.PageSize)
	status, err := messagedomain.NormalizeStatusFilter(query.Status)
	if err != nil {
		return messagedomain.TemplateListResponse{}, err
	}
	templateType, err := messagedomain.NormalizeTemplateTypeFilter(query.Type)
	if err != nil {
		return messagedomain.TemplateListResponse{}, err
	}

	items, total, err := s.repo.ListTemplates(query, page, pageSize, status, templateType)
	if err != nil {
		return messagedomain.TemplateListResponse{}, err
	}

	result := make([]messagedomain.TemplateResponse, 0, len(items))
	for _, item := range items {
		result = append(result, messagedomain.BuildTemplateResponse(item))
	}

	return messagedomain.TemplateListResponse{Items: result, Total: total, Page: page, PageSize: pageSize}, nil
}

// CreateTemplate 创建消息模板。
func (s *Service) CreateTemplate(req messagedomain.CreateTemplateRequest) (messagedomain.TemplateResponse, error) {
	req, err := messagedomain.NormalizeCreateTemplateRequest(req)
	if err != nil {
		return messagedomain.TemplateResponse{}, err
	}

	created := messagedomain.TemplateEntity{
		Code:      req.Code,
		Name:      req.Name,
		Title:     req.Title,
		Content:   req.Content,
		Type:      req.Type,
		Variables: req.Variables,
		Sort:      req.Sort,
		Status:    req.Status,
		IsSystem:  req.IsSystem,
		Remark:    req.Remark,
	}

	if err := s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		exists, err := s.repo.TemplateCodeExists(tx, req.Code)
		if err != nil {
			return err
		}
		if exists {
			return errorsx.BadRequest("模板编码已存在")
		}
		return s.repo.CreateTemplate(tx, &created)
	}); err != nil {
		return messagedomain.TemplateResponse{}, err
	}

	return messagedomain.BuildTemplateResponse(created), nil
}

// UpdateTemplate 更新消息模板。
func (s *Service) UpdateTemplate(templateID uint, req messagedomain.UpdateTemplateRequest) (messagedomain.TemplateResponse, error) {
	req, err := messagedomain.NormalizeUpdateTemplateRequest(req)
	if err != nil {
		return messagedomain.TemplateResponse{}, err
	}

	var updated messagedomain.TemplateEntity
	if err := s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		item, err := s.repo.FindTemplateByID(tx, templateID)
		if err != nil {
			return err
		}
		if err := s.repo.UpdateTemplateBase(tx, &item, req); err != nil {
			return err
		}
		updated = item
		return nil
	}); err != nil {
		return messagedomain.TemplateResponse{}, err
	}

	return messagedomain.BuildTemplateResponse(updated), nil
}

// UpdateTemplateStatus 更新消息模板状态。
func (s *Service) UpdateTemplateStatus(templateID uint, status model.MessageConfigStatus) error {
	status, err := messagedomain.NormalizeStatus(status, false)
	if err != nil {
		return err
	}

	return s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		item, err := s.repo.FindTemplateByID(tx, templateID)
		if err != nil {
			return err
		}
		return s.repo.UpdateTemplateStatus(tx, &item, status)
	})
}

// ListReminders 分页查询消息提醒规则。
func (s *Service) ListReminders(query messagedomain.ReminderListQuery) (messagedomain.ReminderListResponse, error) {
	page, pageSize := paging.NormalizePage(query.Page, query.PageSize)
	status, err := messagedomain.NormalizeStatusFilter(query.Status)
	if err != nil {
		return messagedomain.ReminderListResponse{}, err
	}
	receiverType, err := messagedomain.NormalizeReceiverTypeFilter(query.ReceiverType)
	if err != nil {
		return messagedomain.ReminderListResponse{}, err
	}

	items, total, err := s.repo.ListReminders(query, page, pageSize, status, receiverType)
	if err != nil {
		return messagedomain.ReminderListResponse{}, err
	}

	result := make([]messagedomain.ReminderResponse, 0, len(items))
	for _, item := range items {
		result = append(result, messagedomain.BuildReminderResponse(item.Reminder, item.Template))
	}

	return messagedomain.ReminderListResponse{Items: result, Total: total, Page: page, PageSize: pageSize}, nil
}

// CreateReminder 创建消息提醒规则。
func (s *Service) CreateReminder(req messagedomain.CreateReminderRequest) (messagedomain.ReminderResponse, error) {
	req, err := messagedomain.NormalizeCreateReminderRequest(req)
	if err != nil {
		return messagedomain.ReminderResponse{}, err
	}

	created := messagedomain.ReminderEntity{
		Code:           req.Code,
		Name:           req.Name,
		TriggerEvent:   req.TriggerEvent,
		TemplateID:     req.TemplateID,
		Channels:       req.Channels,
		ReceiverType:   req.ReceiverType,
		ReceiverValues: req.ReceiverValues,
		AdvanceMinutes: req.AdvanceMinutes,
		LinkURL:        req.LinkURL,
		Sort:           req.Sort,
		Status:         req.Status,
		IsSystem:       req.IsSystem,
		Remark:         req.Remark,
	}

	var template messagedomain.TemplateEntity
	if err := s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		exists, err := s.repo.ReminderCodeExists(tx, req.Code)
		if err != nil {
			return err
		}
		if exists {
			return errorsx.BadRequest("提醒编码已存在")
		}

		template, err = s.repo.FindTemplateByID(tx, req.TemplateID)
		if err != nil {
			return err
		}
		return s.repo.CreateReminder(tx, &created)
	}); err != nil {
		return messagedomain.ReminderResponse{}, err
	}

	return messagedomain.BuildReminderResponse(created, template), nil
}

// UpdateReminder 更新消息提醒规则。
func (s *Service) UpdateReminder(reminderID uint, req messagedomain.UpdateReminderRequest) (messagedomain.ReminderResponse, error) {
	req, err := messagedomain.NormalizeUpdateReminderRequest(req)
	if err != nil {
		return messagedomain.ReminderResponse{}, err
	}

	var updated messagedomain.ReminderEntity
	var template messagedomain.TemplateEntity
	if err := s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		item, err := s.repo.FindReminderByID(tx, reminderID)
		if err != nil {
			return err
		}
		template, err = s.repo.FindTemplateByID(tx, req.TemplateID)
		if err != nil {
			return err
		}
		if err := s.repo.UpdateReminderBase(tx, &item, req); err != nil {
			return err
		}
		updated = item
		return nil
	}); err != nil {
		return messagedomain.ReminderResponse{}, err
	}

	return messagedomain.BuildReminderResponse(updated, template), nil
}

// UpdateReminderStatus 更新消息提醒规则状态。
func (s *Service) UpdateReminderStatus(reminderID uint, status model.MessageConfigStatus) error {
	status, err := messagedomain.NormalizeStatus(status, false)
	if err != nil {
		return err
	}

	return s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		item, err := s.repo.FindReminderByID(tx, reminderID)
		if err != nil {
			return err
		}
		return s.repo.UpdateReminderStatus(tx, &item, status)
	})
}
