package datascope

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// Scope 表示一条数据权限规则。
type Scope string

const (
	// ScopeAll 表示允许查看全部数据。
	ScopeAll Scope = "all"
	// ScopeDept 表示仅查看本部门数据。
	ScopeDept Scope = "dept"
	// ScopeDeptAndChildren 表示查看本部门及子部门数据。
	ScopeDeptAndChildren Scope = "dept_and_children"
	// ScopeSelf 表示仅查看本人数据。
	ScopeSelf Scope = "self"
	// ScopeCustomDept 表示按角色显式授权的部门范围查看数据。
	ScopeCustomDept Scope = "custom_dept"
)

// Grant 表示角色授予给当前用户的一条数据权限。
type Grant struct {
	Scope         Scope
	DepartmentIDs []uint
}

// Actor 表示后续业务查询可复用的当前登录人上下文。
type Actor struct {
	UserID       uint
	Username     string
	DepartmentID uint
	RoleCodes    []string
	Grants       []Grant
	IsSuperAdmin bool
}

// Summary 表示将多角色数据权限合并后的结果。
type Summary struct {
	AllowAll            bool
	RequireSelf         bool
	IncludeDepartment   bool
	IncludeDeptTree     bool
	CustomDepartmentIDs []uint
}

// Merge 将多角色数据权限按并集规则压缩为一份摘要。
func Merge(grants []Grant, isSuperAdmin bool) Summary {
	if isSuperAdmin {
		return Summary{AllowAll: true}
	}

	summary := Summary{}
	customDeptSet := make(map[uint]struct{})

	for _, grant := range grants {
		switch grant.Scope {
		case ScopeAll:
			return Summary{AllowAll: true}
		case ScopeDept:
			summary.IncludeDepartment = true
		case ScopeDeptAndChildren:
			summary.IncludeDeptTree = true
		case ScopeSelf:
			summary.RequireSelf = true
		case ScopeCustomDept:
			for _, departmentID := range grant.DepartmentIDs {
				if departmentID == 0 {
					continue
				}
				customDeptSet[departmentID] = struct{}{}
			}
		}
	}

	if len(customDeptSet) > 0 {
		summary.CustomDepartmentIDs = make([]uint, 0, len(customDeptSet))
		for departmentID := range customDeptSet {
			summary.CustomDepartmentIDs = append(summary.CustomDepartmentIDs, departmentID)
		}
	}

	return summary
}

// UserQueryScope 为“用户列表/详情”这类既有部门归属、又可退化到本人数据的资源生成查询作用域。
func UserQueryScope(db *gorm.DB, actor Actor, departmentColumn string, ownerColumn string) func(*gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		summary := Merge(actor.Grants, actor.IsSuperAdmin)
		if summary.AllowAll {
			return tx
		}

		conditions := make([]string, 0, 4)
		args := make([]any, 0, 8)

		if summary.IncludeDepartment && actor.DepartmentID != 0 && departmentColumn != "" {
			conditions = append(conditions, fmt.Sprintf("%s = ?", departmentColumn))
			args = append(args, actor.DepartmentID)
		}

		if summary.IncludeDeptTree && actor.DepartmentID != 0 && departmentColumn != "" {
			departmentIDs, err := expandDepartmentTree(db, []uint{actor.DepartmentID})
			if err != nil {
				tx.AddError(err)
				return tx
			}
			if len(departmentIDs) > 0 {
				conditions = append(conditions, fmt.Sprintf("%s IN ?", departmentColumn))
				args = append(args, departmentIDs)
			}
		}

		if len(summary.CustomDepartmentIDs) > 0 && departmentColumn != "" {
			conditions = append(conditions, fmt.Sprintf("%s IN ?", departmentColumn))
			args = append(args, summary.CustomDepartmentIDs)
		}

		if summary.RequireSelf && actor.UserID != 0 && ownerColumn != "" {
			conditions = append(conditions, fmt.Sprintf("%s = ?", ownerColumn))
			args = append(args, actor.UserID)
		}

		if len(conditions) == 0 {
			return tx.Where("1 = 0")
		}

		return tx.Where("("+strings.Join(conditions, " OR ")+")", args...)
	}
}

// DepartmentQueryScope 为“部门树/部门列表”这类以部门自身为资源的数据生成查询作用域。
// 当角色范围是 self 时，这里退化为“当前用户所在部门”，避免部门管理完全看不到自己的归属部门。
func DepartmentQueryScope(db *gorm.DB, actor Actor, departmentIDColumn string) func(*gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		summary := Merge(actor.Grants, actor.IsSuperAdmin)
		if summary.AllowAll {
			return tx
		}

		departmentIDs, err := accessibleDepartmentIDs(db, actor, summary)
		if err != nil {
			tx.AddError(err)
			return tx
		}
		if len(departmentIDs) == 0 || departmentIDColumn == "" {
			return tx.Where("1 = 0")
		}

		return tx.Where(fmt.Sprintf("%s IN ?", departmentIDColumn), departmentIDs)
	}
}

func accessibleDepartmentIDs(db *gorm.DB, actor Actor, summary Summary) ([]uint, error) {
	idSet := make(map[uint]struct{})

	if summary.IncludeDepartment && actor.DepartmentID != 0 {
		idSet[actor.DepartmentID] = struct{}{}
	}

	if summary.IncludeDeptTree && actor.DepartmentID != 0 {
		departmentIDs, err := expandDepartmentTree(db, []uint{actor.DepartmentID})
		if err != nil {
			return nil, err
		}
		for _, departmentID := range departmentIDs {
			if departmentID == 0 {
				continue
			}
			idSet[departmentID] = struct{}{}
		}
	}

	for _, departmentID := range summary.CustomDepartmentIDs {
		if departmentID == 0 {
			continue
		}
		idSet[departmentID] = struct{}{}
	}

	if summary.RequireSelf && actor.DepartmentID != 0 {
		idSet[actor.DepartmentID] = struct{}{}
	}

	result := make([]uint, 0, len(idSet))
	for departmentID := range idSet {
		result = append(result, departmentID)
	}

	return result, nil
}

func expandDepartmentTree(db *gorm.DB, departmentIDs []uint) ([]uint, error) {
	if len(departmentIDs) == 0 {
		return nil, nil
	}

	idSet := make(map[uint]struct{}, len(departmentIDs))
	for _, departmentID := range departmentIDs {
		if departmentID == 0 {
			continue
		}
		idSet[departmentID] = struct{}{}
	}

	likeParts := make([]string, 0, len(idSet))
	args := make([]any, 0, len(idSet))
	for departmentID := range idSet {
		likeParts = append(likeParts, "ancestors = ? OR ancestors LIKE ? OR ancestors LIKE ? OR ancestors LIKE ?")
		idText := fmt.Sprintf("%d", departmentID)
		args = append(args, idText, idText+",%", "%,"+idText, "%,"+idText+",%")
	}

	query := db.Table("sys_department").Select("id")
	if len(likeParts) > 0 {
		query = query.Where(strings.Join(likeParts, " OR "), args...)
	}

	var descendants []uint
	if err := query.Find(&descendants).Error; err != nil {
		return nil, err
	}

	for _, departmentID := range descendants {
		idSet[departmentID] = struct{}{}
	}

	result := make([]uint, 0, len(idSet))
	for departmentID := range idSet {
		result = append(result, departmentID)
	}

	return result, nil
}
