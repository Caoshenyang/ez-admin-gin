package menu

import (
	"errors"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"

	"gorm.io/gorm"
)

// Repository 负责菜单模块的查询拼装和持久化。
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建菜单仓储。
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// List 返回完整菜单树所需的节点列表。
func (r *Repository) List() ([]model.Menu, error) {
	var items []model.Menu
	if err := r.db.Order("sort ASC, id ASC").Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

// FindByID 查询单个菜单节点。
func (r *Repository) FindByID(db *gorm.DB, menuID uint) (model.Menu, error) {
	var item model.Menu
	err := db.First(&item, menuID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Menu{}, apperror.NotFound("菜单不存在")
		}
		return model.Menu{}, err
	}

	return item, nil
}

// CodeExists 判断菜单编码是否已存在。
func (r *Repository) CodeExists(db *gorm.DB, code string) (bool, error) {
	var item model.Menu
	err := db.Unscoped().Where("code = ?", code).First(&item).Error
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return false, err
}

// ParentUsable 校验父级菜单是否可用，并守住目录/菜单/按钮之间的层级约束。
func (r *Repository) ParentUsable(db *gorm.DB, parentID uint, menuType model.MenuType, currentID uint) error {
	if parentID == 0 {
		if menuType != model.MenuTypeDirectory {
			return apperror.BadRequest("根节点只能是目录")
		}
		return nil
	}

	if parentID == currentID {
		return apperror.BadRequest("父级菜单不能选择自己")
	}

	var parent model.Menu
	if err := db.First(&parent, parentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.BadRequest("父级菜单不存在")
		}
		return err
	}

	if parent.Type == model.MenuTypeButton {
		return apperror.BadRequest("按钮下面不能再添加子节点")
	}
	if menuType == model.MenuTypeButton && parent.Type != model.MenuTypeMenu {
		return apperror.BadRequest("按钮只能挂在菜单下面")
	}

	return nil
}

// CanDelete 校验菜单是否可删除。
func (r *Repository) CanDelete(db *gorm.DB, menuID uint) error {
	var childCount int64
	if err := db.Model(&model.Menu{}).Where("parent_id = ?", menuID).Count(&childCount).Error; err != nil {
		return err
	}
	if childCount > 0 {
		return apperror.BadRequest("请先删除子菜单")
	}

	var roleMenuCount int64
	if err := db.Model(&model.RoleMenu{}).Where("menu_id = ?", menuID).Count(&roleMenuCount).Error; err != nil {
		return err
	}
	if roleMenuCount > 0 {
		return apperror.BadRequest("菜单已分配给角色，不能删除")
	}

	return nil
}

// Create 创建菜单节点。
func (r *Repository) Create(db *gorm.DB, item *model.Menu) error {
	return db.Create(item).Error
}

// UpdateBase 更新菜单基础字段。
func (r *Repository) UpdateBase(db *gorm.DB, item *model.Menu, req UpdateRequest) error {
	if err := db.Model(item).Updates(map[string]any{
		"parent_id":  req.ParentID,
		"type":       req.Type,
		"title":      req.Title,
		"path":       req.Path,
		"component":  req.Component,
		"icon":       req.Icon,
		"sort":       req.Sort,
		"status":     req.Status,
		"remark":     req.Remark,
	}).Error; err != nil {
		return err
	}

	item.ParentID = req.ParentID
	item.Type = req.Type
	item.Title = req.Title
	item.Path = req.Path
	item.Component = req.Component
	item.Icon = req.Icon
	item.Sort = req.Sort
	item.Status = req.Status
	item.Remark = req.Remark
	return nil
}

// UpdateStatus 更新菜单状态。
func (r *Repository) UpdateStatus(db *gorm.DB, item *model.Menu, status model.MenuStatus) error {
	if err := db.Model(item).Update("status", status).Error; err != nil {
		return err
	}
	item.Status = status
	return nil
}

// Delete 删除菜单节点。
func (r *Repository) Delete(db *gorm.DB, item *model.Menu) error {
	return db.Delete(item).Error
}
