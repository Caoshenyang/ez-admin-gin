package domain

import (
	"strings"
	"time"

	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/datascope"
	"ez-admin-gin/server/internal/platform/model"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	UserID      uint   `json:"user_id"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresAt   string `json:"expires_at"`
}

type MeResponse struct {
	UserID       uint              `json:"user_id"`
	Username     string            `json:"username"`
	DepartmentID uint              `json:"department_id"`
	RoleCodes    []string          `json:"role_codes"`
	IsSuperAdmin bool              `json:"is_super_admin"`
	DataScope    MeDataScopeResult `json:"data_scope"`
}

type AccountProfileResponse struct {
	UserID         uint              `json:"user_id"`
	Username       string            `json:"username"`
	Nickname       string            `json:"nickname"`
	DepartmentID   uint              `json:"department_id"`
	DepartmentName string            `json:"department_name"`
	Status         model.UserStatus  `json:"status"`
	RoleCodes      []string          `json:"role_codes"`
	IsSuperAdmin   bool              `json:"is_super_admin"`
	DataScope      MeDataScopeResult `json:"data_scope"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

type UpdateAccountProfileRequest struct {
	Nickname string `json:"nickname"`
}

type UpdateAccountPasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type MeDataScopeResult struct {
	AllowAll            bool   `json:"allow_all"`
	RequireSelf         bool   `json:"require_self"`
	IncludeDepartment   bool   `json:"include_department"`
	IncludeDeptTree     bool   `json:"include_dept_tree"`
	CustomDepartmentIDs []uint `json:"custom_department_ids"`
}

type MenuResponse struct {
	ID        uint           `json:"id"`
	ParentID  uint           `json:"parent_id"`
	Type      model.MenuType `json:"type"`
	Code      string         `json:"code"`
	Title     string         `json:"title"`
	Path      string         `json:"path"`
	Component string         `json:"component"`
	Icon      string         `json:"icon"`
	Sort      int            `json:"sort"`
	Children  []MenuResponse `json:"children,omitempty"`
}

type DashboardCurrentUser struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

type DashboardHealth struct {
	Env      string `json:"env"`
	Database string `json:"database"`
	Redis    string `json:"redis"`
}

type DashboardMetrics struct {
	UserTotal               int64 `json:"user_total"`
	EnabledUserTotal        int64 `json:"enabled_user_total"`
	EnabledRoleTotal        int64 `json:"enabled_role_total"`
	ConfigTotal             int64 `json:"config_total"`
	NoticeTotal             int64 `json:"notice_total"`
	FileTotal               int64 `json:"file_total"`
	TodayOperationTotal     int64 `json:"today_operation_total"`
	TodayRiskOperationTotal int64 `json:"today_risk_operation_total"`
	TodayLoginFailedTotal   int64 `json:"today_login_failed_total"`
}

type DashboardOperationItem struct {
	ID         uint      `json:"id"`
	Username   string    `json:"username"`
	Method     string    `json:"method"`
	Path       string    `json:"path"`
	StatusCode int       `json:"status_code"`
	Success    bool      `json:"success"`
	LatencyMs  int64     `json:"latency_ms"`
	CreatedAt  time.Time `json:"created_at"`
}

type DashboardLoginItem struct {
	ID        uint                 `json:"id"`
	Username  string               `json:"username"`
	Status    model.LoginLogStatus `json:"status"`
	Message   string               `json:"message"`
	IP        string               `json:"ip"`
	CreatedAt time.Time            `json:"created_at"`
}

type DashboardNoticeItem struct {
	ID        uint               `json:"id"`
	Title     string             `json:"title"`
	Status    model.NoticeStatus `json:"status"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type DashboardResponse struct {
	CurrentUser      DashboardCurrentUser     `json:"current_user"`
	Health           DashboardHealth          `json:"health"`
	Metrics          DashboardMetrics         `json:"metrics"`
	RecentOperations []DashboardOperationItem `json:"recent_operations"`
	RecentLogins     []DashboardLoginItem     `json:"recent_logins"`
	LatestNotices    []DashboardNoticeItem    `json:"latest_notices"`
}

func NormalizeLoginRequest(req LoginRequest) (LoginRequest, error) {
	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || req.Password == "" {
		return LoginRequest{}, errorsx.BadRequest("用户名和密码不能为空")
	}
	return req, nil
}

func BuildMeResponse(actor datascope.Actor) MeResponse {
	summary := datascope.Merge(actor.Grants, actor.IsSuperAdmin)
	return MeResponse{
		UserID:       actor.UserID,
		Username:     actor.Username,
		DepartmentID: actor.DepartmentID,
		RoleCodes:    actor.RoleCodes,
		IsSuperAdmin: actor.IsSuperAdmin,
		DataScope: MeDataScopeResult{
			AllowAll:            summary.AllowAll,
			RequireSelf:         summary.RequireSelf,
			IncludeDepartment:   summary.IncludeDepartment,
			IncludeDeptTree:     summary.IncludeDeptTree,
			CustomDepartmentIDs: summary.CustomDepartmentIDs,
		},
	}
}

func NormalizeUpdateAccountProfileRequest(req UpdateAccountProfileRequest) (UpdateAccountProfileRequest, error) {
	req.Nickname = strings.TrimSpace(req.Nickname)
	if req.Nickname == "" {
		return UpdateAccountProfileRequest{}, errorsx.BadRequest("昵称不能为空")
	}
	if len(req.Nickname) > 64 {
		return UpdateAccountProfileRequest{}, errorsx.BadRequest("昵称不能超过 64 个字符")
	}
	return req, nil
}

func NormalizeUpdateAccountPasswordRequest(req UpdateAccountPasswordRequest) (UpdateAccountPasswordRequest, error) {
	if strings.TrimSpace(req.OldPassword) == "" {
		return UpdateAccountPasswordRequest{}, errorsx.BadRequest("当前密码不能为空")
	}
	if len(req.NewPassword) < 8 || len(req.NewPassword) > 72 {
		return UpdateAccountPasswordRequest{}, errorsx.BadRequest("新密码长度需要在 8 到 72 个字符之间")
	}
	if req.OldPassword == req.NewPassword {
		return UpdateAccountPasswordRequest{}, errorsx.BadRequest("新密码不能与当前密码相同")
	}
	return req, nil
}

func BuildAccountProfileResponse(actor datascope.Actor, user model.User, departmentName string) AccountProfileResponse {
	summary := datascope.Merge(actor.Grants, actor.IsSuperAdmin)
	return AccountProfileResponse{
		UserID:         user.ID,
		Username:       user.Username,
		Nickname:       user.Nickname,
		DepartmentID:   user.DepartmentID,
		DepartmentName: departmentName,
		Status:         user.Status,
		RoleCodes:      actor.RoleCodes,
		IsSuperAdmin:   actor.IsSuperAdmin,
		DataScope: MeDataScopeResult{
			AllowAll:            summary.AllowAll,
			RequireSelf:         summary.RequireSelf,
			IncludeDepartment:   summary.IncludeDepartment,
			IncludeDeptTree:     summary.IncludeDeptTree,
			CustomDepartmentIDs: summary.CustomDepartmentIDs,
		},
		UpdatedAt: user.UpdatedAt,
	}
}
