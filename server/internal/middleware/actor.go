package middleware

import (
	"ez-admin-gin/server/internal/model"
	"ez-admin-gin/server/internal/platform/datascope"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const currentActorKey = "current_actor"

type actorRoleRow struct {
	RoleID    uint
	Code      string
	DataScope string
}

type actorCustomScopeRow struct {
	RoleID       uint
	DepartmentID uint
}

// LoadActor 在认证通过后加载当前登录人的组织与数据权限上下文。
func LoadActor(db *gorm.DB, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := CurrentUserID(c)
		if !ok {
			c.Next()
			return
		}

		actor, err := buildActor(db, userID)
		if err != nil {
			if log != nil {
				log.Error("load current actor failed", zap.Error(err), zap.Uint("user_id", userID))
			}
			c.Next()
			return
		}

		c.Set(currentActorKey, actor)
		c.Next()
	}
}

// CurrentActor 从 Gin 上下文中读取当前登录人上下文。
func CurrentActor(c *gin.Context) (datascope.Actor, bool) {
	value, ok := c.Get(currentActorKey)
	if !ok {
		return datascope.Actor{}, false
	}

	actor, ok := value.(datascope.Actor)
	return actor, ok
}

func buildActor(db *gorm.DB, userID uint) (datascope.Actor, error) {
	var user model.User
	if err := db.Select("id", "username", "department_id").First(&user, userID).Error; err != nil {
		return datascope.Actor{}, err
	}

	var roleRows []actorRoleRow
	if err := db.
		Table("sys_role AS r").
		Select("r.id AS role_id, r.code, r.data_scope").
		Joins("JOIN sys_user_role AS ur ON ur.role_id = r.id").
		Where("ur.user_id = ?", userID).
		Where("r.status = ?", model.RoleStatusEnabled).
		Where("r.deleted_at IS NULL").
		Order("r.id ASC").
		Scan(&roleRows).Error; err != nil {
		return datascope.Actor{}, err
	}

	roleCodes := make([]string, 0, len(roleRows))
	grants := make([]datascope.Grant, 0, len(roleRows))
	isSuperAdmin := false

	for _, row := range roleRows {
		roleCodes = append(roleCodes, row.Code)
		if row.Code == "super_admin" {
			isSuperAdmin = true
		}
		grants = append(grants, datascope.Grant{
			Scope: datascope.Scope(row.DataScope),
		})
	}

	if err := attachCustomDeptGrants(db, grants, roleRows); err != nil {
		return datascope.Actor{}, err
	}

	return datascope.Actor{
		UserID:       user.ID,
		Username:     user.Username,
		DepartmentID: user.DepartmentID,
		RoleCodes:    roleCodes,
		Grants:       grants,
		IsSuperAdmin: isSuperAdmin,
	}, nil
}

func attachCustomDeptGrants(db *gorm.DB, grants []datascope.Grant, roleRows []actorRoleRow) error {
	roleIndexByID := make(map[uint]int, len(roleRows))
	customRoleIDs := make([]uint, 0, len(roleRows))

	for idx, row := range roleRows {
		roleIndexByID[row.RoleID] = idx
		if datascope.Scope(row.DataScope) == datascope.ScopeCustomDept {
			customRoleIDs = append(customRoleIDs, row.RoleID)
		}
	}

	if len(customRoleIDs) == 0 {
		return nil
	}

	var rows []actorCustomScopeRow
	if err := db.
		Table("sys_role_data_scope").
		Select("role_id, department_id").
		Where("role_id IN ?", customRoleIDs).
		Order("role_id ASC, department_id ASC").
		Scan(&rows).Error; err != nil {
		return err
	}

	for _, row := range rows {
		grantIndex, ok := roleIndexByID[row.RoleID]
		if !ok {
			continue
		}
		grants[grantIndex].DepartmentIDs = append(grants[grantIndex].DepartmentIDs, row.DepartmentID)
	}

	return nil
}
