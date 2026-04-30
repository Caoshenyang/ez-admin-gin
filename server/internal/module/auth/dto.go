package auth

import (
	"strings"
	"time"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
	"ez-admin-gin/server/internal/platform/datascope"
)

// LoginRequest 表示登录请求体。
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 表示登录响应体。
type LoginResponse struct {
	UserID      uint   `json:"user_id"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresAt   string `json:"expires_at"`
}

// MeResponse 表示当前登录用户信息。
type MeResponse struct {
	UserID       uint              `json:"user_id"`
	Username     string            `json:"username"`
	DepartmentID uint              `json:"department_id"`
	RoleCodes    []string          `json:"role_codes"`
	IsSuperAdmin bool              `json:"is_super_admin"`
	DataScope    MeDataScopeResult `json:"data_scope"`
}

// MeDataScopeResult 表示当前登录人的聚合数据范围摘要。
type MeDataScopeResult struct {
	AllowAll            bool   `json:"allow_all"`
	RequireSelf         bool   `json:"require_self"`
	IncludeDepartment   bool   `json:"include_department"`
	IncludeDeptTree     bool   `json:"include_dept_tree"`
	CustomDepartmentIDs []uint `json:"custom_department_ids"`
}

// MenuResponse 表示当前登录用户可见菜单节点。
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

// DashboardCurrentUser 表示工作台当前用户摘要。
type DashboardCurrentUser struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

// DashboardHealth 表示工作台健康概览。
type DashboardHealth struct {
	Env      string `json:"env"`
	Database string `json:"database"`
	Redis    string `json:"redis"`
}

// DashboardMetrics 表示工作台指标摘要。
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

// DashboardOperationItem 表示最近操作记录。
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

// DashboardLoginItem 表示最近登录记录。
type DashboardLoginItem struct {
	ID        uint                 `json:"id"`
	Username  string               `json:"username"`
	Status    model.LoginLogStatus `json:"status"`
	Message   string               `json:"message"`
	IP        string               `json:"ip"`
	CreatedAt time.Time            `json:"created_at"`
}

// DashboardNoticeItem 表示最近公告摘要。
type DashboardNoticeItem struct {
	ID        uint               `json:"id"`
	Title     string             `json:"title"`
	Status    model.NoticeStatus `json:"status"`
	UpdatedAt time.Time          `json:"updated_at"`
}

// DashboardResponse 表示工作台响应体。
type DashboardResponse struct {
	CurrentUser      DashboardCurrentUser     `json:"current_user"`
	Health           DashboardHealth          `json:"health"`
	Metrics          DashboardMetrics         `json:"metrics"`
	RecentOperations []DashboardOperationItem `json:"recent_operations"`
	RecentLogins     []DashboardLoginItem     `json:"recent_logins"`
	LatestNotices    []DashboardNoticeItem    `json:"latest_notices"`
}

// NormalizeLoginRequest 收口登录参数。
func NormalizeLoginRequest(req LoginRequest) (LoginRequest, error) {
	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || req.Password == "" {
		return LoginRequest{}, apperror.BadRequest("用户名和密码不能为空")
	}

	return req, nil
}

// BuildMeResponse 把 Actor 压成当前登录用户返回结构。
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
