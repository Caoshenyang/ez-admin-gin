// Package infra 实现消息提醒配置的数据访问层。
package infra

import (
	"errors"
	"strings"

	messagedomain "ez-admin-gin/server/internal/modules/system/message/domain"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

// Repository 封装消息模板和提醒规则表的数据访问操作。
type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ListTemplates 按关键词、类型和状态分页查询消息模板。
func (r *Repository) ListTemplates(query messagedomain.TemplateListQuery, page int, pageSize int, status *model.MessageConfigStatus, templateType *model.MessageTemplateType) ([]messagedomain.TemplateEntity, int64, error) {
	queryDB := r.db.Model(&messagedomain.TemplateEntity{})

	keyword := strings.TrimSpace(query.Keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		queryDB = queryDB.Where("code LIKE ? OR name LIKE ? OR title LIKE ?", like, like, like)
	}
	if status != nil {
		queryDB = queryDB.Where("status = ?", *status)
	}
	if templateType != nil {
		queryDB = queryDB.Where("type = ?", *templateType)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []messagedomain.TemplateEntity
	if err := queryDB.Order("sort ASC, id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// FindTemplateByID 在指定事务中按主键查找消息模板，不存在时返回 NotFound 错误。
func (r *Repository) FindTemplateByID(db *gorm.DB, templateID uint) (messagedomain.TemplateEntity, error) {
	var item messagedomain.TemplateEntity
	err := db.First(&item, templateID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return messagedomain.TemplateEntity{}, errorsx.NotFound("消息模板不存在")
		}
		return messagedomain.TemplateEntity{}, err
	}
	return item, nil
}

// TemplateCodeExists 检查指定模板编码是否已存在（包含已软删除的记录）。
func (r *Repository) TemplateCodeExists(db *gorm.DB, code string) (bool, error) {
	var item messagedomain.TemplateEntity
	err := db.Unscoped().Where("code = ?", code).First(&item).Error
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, err
}

// CreateTemplate 在指定事务中插入一条消息模板。
func (r *Repository) CreateTemplate(db *gorm.DB, item *messagedomain.TemplateEntity) error {
	return db.Create(item).Error
}

// UpdateTemplateBase 更新消息模板基本字段。
func (r *Repository) UpdateTemplateBase(db *gorm.DB, item *messagedomain.TemplateEntity, req messagedomain.UpdateTemplateRequest) error {
	if err := db.Model(item).Updates(map[string]any{
		"name":      req.Name,
		"title":     req.Title,
		"content":   req.Content,
		"type":      req.Type,
		"variables": req.Variables,
		"sort":      req.Sort,
		"status":    req.Status,
		"is_system": req.IsSystem,
		"remark":    req.Remark,
	}).Error; err != nil {
		return err
	}

	item.Name = req.Name
	item.Title = req.Title
	item.Content = req.Content
	item.Type = req.Type
	item.Variables = req.Variables
	item.Sort = req.Sort
	item.Status = req.Status
	item.IsSystem = req.IsSystem
	item.Remark = req.Remark
	return nil
}

// UpdateTemplateStatus 更新消息模板状态字段。
func (r *Repository) UpdateTemplateStatus(db *gorm.DB, item *messagedomain.TemplateEntity, status model.MessageConfigStatus) error {
	if err := db.Model(item).Update("status", status).Error; err != nil {
		return err
	}
	item.Status = status
	return nil
}

// ListReminders 按关键词、触发事件、模板、接收人类型和状态分页查询提醒规则。
func (r *Repository) ListReminders(query messagedomain.ReminderListQuery, page int, pageSize int, status *model.MessageConfigStatus, receiverType *model.MessageReceiverType) ([]messagedomain.ReminderListItem, int64, error) {
	queryDB := r.db.Model(&messagedomain.ReminderEntity{})

	keyword := strings.TrimSpace(query.Keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		queryDB = queryDB.Where("code LIKE ? OR name LIKE ? OR trigger_event LIKE ?", like, like, like)
	}
	triggerEvent := strings.TrimSpace(query.TriggerEvent)
	if triggerEvent != "" {
		queryDB = queryDB.Where("trigger_event = ?", triggerEvent)
	}
	if query.TemplateID != 0 {
		queryDB = queryDB.Where("template_id = ?", query.TemplateID)
	}
	if status != nil {
		queryDB = queryDB.Where("status = ?", *status)
	}
	if receiverType != nil {
		queryDB = queryDB.Where("receiver_type = ?", *receiverType)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var reminders []messagedomain.ReminderEntity
	if err := queryDB.Order("sort ASC, id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&reminders).Error; err != nil {
		return nil, 0, err
	}

	templateIDs := make([]uint, 0, len(reminders))
	seen := map[uint]bool{}
	for _, item := range reminders {
		if item.TemplateID != 0 && !seen[item.TemplateID] {
			templateIDs = append(templateIDs, item.TemplateID)
			seen[item.TemplateID] = true
		}
	}

	templates := map[uint]messagedomain.TemplateEntity{}
	if len(templateIDs) > 0 {
		var rows []messagedomain.TemplateEntity
		if err := r.db.Where("id IN ?", templateIDs).Find(&rows).Error; err != nil {
			return nil, 0, err
		}
		for _, item := range rows {
			templates[item.ID] = item
		}
	}

	items := make([]messagedomain.ReminderListItem, 0, len(reminders))
	for _, reminder := range reminders {
		items = append(items, messagedomain.ReminderListItem{
			Reminder: reminder,
			Template: templates[reminder.TemplateID],
		})
	}

	return items, total, nil
}

// FindReminderByID 在指定事务中按主键查找提醒规则，不存在时返回 NotFound 错误。
func (r *Repository) FindReminderByID(db *gorm.DB, reminderID uint) (messagedomain.ReminderEntity, error) {
	var item messagedomain.ReminderEntity
	err := db.First(&item, reminderID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return messagedomain.ReminderEntity{}, errorsx.NotFound("提醒规则不存在")
		}
		return messagedomain.ReminderEntity{}, err
	}
	return item, nil
}

// ReminderCodeExists 检查指定提醒编码是否已存在（包含已软删除的记录）。
func (r *Repository) ReminderCodeExists(db *gorm.DB, code string) (bool, error) {
	var item messagedomain.ReminderEntity
	err := db.Unscoped().Where("code = ?", code).First(&item).Error
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, err
}

// CreateReminder 在指定事务中插入一条提醒规则。
func (r *Repository) CreateReminder(db *gorm.DB, item *messagedomain.ReminderEntity) error {
	return db.Create(item).Error
}

// UpdateReminderBase 更新提醒规则基本字段。
func (r *Repository) UpdateReminderBase(db *gorm.DB, item *messagedomain.ReminderEntity, req messagedomain.UpdateReminderRequest) error {
	if err := db.Model(item).Updates(map[string]any{
		"name":            req.Name,
		"trigger_event":   req.TriggerEvent,
		"template_id":     req.TemplateID,
		"channels":        req.Channels,
		"receiver_type":   req.ReceiverType,
		"receiver_values": req.ReceiverValues,
		"advance_minutes": req.AdvanceMinutes,
		"link_url":        req.LinkURL,
		"sort":            req.Sort,
		"status":          req.Status,
		"is_system":       req.IsSystem,
		"remark":          req.Remark,
	}).Error; err != nil {
		return err
	}

	item.Name = req.Name
	item.TriggerEvent = req.TriggerEvent
	item.TemplateID = req.TemplateID
	item.Channels = req.Channels
	item.ReceiverType = req.ReceiverType
	item.ReceiverValues = req.ReceiverValues
	item.AdvanceMinutes = req.AdvanceMinutes
	item.LinkURL = req.LinkURL
	item.Sort = req.Sort
	item.Status = req.Status
	item.IsSystem = req.IsSystem
	item.Remark = req.Remark
	return nil
}

// UpdateReminderStatus 更新提醒规则状态字段。
func (r *Repository) UpdateReminderStatus(db *gorm.DB, item *messagedomain.ReminderEntity, status model.MessageConfigStatus) error {
	if err := db.Model(item).Update("status", status).Error; err != nil {
		return err
	}
	item.Status = status
	return nil
}
