package system

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/middleware"
	"ez-admin-gin/server/internal/model"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserHandler 负责后台用户管理接口。
type UserHandler struct {
	db  *gorm.DB
	log *zap.Logger
}

// NewUserHandler 创建用户管理 Handler。
func NewUserHandler(db *gorm.DB, log *zap.Logger) *UserHandler {
	return &UserHandler{
		db:  db,
		log: log,
	}
}

type userListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Keyword  string `form:"keyword"`
	Status   int    `form:"status"`
}

type createUserRequest struct {
	Username string           `json:"username"`
	Password string           `json:"password"`
	Nickname string           `json:"nickname"`
	Status   model.UserStatus `json:"status"`
	RoleIDs  []uint           `json:"role_ids"`
}

type updateUserRequest struct {
	Nickname string           `json:"nickname"`
	Status   model.UserStatus `json:"status"`
}

type updateUserStatusRequest struct {
	Status model.UserStatus `json:"status"`
}

type updateUserRolesRequest struct {
	RoleIDs []uint `json:"role_ids"`
}

type userResponse struct {
	ID        uint             `json:"id"`
	Username  string           `json:"username"`
	Nickname  string           `json:"nickname"`
	Status    model.UserStatus `json:"status"`
	RoleIDs   []uint           `json:"role_ids"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

type userListResponse struct {
	Items    []userResponse `json:"items"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

// List 返回后台用户分页列表。
func (h *UserHandler) List(c *gin.Context) {
	var query userListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, apperror.BadRequest("查询参数不正确"), h.log)
		return
	}

	page, pageSize := normalizePage(query.Page, query.PageSize)
	queryDB := h.db.Model(&model.User{})

	keyword := strings.TrimSpace(query.Keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		queryDB = queryDB.Where("username LIKE ? OR nickname LIKE ?", like, like)
	}

	if query.Status != 0 {
		status := model.UserStatus(query.Status)
		if !validUserStatus(status) {
			response.Error(c, apperror.BadRequest("用户状态不正确"), h.log)
			return
		}
		queryDB = queryDB.Where("status = ?", status)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		response.Error(c, apperror.Internal("查询用户总数失败", err), h.log)
		return
	}

	var users []model.User
	if err := queryDB.
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&users).Error; err != nil {
		response.Error(c, apperror.Internal("查询用户列表失败", err), h.log)
		return
	}

	roleIDs, err := h.userRoleIDs(users)
	if err != nil {
		response.Error(c, apperror.Internal("查询用户角色失败", err), h.log)
		return
	}

	items := make([]userResponse, 0, len(users))
	for _, user := range users {
		items = append(items, buildUserResponse(user, roleIDs[user.ID]))
	}

	response.Success(c, userListResponse{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// Create 创建后台用户。
func (h *UserHandler) Create(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	username, password, nickname, status, roleIDs, err := normalizeCreateUserRequest(req)
	if err != nil {
		response.Error(c, err, h.log)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Error(c, apperror.Internal("生成密码哈希失败", err), h.log)
		return
	}

	var created model.User
	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := ensureUsernameAvailable(tx, username); err != nil {
			return err
		}

		if err := ensureRolesUsable(tx, roleIDs); err != nil {
			return err
		}

		user := model.User{
			Username:     username,
			PasswordHash: string(passwordHash),
			Nickname:     nickname,
			Status:       status,
		}

		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		if err := replaceUserRoles(tx, user.ID, roleIDs); err != nil {
			return err
		}

		created = user
		return nil
	})
	if err != nil {
		writeError(c, err, "创建用户失败", h.log)
		return
	}

	response.Success(c, buildUserResponse(created, roleIDs))
}

// Update 编辑用户基础信息。
func (h *UserHandler) Update(c *gin.Context) {
	userID, ok := userIDParam(c, h.log)
	if !ok {
		return
	}

	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	nickname, status, err := normalizeUpdateUserRequest(req)
	if err != nil {
		response.Error(c, err, h.log)
		return
	}

	if currentUserID, ok := middleware.CurrentUserID(c); ok && currentUserID == userID && status == model.UserStatusDisabled {
		response.Error(c, apperror.BadRequest("不能禁用当前登录用户"), h.log)
		return
	}

	var user model.User
	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&user, userID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("用户不存在")
			}
			return err
		}

		if err := tx.Model(&user).Updates(map[string]any{
			"nickname": nickname,
			"status":   status,
		}).Error; err != nil {
			return err
		}

		user.Nickname = nickname
		user.Status = status
		return nil
	})
	if err != nil {
		writeError(c, err, "更新用户失败", h.log)
		return
	}

	roleIDs, err := h.userRoleIDs([]model.User{user})
	if err != nil {
		response.Error(c, apperror.Internal("查询用户角色失败", err), h.log)
		return
	}

	response.Success(c, buildUserResponse(user, roleIDs[user.ID]))
}

// UpdateStatus 修改用户启用状态。
func (h *UserHandler) UpdateStatus(c *gin.Context) {
	userID, ok := userIDParam(c, h.log)
	if !ok {
		return
	}

	var req updateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	if !validUserStatus(req.Status) {
		response.Error(c, apperror.BadRequest("用户状态不正确"), h.log)
		return
	}

	if currentUserID, ok := middleware.CurrentUserID(c); ok && currentUserID == userID && req.Status == model.UserStatusDisabled {
		response.Error(c, apperror.BadRequest("不能禁用当前登录用户"), h.log)
		return
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		var user model.User
		if err := tx.First(&user, userID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("用户不存在")
			}
			return err
		}

		return tx.Model(&user).Update("status", req.Status).Error
	})
	if err != nil {
		writeError(c, err, "更新用户状态失败", h.log)
		return
	}

	response.Success(c, gin.H{
		"id":     userID,
		"status": req.Status,
	})
}

// UpdateRoles 更新用户绑定的角色。
func (h *UserHandler) UpdateRoles(c *gin.Context) {
	userID, ok := userIDParam(c, h.log)
	if !ok {
		return
	}

	if currentUserID, ok := middleware.CurrentUserID(c); ok && currentUserID == userID {
		response.Error(c, apperror.BadRequest("不能修改当前登录用户的角色"), h.log)
		return
	}

	var req updateUserRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	roleIDs, err := normalizeRoleIDs(req.RoleIDs)
	if err != nil {
		response.Error(c, err, h.log)
		return
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		var user model.User
		if err := tx.First(&user, userID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("用户不存在")
			}
			return err
		}

		if err := ensureRolesUsable(tx, roleIDs); err != nil {
			return err
		}

		return replaceUserRoles(tx, userID, roleIDs)
	})
	if err != nil {
		writeError(c, err, "更新用户角色失败", h.log)
		return
	}

	response.Success(c, gin.H{
		"id":       userID,
		"role_ids": roleIDs,
	})
}

func (h *UserHandler) userRoleIDs(users []model.User) (map[uint][]uint, error) {
	result := make(map[uint][]uint, len(users))
	if len(users) == 0 {
		return result, nil
	}

	userIDs := make([]uint, 0, len(users))
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}

	var rows []model.UserRole
	if err := h.db.Where("user_id IN ?", userIDs).Order("role_id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		result[row.UserID] = append(result[row.UserID], row.RoleID)
	}

	return result, nil
}

func normalizeCreateUserRequest(req createUserRequest) (string, string, string, model.UserStatus, []uint, error) {
	username := strings.TrimSpace(req.Username)
	if username == "" {
		return "", "", "", 0, nil, apperror.BadRequest("用户名不能为空")
	}
	if len(username) > 64 {
		return "", "", "", 0, nil, apperror.BadRequest("用户名不能超过 64 个字符")
	}

	if len(req.Password) < 8 || len(req.Password) > 72 {
		return "", "", "", 0, nil, apperror.BadRequest("密码长度需要在 8 到 72 个字符之间")
	}

	nickname := strings.TrimSpace(req.Nickname)
	if nickname == "" {
		nickname = username
	}
	if len(nickname) > 64 {
		return "", "", "", 0, nil, apperror.BadRequest("昵称不能超过 64 个字符")
	}

	status := req.Status
	if status == 0 {
		status = model.UserStatusEnabled
	}
	if !validUserStatus(status) {
		return "", "", "", 0, nil, apperror.BadRequest("用户状态不正确")
	}

	roleIDs, err := normalizeRoleIDs(req.RoleIDs)
	if err != nil {
		return "", "", "", 0, nil, err
	}

	return username, req.Password, nickname, status, roleIDs, nil
}

func normalizeUpdateUserRequest(req updateUserRequest) (string, model.UserStatus, error) {
	nickname := strings.TrimSpace(req.Nickname)
	if nickname == "" {
		return "", 0, apperror.BadRequest("昵称不能为空")
	}
	if len(nickname) > 64 {
		return "", 0, apperror.BadRequest("昵称不能超过 64 个字符")
	}

	if !validUserStatus(req.Status) {
		return "", 0, apperror.BadRequest("用户状态不正确")
	}

	return nickname, req.Status, nil
}

func normalizeRoleIDs(roleIDs []uint) ([]uint, error) {
	unique := make([]uint, 0, len(roleIDs))
	seen := make(map[uint]struct{}, len(roleIDs))

	for _, roleID := range roleIDs {
		if roleID == 0 {
			return nil, apperror.BadRequest("角色 ID 不正确")
		}
		if _, ok := seen[roleID]; ok {
			continue
		}

		seen[roleID] = struct{}{}
		unique = append(unique, roleID)
	}

	return unique, nil
}

func normalizePage(page int, pageSize int) (int, int) {
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

func validUserStatus(status model.UserStatus) bool {
	return status == model.UserStatusEnabled || status == model.UserStatusDisabled
}

func userIDParam(c *gin.Context, log *zap.Logger) (uint, bool) {
	rawID := c.Param("id")
	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil || id == 0 {
		response.Error(c, apperror.BadRequest("用户 ID 不正确"), log)
		return 0, false
	}

	return uint(id), true
}

func ensureUsernameAvailable(db *gorm.DB, username string) error {
	var user model.User
	err := db.Unscoped().Where("username = ?", username).First(&user).Error
	if err == nil {
		return apperror.BadRequest("用户名已存在")
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return err
}

func ensureRolesUsable(db *gorm.DB, roleIDs []uint) error {
	if len(roleIDs) == 0 {
		return nil
	}

	var count int64
	err := db.Model(&model.Role{}).
		Where("id IN ?", roleIDs).
		Where("status = ?", model.RoleStatusEnabled).
		Count(&count).Error
	if err != nil {
		return err
	}

	if count != int64(len(roleIDs)) {
		return apperror.BadRequest("角色不存在或已禁用")
	}

	return nil
}

func replaceUserRoles(db *gorm.DB, userID uint, roleIDs []uint) error {
	if err := db.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
		return err
	}

	if len(roleIDs) == 0 {
		return nil
	}

	rows := make([]model.UserRole, 0, len(roleIDs))
	for _, roleID := range roleIDs {
		rows = append(rows, model.UserRole{
			UserID: userID,
			RoleID: roleID,
		})
	}

	return db.Create(&rows).Error
}

func buildUserResponse(user model.User, roleIDs []uint) userResponse {
	return userResponse{
		ID:        user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Status:    user.Status,
		RoleIDs:   roleIDs,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func writeError(c *gin.Context, err error, fallbackMessage string, log *zap.Logger) {
	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		response.Error(c, appErr, log)
		return
	}

	response.Error(c, apperror.Internal(fallbackMessage, err), log)
}
