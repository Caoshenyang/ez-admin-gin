package auth

import (
	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/middleware"
	"ez-admin-gin/server/internal/platform/datascope"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// MeHandler 负责当前用户相关接口。
type MeHandler struct {
	log *zap.Logger
}

// NewMeHandler 创建当前用户 Handler。
func NewMeHandler(log *zap.Logger) *MeHandler {
	return &MeHandler{
		log: log,
	}
}

type meResponse struct {
	UserID       uint                `json:"user_id"`
	Username     string              `json:"username"`
	DepartmentID uint                `json:"department_id"`
	RoleCodes    []string            `json:"role_codes"`
	IsSuperAdmin bool                `json:"is_super_admin"`
	DataScope    meDataScopeResponse `json:"data_scope"`
}

type meDataScopeResponse struct {
	AllowAll            bool   `json:"allow_all"`
	RequireSelf         bool   `json:"require_self"`
	IncludeDepartment   bool   `json:"include_department"`
	IncludeDeptTree     bool   `json:"include_dept_tree"`
	CustomDepartmentIDs []uint `json:"custom_department_ids"`
}

// Me 返回当前登录用户的基础信息。
func (h *MeHandler) Me(c *gin.Context) {
	if actor, ok := middleware.CurrentActor(c); ok {
		response.Success(c, buildMeResponse(actor))
		return
	}

	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, apperror.Unauthorized("请先登录"), h.log)
		return
	}

	username, _ := middleware.CurrentUsername(c)

	response.Success(c, meResponse{
		UserID:       userID,
		Username:     username,
		DepartmentID: 0,
		RoleCodes:    nil,
		IsSuperAdmin: false,
		DataScope:    meDataScopeResponse{},
	})
}

func buildMeResponse(actor datascope.Actor) meResponse {
	summary := datascope.Merge(actor.Grants, actor.IsSuperAdmin)
	return meResponse{
		UserID:       actor.UserID,
		Username:     actor.Username,
		DepartmentID: actor.DepartmentID,
		RoleCodes:    actor.RoleCodes,
		IsSuperAdmin: actor.IsSuperAdmin,
		DataScope: meDataScopeResponse{
			AllowAll:            summary.AllowAll,
			RequireSelf:         summary.RequireSelf,
			IncludeDepartment:   summary.IncludeDepartment,
			IncludeDeptTree:     summary.IncludeDeptTree,
			CustomDepartmentIDs: summary.CustomDepartmentIDs,
		},
	}
}
