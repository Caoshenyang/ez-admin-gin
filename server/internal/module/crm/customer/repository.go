package customer

import (
	"errors"
	"strings"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
	"ez-admin-gin/server/internal/platform/datascope"

	"gorm.io/gorm"
)

// Repository 负责 CRM 客户的查询和持久化。
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建 CRM 客户仓储。
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// List 返回当前数据范围内的客户分页结果和总数。
func (r *Repository) List(actor datascope.Actor, query ListQuery, page int, pageSize int, status *model.CustomerStatus) ([]View, int64, error) {
	queryDB := r.listBase(actor, query, status)

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []View
	if err := queryDB.
		Select(`
			c.id, c.name, c.contact_name, c.phone, c.level, c.source,
			c.department_id, COALESCE(d.name, '') AS department_name,
			c.owner_user_id, COALESCE(u.username, '') AS owner_username, COALESCE(u.nickname, '') AS owner_nickname,
			c.status, c.remark, c.created_at, c.updated_at
		`).
		Order("c.id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// FindByIDInScope 查询当前数据范围内的客户实体。
func (r *Repository) FindByIDInScope(tx *gorm.DB, actor datascope.Actor, customerID uint) (Entity, error) {
	var item Entity
	err := applyDataScope(r.dbOr(tx).Model(&Entity{}), actor).First(&item, customerID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Entity{}, apperror.NotFound("客户不存在")
		}
		return Entity{}, err
	}

	return item, nil
}

// FindViewByID 查询客户联合视图。
func (r *Repository) FindViewByID(tx *gorm.DB, customerID uint) (View, error) {
	var item View
	err := r.dbOr(tx).
		Table("sys_customer AS c").
		Select(`
			c.id, c.name, c.contact_name, c.phone, c.level, c.source,
			c.department_id, COALESCE(d.name, '') AS department_name,
			c.owner_user_id, COALESCE(u.username, '') AS owner_username, COALESCE(u.nickname, '') AS owner_nickname,
			c.status, c.remark, c.created_at, c.updated_at
		`).
		Joins("LEFT JOIN sys_department AS d ON d.id = c.department_id AND d.deleted_at IS NULL").
		Joins("LEFT JOIN sys_user AS u ON u.id = c.owner_user_id AND u.deleted_at IS NULL").
		Where("c.id = ?", customerID).
		Where("c.deleted_at IS NULL").
		Scan(&item).Error
	if err != nil {
		return View{}, err
	}
	if item.ID == 0 {
		return View{}, gorm.ErrRecordNotFound
	}

	return item, nil
}

// Create 创建客户记录。
func (r *Repository) Create(tx *gorm.DB, item *Entity) error {
	return r.dbOr(tx).Create(item).Error
}

// UpdateBase 更新客户基础字段。
func (r *Repository) UpdateBase(tx *gorm.DB, item *Entity, req UpdateRequest) error {
	item.Name = req.Name
	item.ContactName = req.ContactName
	item.Phone = req.Phone
	item.Level = req.Level
	item.Source = req.Source
	item.Status = req.Status
	item.Remark = req.Remark

	return r.dbOr(tx).Model(item).Updates(map[string]any{
		"name":         item.Name,
		"contact_name": item.ContactName,
		"phone":        item.Phone,
		"level":        item.Level,
		"source":       item.Source,
		"status":       item.Status,
		"remark":       item.Remark,
	}).Error
}

// UpdateStatus 单独更新客户状态。
func (r *Repository) UpdateStatus(tx *gorm.DB, item *Entity, status model.CustomerStatus) error {
	item.Status = status
	return r.dbOr(tx).Model(item).Update("status", status).Error
}

// DepartmentUsable 校验部门是否存在且可用。
func (r *Repository) DepartmentUsable(tx *gorm.DB, departmentID uint) error {
	if departmentID == 0 {
		return nil
	}

	var count int64
	err := r.dbOr(tx).
		Model(&model.Department{}).
		Where("id = ?", departmentID).
		Where("status = ?", model.DepartmentStatusEnabled).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count != 1 {
		return apperror.BadRequest("当前用户所属部门不存在或已禁用")
	}

	return nil
}

// OwnerUsable 校验负责人是否存在且可用。
func (r *Repository) OwnerUsable(tx *gorm.DB, ownerUserID uint) error {
	if ownerUserID == 0 {
		return apperror.BadRequest("当前用户信息不完整，不能创建客户")
	}

	var count int64
	err := r.dbOr(tx).
		Model(&model.User{}).
		Where("id = ?", ownerUserID).
		Where("status = ?", model.UserStatusEnabled).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count != 1 {
		return apperror.BadRequest("当前用户不存在或已禁用")
	}

	return nil
}

func (r *Repository) listBase(actor datascope.Actor, query ListQuery, status *model.CustomerStatus) *gorm.DB {
	queryDB := applyDataScope(
		r.db.Table("sys_customer AS c").
			Joins("LEFT JOIN sys_department AS d ON d.id = c.department_id AND d.deleted_at IS NULL").
			Joins("LEFT JOIN sys_user AS u ON u.id = c.owner_user_id AND u.deleted_at IS NULL").
			Where("c.deleted_at IS NULL"),
		actor,
	)

	keyword := strings.TrimSpace(query.Keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		queryDB = queryDB.Where(
			"c.name LIKE ? OR c.contact_name LIKE ? OR c.phone LIKE ? OR u.username LIKE ? OR u.nickname LIKE ?",
			like, like, like, like, like,
		)
	}

	level := strings.TrimSpace(query.Level)
	if level != "" {
		queryDB = queryDB.Where("c.level = ?", level)
	}

	source := strings.TrimSpace(query.Source)
	if source != "" {
		queryDB = queryDB.Where("c.source = ?", source)
	}

	if status != nil {
		queryDB = queryDB.Where("c.status = ?", *status)
	}

	return queryDB
}

func (r *Repository) dbOr(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}
