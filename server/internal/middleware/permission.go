package middleware

import (
	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
	"ez-admin-gin/server/internal/permission"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Permission 根据当前用户角色判断接口访问权限。
func Permission(db *gorm.DB, enforcer *permission.Enforcer, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleCodes, err := permissionRoleCodes(c, db)
		if err != nil {
			response.Error(c, apperror.Internal("权限校验失败", err), log)
			c.Abort()
			return
		}

		if len(roleCodes) == 0 {
			response.Error(c, apperror.Forbidden("没有权限访问"), log)
			c.Abort()
			return
		}

		obj := c.FullPath()
		if obj == "" {
			obj = c.Request.URL.Path
		}
		act := c.Request.Method

		for _, roleCode := range roleCodes {
			allowed, err := enforcer.Enforce(roleCode, obj, act)
			if err != nil {
				response.Error(c, apperror.Internal("权限校验失败", err), log)
				c.Abort()
				return
			}

			if allowed {
				c.Next()
				return
			}
		}

		response.Error(c, apperror.Forbidden("没有权限访问"), log)
		c.Abort()
	}
}

func permissionRoleCodes(c *gin.Context, db *gorm.DB) ([]string, error) {
	if actor, ok := CurrentActor(c); ok {
		return actor.RoleCodes, nil
	}

	userID, ok := CurrentUserID(c)
	if !ok {
		return nil, apperror.Unauthorized("请先登录")
	}

	return currentRoleCodes(db, userID)
}

// currentRoleCodes 查询当前用户拥有的启用角色编码。
func currentRoleCodes(db *gorm.DB, userID uint) ([]string, error) {
	var roleCodes []string
	err := db.
		Table("sys_role AS r").
		Select("r.code").
		Joins("JOIN sys_user_role AS ur ON ur.role_id = r.id").
		Where("ur.user_id = ?", userID).
		Where("r.status = ?", model.RoleStatusEnabled).
		Where("r.deleted_at IS NULL").
		Pluck("r.code", &roleCodes).Error
	if err != nil {
		return nil, err
	}

	return roleCodes, nil
}
