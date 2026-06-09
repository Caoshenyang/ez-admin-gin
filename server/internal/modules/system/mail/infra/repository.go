// Package infra 实现邮件模块的数据访问层。
package infra

import (
	"errors"
	"strings"
	"time"

	maildomain "ez-admin-gin/server/internal/modules/system/mail/domain"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

// Repository 封装系统邮箱、邮件模板和发送日志的数据访问操作。
type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListAccounts(query maildomain.AccountListQuery, page int, pageSize int, status *model.MailAccountStatus) ([]maildomain.AccountEntity, int64, error) {
	queryDB := r.db.Model(&maildomain.AccountEntity{})

	keyword := strings.TrimSpace(query.Keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		queryDB = queryDB.Where("name LIKE ? OR host LIKE ? OR from_email LIKE ?", like, like, like)
	}
	if status != nil {
		queryDB = queryDB.Where("status = ?", *status)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []maildomain.AccountEntity
	if err := queryDB.Order("is_default DESC, id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *Repository) FindAccountByID(db *gorm.DB, accountID uint) (maildomain.AccountEntity, error) {
	var item maildomain.AccountEntity
	err := db.First(&item, accountID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return maildomain.AccountEntity{}, errorsx.NotFound("邮箱账号不存在")
		}
		return maildomain.AccountEntity{}, err
	}
	return item, nil
}

func (r *Repository) FindDefaultAccount() (maildomain.AccountEntity, error) {
	var item maildomain.AccountEntity
	err := r.db.Where("is_default = ?", true).Where("status = ?", model.MailAccountStatusEnabled).Order("id DESC").First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return maildomain.AccountEntity{}, errorsx.NotFound("默认邮箱账号不存在或已禁用")
		}
		return maildomain.AccountEntity{}, err
	}
	return item, nil
}

func (r *Repository) FindEnabledAccountByID(accountID uint) (maildomain.AccountEntity, error) {
	var item maildomain.AccountEntity
	err := r.db.Where("status = ?", model.MailAccountStatusEnabled).First(&item, accountID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return maildomain.AccountEntity{}, errorsx.NotFound("邮箱账号不存在或已禁用")
		}
		return maildomain.AccountEntity{}, err
	}
	return item, nil
}

func (r *Repository) AccountNameExists(db *gorm.DB, name string, excludeID uint) (bool, error) {
	queryDB := db.Unscoped().Model(&maildomain.AccountEntity{}).Where("name = ?", name)
	if excludeID > 0 {
		queryDB = queryDB.Where("id <> ?", excludeID)
	}
	var count int64
	if err := queryDB.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) CreateAccount(db *gorm.DB, item *maildomain.AccountEntity) error {
	return db.Create(item).Error
}

func (r *Repository) UpdateAccount(db *gorm.DB, item *maildomain.AccountEntity, req maildomain.UpdateAccountRequest) error {
	updates := map[string]any{
		"name":       req.Name,
		"host":       req.Host,
		"port":       req.Port,
		"username":   req.Username,
		"from_email": req.FromEmail,
		"from_name":  req.FromName,
		"encryption": req.Encryption,
		"is_default": req.IsDefault,
		"status":     req.Status,
		"remark":     req.Remark,
	}
	if req.Password != "" {
		updates["password"] = req.Password
	}
	if err := db.Model(item).Updates(updates).Error; err != nil {
		return err
	}

	item.Name = req.Name
	item.Host = req.Host
	item.Port = req.Port
	item.Username = req.Username
	if req.Password != "" {
		item.Password = req.Password
	}
	item.FromEmail = req.FromEmail
	item.FromName = req.FromName
	item.Encryption = req.Encryption
	item.IsDefault = req.IsDefault
	item.Status = req.Status
	item.Remark = req.Remark
	return nil
}

func (r *Repository) UpdateAccountStatus(db *gorm.DB, item *maildomain.AccountEntity, status model.MailAccountStatus) error {
	if err := db.Model(item).Update("status", status).Error; err != nil {
		return err
	}
	item.Status = status
	return nil
}

func (r *Repository) UpdateAccountTestResult(db *gorm.DB, item *maildomain.AccountEntity, message string) error {
	now := time.Now()
	if len(message) > 255 {
		message = message[:255]
	}
	if err := db.Model(item).Updates(map[string]any{
		"last_test_at":  now,
		"last_test_msg": message,
	}).Error; err != nil {
		return err
	}
	item.LastTestAt = &now
	item.LastTestMsg = message
	return nil
}

func (r *Repository) ClearDefaultAccounts(db *gorm.DB, excludeID uint) error {
	queryDB := db.Model(&maildomain.AccountEntity{}).Where("is_default = ?", true)
	if excludeID > 0 {
		queryDB = queryDB.Where("id <> ?", excludeID)
	}
	return queryDB.Update("is_default", false).Error
}

func (r *Repository) DeleteAccount(db *gorm.DB, item *maildomain.AccountEntity) error {
	return db.Delete(item).Error
}

func (r *Repository) ListTemplates(query maildomain.TemplateListQuery, page int, pageSize int, status *model.MailTemplateStatus) ([]maildomain.TemplateEntity, int64, error) {
	queryDB := r.db.Model(&maildomain.TemplateEntity{})

	keyword := strings.TrimSpace(query.Keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		queryDB = queryDB.Where("code LIKE ? OR name LIKE ? OR subject LIKE ?", like, like, like)
	}
	if status != nil {
		queryDB = queryDB.Where("status = ?", *status)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []maildomain.TemplateEntity
	if err := queryDB.Order("sort ASC, id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *Repository) FindTemplateByID(db *gorm.DB, templateID uint) (maildomain.TemplateEntity, error) {
	var item maildomain.TemplateEntity
	err := db.First(&item, templateID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return maildomain.TemplateEntity{}, errorsx.NotFound("邮件模板不存在")
		}
		return maildomain.TemplateEntity{}, err
	}
	return item, nil
}

func (r *Repository) FindEnabledTemplateByCode(code string) (maildomain.TemplateEntity, error) {
	var item maildomain.TemplateEntity
	err := r.db.Where("code = ?", code).Where("status = ?", model.MailTemplateStatusEnabled).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return maildomain.TemplateEntity{}, errorsx.NotFound("邮件模板不存在或已禁用")
		}
		return maildomain.TemplateEntity{}, err
	}
	return item, nil
}

func (r *Repository) TemplateCodeExists(db *gorm.DB, code string) (bool, error) {
	var item maildomain.TemplateEntity
	err := db.Unscoped().Where("code = ?", code).First(&item).Error
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, err
}

func (r *Repository) CreateTemplate(db *gorm.DB, item *maildomain.TemplateEntity) error {
	return db.Create(item).Error
}

func (r *Repository) UpdateTemplate(db *gorm.DB, item *maildomain.TemplateEntity, req maildomain.UpdateTemplateRequest) error {
	variables := maildomain.EncodeVariables(req.Variables)
	if err := db.Model(item).Updates(map[string]any{
		"name":      req.Name,
		"subject":   req.Subject,
		"content":   req.Content,
		"is_html":   req.IsHTML,
		"variables": variables,
		"sort":      req.Sort,
		"status":    req.Status,
		"remark":    req.Remark,
	}).Error; err != nil {
		return err
	}
	item.Name = req.Name
	item.Subject = req.Subject
	item.Content = req.Content
	item.IsHTML = req.IsHTML
	item.Variables = variables
	item.Sort = req.Sort
	item.Status = req.Status
	item.Remark = req.Remark
	return nil
}

func (r *Repository) UpdateTemplateStatus(db *gorm.DB, item *maildomain.TemplateEntity, status model.MailTemplateStatus) error {
	if err := db.Model(item).Update("status", status).Error; err != nil {
		return err
	}
	item.Status = status
	return nil
}

func (r *Repository) DeleteTemplate(db *gorm.DB, item *maildomain.TemplateEntity) error {
	return db.Delete(item).Error
}

func (r *Repository) ListLogs(query maildomain.LogListQuery, page int, pageSize int, status *model.MailLogStatus) ([]maildomain.LogEntity, int64, error) {
	queryDB := r.db.Model(&maildomain.LogEntity{})

	keyword := strings.TrimSpace(query.Keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		queryDB = queryDB.Where("subject LIKE ? OR to_emails LIKE ? OR from_email LIKE ?", like, like, like)
	}
	if query.AccountID > 0 {
		queryDB = queryDB.Where("account_id = ?", query.AccountID)
	}
	templateCode := strings.TrimSpace(query.TemplateCode)
	if templateCode != "" {
		queryDB = queryDB.Where("template_code = ?", templateCode)
	}
	if status != nil {
		queryDB = queryDB.Where("status = ?", *status)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []maildomain.LogEntity
	if err := queryDB.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *Repository) CreateLog(db *gorm.DB, item *maildomain.LogEntity) error {
	return db.Create(item).Error
}
