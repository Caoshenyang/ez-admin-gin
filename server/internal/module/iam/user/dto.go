package user

import (
	"strings"
	"time"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
)

// ListQuery 表示用户列表的查询参数。
type ListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Keyword  string `form:"keyword"`
	Status   int    `form:"status"`
}

// CreateRequest 表示创建用户的请求体。
type CreateRequest struct {
	Username     string           `json:"username"`
	Password     string           `json:"password"`
	Nickname     string           `json:"nickname"`
	DepartmentID uint             `json:"department_id"`
	Status       model.UserStatus `json:"status"`
	RoleIDs      []uint           `json:"role_ids"`
	PostIDs      []uint           `json:"post_ids"`
}

// UpdateRequest 表示编辑用户基础信息的请求体。
type UpdateRequest struct {
	Nickname     string           `json:"nickname"`
	DepartmentID uint             `json:"department_id"`
	Status       model.UserStatus `json:"status"`
	PostIDs      []uint           `json:"post_ids"`
}

// UpdateStatusRequest 表示单独修改用户状态的请求体。
type UpdateStatusRequest struct {
	Status model.UserStatus `json:"status"`
}

// UpdateRolesRequest 表示更新用户角色集合的请求体。
type UpdateRolesRequest struct {
	RoleIDs []uint `json:"role_ids"`
}

// Response 表示用户管理接口返回的用户对象。
type Response struct {
	ID           uint             `json:"id"`
	Username     string           `json:"username"`
	Nickname     string           `json:"nickname"`
	DepartmentID uint             `json:"department_id"`
	Status       model.UserStatus `json:"status"`
	RoleIDs      []uint           `json:"role_ids"`
	PostIDs      []uint           `json:"post_ids"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

// ListResponse 表示用户分页结果。
type ListResponse struct {
	Items    []Response `json:"items"`
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}

// NormalizePage 统一修正分页参数，避免各个 Handler 重复写边界逻辑。
func NormalizePage(page int, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return page, pageSize
}

// NormalizeCreateRequest 统一校验并收敛创建用户参数。
func NormalizeCreateRequest(req CreateRequest) (CreateRequest, error) {
	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" {
		return CreateRequest{}, apperror.BadRequest("用户名不能为空")
	}
	if len(req.Username) > 64 {
		return CreateRequest{}, apperror.BadRequest("用户名不能超过 64 个字符")
	}

	if len(req.Password) < 8 || len(req.Password) > 72 {
		return CreateRequest{}, apperror.BadRequest("密码长度需要在 8 到 72 个字符之间")
	}

	req.Nickname = strings.TrimSpace(req.Nickname)
	if req.Nickname == "" {
		req.Nickname = req.Username
	}
	if len(req.Nickname) > 64 {
		return CreateRequest{}, apperror.BadRequest("昵称不能超过 64 个字符")
	}

	if req.Status == 0 {
		req.Status = model.UserStatusEnabled
	}
	if !ValidStatus(req.Status) {
		return CreateRequest{}, apperror.BadRequest("用户状态不正确")
	}

	roleIDs, err := NormalizeRoleIDs(req.RoleIDs)
	if err != nil {
		return CreateRequest{}, err
	}
	req.RoleIDs = roleIDs
	postIDs, err := NormalizePostIDs(req.PostIDs)
	if err != nil {
		return CreateRequest{}, err
	}
	req.PostIDs = postIDs

	return req, nil
}

// NormalizeUpdateRequest 统一校验编辑用户参数。
func NormalizeUpdateRequest(req UpdateRequest) (UpdateRequest, error) {
	req.Nickname = strings.TrimSpace(req.Nickname)
	if req.Nickname == "" {
		return UpdateRequest{}, apperror.BadRequest("昵称不能为空")
	}
	if len(req.Nickname) > 64 {
		return UpdateRequest{}, apperror.BadRequest("昵称不能超过 64 个字符")
	}
	if !ValidStatus(req.Status) {
		return UpdateRequest{}, apperror.BadRequest("用户状态不正确")
	}
	postIDs, err := NormalizePostIDs(req.PostIDs)
	if err != nil {
		return UpdateRequest{}, err
	}
	req.PostIDs = postIDs

	return req, nil
}

// NormalizeRoleIDs 对角色 ID 做去重和基础校验。
func NormalizeRoleIDs(roleIDs []uint) ([]uint, error) {
	return normalizeUintIDs(roleIDs, "角色 ID 不正确")
}

// NormalizePostIDs 对岗位 ID 做去重和基础校验。
func NormalizePostIDs(postIDs []uint) ([]uint, error) {
	return normalizeUintIDs(postIDs, "岗位 ID 不正确")
}

func normalizeUintIDs(ids []uint, invalidMessage string) ([]uint, error) {
	unique := make([]uint, 0, len(ids))
	seen := make(map[uint]struct{}, len(ids))

	for _, id := range ids {
		if id == 0 {
			return nil, apperror.BadRequest(invalidMessage)
		}
		if _, ok := seen[id]; ok {
			continue
		}

		seen[id] = struct{}{}
		unique = append(unique, id)
	}

	return unique, nil
}

// ValidStatus 统一判断用户状态值是否合法。
func ValidStatus(status model.UserStatus) bool {
	return status == model.UserStatusEnabled || status == model.UserStatusDisabled
}

// BuildResponse 把模型对象和角色集合压成 API 返回结构。
func BuildResponse(user model.User, roleIDs []uint, postIDs []uint) Response {
	return Response{
		ID:           user.ID,
		Username:     user.Username,
		Nickname:     user.Nickname,
		DepartmentID: user.DepartmentID,
		Status:       user.Status,
		RoleIDs:      roleIDs,
		PostIDs:      postIDs,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}
