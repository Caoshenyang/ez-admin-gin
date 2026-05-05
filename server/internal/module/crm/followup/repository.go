package followup

import (
	"errors"
	"strings"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
	"ez-admin-gin/server/internal/platform/datascope"

	"gorm.io/gorm"
)

// Repository 负责 CRM 客户跟进的查询和持久化。
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建 CRM 客户跟进仓储。
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// ListCustomerOptions 返回当前数据范围内可选客户。
func (r *Repository) ListCustomerOptions(actor datascope.Actor, keyword string, limit int) ([]CustomerOption, error) {
	queryDB := r.db.
		Table("sys_customer AS c").
		Select(`
			c.id, c.name, c.department_id, COALESCE(d.name, '') AS department_name,
			c.owner_user_id, COALESCE(u.username, '') AS owner_username, COALESCE(u.nickname, '') AS owner_nickname
		`).
		Joins("LEFT JOIN sys_department AS d ON d.id = c.department_id AND d.deleted_at IS NULL").
		Joins("LEFT JOIN sys_user AS u ON u.id = c.owner_user_id AND u.deleted_at IS NULL").
		Where("c.deleted_at IS NULL").
		Where("c.status = ?", model.CustomerStatusEnabled).
		Scopes(datascope.UserQueryScope(r.db, actor, "c.department_id", "c.owner_user_id"))

	trimmedKeyword := strings.TrimSpace(keyword)
	if trimmedKeyword != "" {
		like := "%" + trimmedKeyword + "%"
		queryDB = queryDB.Where(
			"c.name LIKE ? OR c.contact_name LIKE ? OR c.phone LIKE ? OR u.username LIKE ? OR u.nickname LIKE ?",
			like, like, like, like, like,
		)
	}

	var items []CustomerOption
	if err := queryDB.Order("c.id DESC").Limit(limit).Scan(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// List 返回当前数据范围内的客户跟进分页结果和总数。
func (r *Repository) List(actor datascope.Actor, query ListQuery, page int, pageSize int, status *model.CustomerFollowUpStatus) ([]View, int64, error) {
	queryDB := r.listBase(actor, query, status)

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []View
	if err := queryDB.
		Select(`
			f.id, f.customer_id, COALESCE(c.name, '') AS customer_name,
			f.department_id, COALESCE(d.name, '') AS department_name,
			f.owner_user_id, COALESCE(u.username, '') AS owner_username, COALESCE(u.nickname, '') AS owner_nickname,
			f.follow_type, f.subject, f.content, f.result, f.next_follow_at, f.status, f.created_at, f.updated_at
		`).
		Order("f.id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// FindByIDInScope 查询当前数据范围内的客户跟进实体。
func (r *Repository) FindByIDInScope(tx *gorm.DB, actor datascope.Actor, followUpID uint) (Entity, error) {
	var item Entity
	err := applyDataScope(r.dbOr(tx).Table("sys_customer_followup AS f"), actor).
		Select("f.*").
		Where("f.deleted_at IS NULL").
		First(&item, "f.id = ?", followUpID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Entity{}, apperror.NotFound("客户跟进不存在")
		}
		return Entity{}, err
	}

	return item, nil
}

// FindViewByID 查询客户跟进联合视图。
func (r *Repository) FindViewByID(tx *gorm.DB, followUpID uint) (View, error) {
	var item View
	err := r.dbOr(tx).
		Table("sys_customer_followup AS f").
		Select(`
			f.id, f.customer_id, COALESCE(c.name, '') AS customer_name,
			f.department_id, COALESCE(d.name, '') AS department_name,
			f.owner_user_id, COALESCE(u.username, '') AS owner_username, COALESCE(u.nickname, '') AS owner_nickname,
			f.follow_type, f.subject, f.content, f.result, f.next_follow_at, f.status, f.created_at, f.updated_at
		`).
		Joins("LEFT JOIN sys_customer AS c ON c.id = f.customer_id AND c.deleted_at IS NULL").
		Joins("LEFT JOIN sys_department AS d ON d.id = f.department_id AND d.deleted_at IS NULL").
		Joins("LEFT JOIN sys_user AS u ON u.id = f.owner_user_id AND u.deleted_at IS NULL").
		Where("f.id = ?", followUpID).
		Where("f.deleted_at IS NULL").
		Scan(&item).Error
	if err != nil {
		return View{}, err
	}
	if item.ID == 0 {
		return View{}, gorm.ErrRecordNotFound
	}

	return item, nil
}

// FindCustomerInScope 查询当前数据范围内的客户实体。
func (r *Repository) FindCustomerInScope(tx *gorm.DB, actor datascope.Actor, customerID uint) (model.Customer, error) {
	var item model.Customer
	err := r.dbOr(tx).
		Model(&model.Customer{}).
		Scopes(datascope.UserQueryScope(r.dbOr(tx), actor, "department_id", "owner_user_id")).
		First(&item, customerID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Customer{}, apperror.NotFound("客户不存在")
		}
		return model.Customer{}, err
	}

	return item, nil
}

// Create 创建客户跟进记录。
func (r *Repository) Create(tx *gorm.DB, item *Entity) error {
	return r.dbOr(tx).Create(item).Error
}

// UpdateBase 更新客户跟进基础字段。
func (r *Repository) UpdateBase(tx *gorm.DB, item *Entity, req UpdateRequest) error {
	item.FollowType = req.FollowType
	item.Subject = req.Subject
	item.Content = req.Content
	item.Result = req.Result
	item.NextFollowAt = req.NextFollowAt
	item.Status = req.Status

	return r.dbOr(tx).Model(item).Updates(map[string]any{
		"follow_type":    item.FollowType,
		"subject":        item.Subject,
		"content":        item.Content,
		"result":         item.Result,
		"next_follow_at": item.NextFollowAt,
		"status":         item.Status,
	}).Error
}

// UpdateStatus 单独更新客户跟进状态。
func (r *Repository) UpdateStatus(tx *gorm.DB, item *Entity, status model.CustomerFollowUpStatus) error {
	item.Status = status
	return r.dbOr(tx).Model(item).Update("status", status).Error
}

func (r *Repository) listBase(actor datascope.Actor, query ListQuery, status *model.CustomerFollowUpStatus) *gorm.DB {
	queryDB := applyDataScope(
		r.db.Table("sys_customer_followup AS f").
			Joins("LEFT JOIN sys_customer AS c ON c.id = f.customer_id AND c.deleted_at IS NULL").
			Joins("LEFT JOIN sys_department AS d ON d.id = f.department_id AND d.deleted_at IS NULL").
			Joins("LEFT JOIN sys_user AS u ON u.id = f.owner_user_id AND u.deleted_at IS NULL").
			Where("f.deleted_at IS NULL"),
		actor,
	)

	keyword := strings.TrimSpace(query.Keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		queryDB = queryDB.Where(
			"f.subject LIKE ? OR f.content LIKE ? OR f.result LIKE ? OR c.name LIKE ? OR u.username LIKE ? OR u.nickname LIKE ?",
			like, like, like, like, like, like,
		)
	}

	followType := strings.TrimSpace(query.FollowType)
	if followType != "" {
		queryDB = queryDB.Where("f.follow_type = ?", followType)
	}

	if query.CustomerID > 0 {
		queryDB = queryDB.Where("f.customer_id = ?", query.CustomerID)
	}

	if status != nil {
		queryDB = queryDB.Where("f.status = ?", *status)
	}

	return queryDB
}

func (r *Repository) dbOr(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}
