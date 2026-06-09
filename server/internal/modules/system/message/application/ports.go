package application

import (
	messagedomain "ez-admin-gin/server/internal/modules/system/message/domain"
	"ez-admin-gin/server/internal/platform/database"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

type MessageTransactor = database.Transactor

type MessageRepository interface {
	ListTemplates(query messagedomain.TemplateListQuery, page int, pageSize int, status *model.MessageConfigStatus, templateType *model.MessageTemplateType) ([]messagedomain.TemplateEntity, int64, error)
	FindTemplateByID(db *gorm.DB, templateID uint) (messagedomain.TemplateEntity, error)
	TemplateCodeExists(db *gorm.DB, code string) (bool, error)
	CreateTemplate(db *gorm.DB, item *messagedomain.TemplateEntity) error
	UpdateTemplateBase(db *gorm.DB, item *messagedomain.TemplateEntity, req messagedomain.UpdateTemplateRequest) error
	UpdateTemplateStatus(db *gorm.DB, item *messagedomain.TemplateEntity, status model.MessageConfigStatus) error

	ListReminders(query messagedomain.ReminderListQuery, page int, pageSize int, status *model.MessageConfigStatus, receiverType *model.MessageReceiverType) ([]messagedomain.ReminderListItem, int64, error)
	FindReminderByID(db *gorm.DB, reminderID uint) (messagedomain.ReminderEntity, error)
	ReminderCodeExists(db *gorm.DB, code string) (bool, error)
	CreateReminder(db *gorm.DB, item *messagedomain.ReminderEntity) error
	UpdateReminderBase(db *gorm.DB, item *messagedomain.ReminderEntity, req messagedomain.UpdateReminderRequest) error
	UpdateReminderStatus(db *gorm.DB, item *messagedomain.ReminderEntity, status model.MessageConfigStatus) error
}
