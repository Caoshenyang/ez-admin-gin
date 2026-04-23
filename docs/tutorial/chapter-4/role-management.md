---
title: 角色管理
description: "实现角色维护、接口权限分配和菜单权限分配能力。"
---

# 角色管理

上一节已经能管理用户，并能给用户绑定角色。这一节继续补齐角色本身的管理能力：创建角色、编辑角色、禁用角色，并为角色分配接口权限和菜单权限。

::: tip 🎯 本节目标
完成后，`super_admin` 可以访问角色管理接口；系统会初始化角色管理菜单和按钮；通过接口可以维护角色，并把接口权限、菜单权限分配给指定角色。
:::

## 本节会改什么

本节会新增或修改下面这些文件：

```text
server/
├─ internal/
│  ├─ bootstrap/
│  │  └─ bootstrap.go
│  ├─ handler/
│  │  └─ system/
│  │     └─ roles.go
│  └─ router/
│     └─ router.go
```

| 位置 | 用途 |
| --- | --- |
| `internal/handler/system/roles.go` | 角色管理接口 |
| `internal/router/router.go` | 注册角色管理路由 |
| `internal/bootstrap/bootstrap.go` | 初始化角色管理权限和菜单 |

::: info 本节不新增数据库表
角色管理复用 `sys_role`、`sys_role_menu`、`casbin_rule`。其中 `sys_role` 保存角色本身，`sys_role_menu` 保存角色菜单关系，`casbin_rule` 保存接口权限策略。
:::

## 接口规划

本节先实现 6 个接口：

| 方法 | 路径 | 用途 |
| --- | --- | --- |
| `GET` | `/api/v1/system/roles` | 角色列表 |
| `POST` | `/api/v1/system/roles` | 创建角色 |
| `POST` | `/api/v1/system/roles/:id/update` | 编辑角色基础信息 |
| `POST` | `/api/v1/system/roles/:id/status` | 修改角色状态 |
| `POST` | `/api/v1/system/roles/:id/permissions` | 分配接口权限 |
| `POST` | `/api/v1/system/roles/:id/menus` | 分配菜单权限 |

::: warning ⚠️ 不要随意禁用 `super_admin`
`super_admin` 是本教程本地起步的超级管理员角色。为了避免把自己锁在系统外，本节会阻止禁用这个角色。
:::

## 🛠️ 创建角色管理 Handler

创建 `server/internal/handler/system/roles.go`。这是新增文件，直接完整写入即可。

```go
package system

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const superAdminRoleCode = "super_admin"

// RoleHandler 负责后台角色管理接口。
type RoleHandler struct {
	db  *gorm.DB
	log *zap.Logger
}

// NewRoleHandler 创建角色管理 Handler。
func NewRoleHandler(db *gorm.DB, log *zap.Logger) *RoleHandler {
	return &RoleHandler{
		db:  db,
		log: log,
	}
}

type roleListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Keyword  string `form:"keyword"`
	Status   int    `form:"status"`
}

type createRoleRequest struct {
	Code   string           `json:"code"`
	Name   string           `json:"name"`
	Sort   int              `json:"sort"`
	Status model.RoleStatus `json:"status"`
	Remark string           `json:"remark"`
}

type updateRoleRequest struct {
	Name   string           `json:"name"`
	Sort   int              `json:"sort"`
	Status model.RoleStatus `json:"status"`
	Remark string           `json:"remark"`
}

type updateRoleStatusRequest struct {
	Status model.RoleStatus `json:"status"`
}

type rolePermissionItem struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

type updateRolePermissionsRequest struct {
	Permissions []rolePermissionItem `json:"permissions"`
}

type updateRoleMenusRequest struct {
	MenuIDs []uint `json:"menu_ids"`
}

type roleResponse struct {
	ID          uint                 `json:"id"`
	Code        string               `json:"code"`
	Name        string               `json:"name"`
	Sort        int                  `json:"sort"`
	Status      model.RoleStatus     `json:"status"`
	Remark      string               `json:"remark"`
	Permissions []rolePermissionItem `json:"permissions"`
	MenuIDs     []uint               `json:"menu_ids"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

type roleListResponse struct {
	Items    []roleResponse `json:"items"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

// List 返回角色分页列表。
func (h *RoleHandler) List(c *gin.Context) {
	var query roleListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, apperror.BadRequest("查询参数不正确"), h.log)
		return
	}

	page, pageSize := normalizeRolePage(query.Page, query.PageSize)
	queryDB := h.db.Model(&model.Role{})

	keyword := strings.TrimSpace(query.Keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		queryDB = queryDB.Where("code LIKE ? OR name LIKE ?", like, like)
	}

	if query.Status != 0 {
		status := model.RoleStatus(query.Status)
		if !validRoleStatus(status) {
			response.Error(c, apperror.BadRequest("角色状态不正确"), h.log)
			return
		}
		queryDB = queryDB.Where("status = ?", status)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		response.Error(c, apperror.Internal("查询角色总数失败", err), h.log)
		return
	}

	var roles []model.Role
	if err := queryDB.
		Order("sort ASC, id ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&roles).Error; err != nil {
		response.Error(c, apperror.Internal("查询角色列表失败", err), h.log)
		return
	}

	permissions, err := h.rolePermissions(roles)
	if err != nil {
		response.Error(c, apperror.Internal("查询角色接口权限失败", err), h.log)
		return
	}

	menuIDs, err := h.roleMenuIDs(roles)
	if err != nil {
		response.Error(c, apperror.Internal("查询角色菜单权限失败", err), h.log)
		return
	}

	items := make([]roleResponse, 0, len(roles))
	for _, role := range roles {
		items = append(items, buildRoleResponse(role, permissions[role.Code], menuIDs[role.ID]))
	}

	response.Success(c, roleListResponse{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// Create 创建角色。
func (h *RoleHandler) Create(c *gin.Context) {
	var req createRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	code, name, status, remark, err := normalizeCreateRoleRequest(req)
	if err != nil {
		response.Error(c, err, h.log)
		return
	}

	var created model.Role
	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := ensureRoleCodeAvailable(tx, code); err != nil {
			return err
		}

		role := model.Role{
			Code:   code,
			Name:   name,
			Sort:   req.Sort,
			Status: status,
			Remark: remark,
		}

		if err := tx.Create(&role).Error; err != nil {
			return err
		}

		created = role
		return nil
	})
	if err != nil {
		writeRoleError(c, err, "创建角色失败", h.log)
		return
	}

	response.Success(c, buildRoleResponse(created, nil, nil))
}

// Update 编辑角色基础信息。
func (h *RoleHandler) Update(c *gin.Context) {
	roleID, ok := roleIDParam(c, h.log)
	if !ok {
		return
	}

	var req updateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	name, status, remark, err := normalizeUpdateRoleRequest(req)
	if err != nil {
		response.Error(c, err, h.log)
		return
	}

	var role model.Role
	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&role, roleID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("角色不存在")
			}
			return err
		}

		if role.Code == superAdminRoleCode && status == model.RoleStatusDisabled {
			return apperror.BadRequest("不能禁用超级管理员角色")
		}

		if err := tx.Model(&role).Updates(map[string]any{
			"name":   name,
			"sort":   req.Sort,
			"status": status,
			"remark": remark,
		}).Error; err != nil {
			return err
		}

		role.Name = name
		role.Sort = req.Sort
		role.Status = status
		role.Remark = remark
		return nil
	})
	if err != nil {
		writeRoleError(c, err, "更新角色失败", h.log)
		return
	}

	response.Success(c, buildRoleResponse(role, nil, nil))
}

// UpdateStatus 修改角色状态。
func (h *RoleHandler) UpdateStatus(c *gin.Context) {
	roleID, ok := roleIDParam(c, h.log)
	if !ok {
		return
	}

	var req updateRoleStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	if !validRoleStatus(req.Status) {
		response.Error(c, apperror.BadRequest("角色状态不正确"), h.log)
		return
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		var role model.Role
		if err := tx.First(&role, roleID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("角色不存在")
			}
			return err
		}

		if role.Code == superAdminRoleCode && req.Status == model.RoleStatusDisabled {
			return apperror.BadRequest("不能禁用超级管理员角色")
		}

		return tx.Model(&role).Update("status", req.Status).Error
	})
	if err != nil {
		writeRoleError(c, err, "更新角色状态失败", h.log)
		return
	}

	response.Success(c, gin.H{
		"id":     roleID,
		"status": req.Status,
	})
}

// UpdatePermissions 替换角色接口权限。
func (h *RoleHandler) UpdatePermissions(c *gin.Context) {
	roleID, ok := roleIDParam(c, h.log)
	if !ok {
		return
	}

	var req updateRolePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	permissions, err := normalizeRolePermissions(req.Permissions)
	if err != nil {
		response.Error(c, err, h.log)
		return
	}

	var role model.Role
	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&role, roleID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("角色不存在")
			}
			return err
		}

		if role.Code == superAdminRoleCode {
			return apperror.BadRequest("超级管理员角色权限不在这里修改")
		}

		return replaceRolePermissions(tx, role.Code, permissions)
	})
	if err != nil {
		writeRoleError(c, err, "更新角色接口权限失败", h.log)
		return
	}

	response.Success(c, gin.H{
		"id":          roleID,
		"code":        role.Code,
		"permissions": permissions,
	})
}

// UpdateMenus 替换角色菜单权限。
func (h *RoleHandler) UpdateMenus(c *gin.Context) {
	roleID, ok := roleIDParam(c, h.log)
	if !ok {
		return
	}

	var req updateRoleMenusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	menuIDs, err := normalizeMenuIDs(req.MenuIDs)
	if err != nil {
		response.Error(c, err, h.log)
		return
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		var role model.Role
		if err := tx.First(&role, roleID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("角色不存在")
			}
			return err
		}

		if role.Code == superAdminRoleCode {
			return apperror.BadRequest("超级管理员菜单权限不在这里修改")
		}

		if err := ensureMenusUsable(tx, menuIDs); err != nil {
			return err
		}

		return replaceRoleMenus(tx, roleID, menuIDs)
	})
	if err != nil {
		writeRoleError(c, err, "更新角色菜单权限失败", h.log)
		return
	}

	response.Success(c, gin.H{
		"id":       roleID,
		"menu_ids": menuIDs,
	})
}

func (h *RoleHandler) rolePermissions(roles []model.Role) (map[string][]rolePermissionItem, error) {
	result := make(map[string][]rolePermissionItem, len(roles))
	if len(roles) == 0 {
		return result, nil
	}

	roleCodes := make([]string, 0, len(roles))
	for _, role := range roles {
		roleCodes = append(roleCodes, role.Code)
	}

	var rows []model.CasbinRule
	if err := h.db.
		Where("ptype = ?", "p").
		Where("v0 IN ?", roleCodes).
		Order("v1 ASC, v2 ASC").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		result[row.V0] = append(result[row.V0], rolePermissionItem{
			Path:   row.V1,
			Method: row.V2,
		})
	}

	return result, nil
}

func (h *RoleHandler) roleMenuIDs(roles []model.Role) (map[uint][]uint, error) {
	result := make(map[uint][]uint, len(roles))
	if len(roles) == 0 {
		return result, nil
	}

	roleIDs := make([]uint, 0, len(roles))
	for _, role := range roles {
		roleIDs = append(roleIDs, role.ID)
	}

	var rows []model.RoleMenu
	if err := h.db.Where("role_id IN ?", roleIDs).Order("menu_id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		result[row.RoleID] = append(result[row.RoleID], row.MenuID)
	}

	return result, nil
}

func normalizeCreateRoleRequest(req createRoleRequest) (string, string, model.RoleStatus, string, error) {
	code := strings.TrimSpace(req.Code)
	if code == "" {
		return "", "", 0, "", apperror.BadRequest("角色编码不能为空")
	}
	if len(code) > 64 {
		return "", "", 0, "", apperror.BadRequest("角色编码不能超过 64 个字符")
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return "", "", 0, "", apperror.BadRequest("角色名称不能为空")
	}
	if len(name) > 64 {
		return "", "", 0, "", apperror.BadRequest("角色名称不能超过 64 个字符")
	}

	status := req.Status
	if status == 0 {
		status = model.RoleStatusEnabled
	}
	if !validRoleStatus(status) {
		return "", "", 0, "", apperror.BadRequest("角色状态不正确")
	}

	remark := strings.TrimSpace(req.Remark)
	if len(remark) > 255 {
		return "", "", 0, "", apperror.BadRequest("备注不能超过 255 个字符")
	}

	return code, name, status, remark, nil
}

func normalizeUpdateRoleRequest(req updateRoleRequest) (string, model.RoleStatus, string, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return "", 0, "", apperror.BadRequest("角色名称不能为空")
	}
	if len(name) > 64 {
		return "", 0, "", apperror.BadRequest("角色名称不能超过 64 个字符")
	}

	if !validRoleStatus(req.Status) {
		return "", 0, "", apperror.BadRequest("角色状态不正确")
	}

	remark := strings.TrimSpace(req.Remark)
	if len(remark) > 255 {
		return "", 0, "", apperror.BadRequest("备注不能超过 255 个字符")
	}

	return name, req.Status, remark, nil
}

func normalizeRolePermissions(permissions []rolePermissionItem) ([]rolePermissionItem, error) {
	unique := make([]rolePermissionItem, 0, len(permissions))
	seen := make(map[string]struct{}, len(permissions))

	for _, item := range permissions {
		path := strings.TrimSpace(item.Path)
		method := strings.ToUpper(strings.TrimSpace(item.Method))
		if path == "" || method == "" {
			return nil, apperror.BadRequest("接口权限参数不正确")
		}

		key := path + " " + method
		if _, ok := seen[key]; ok {
			continue
		}

		seen[key] = struct{}{}
		unique = append(unique, rolePermissionItem{
			Path:   path,
			Method: method,
		})
	}

	return unique, nil
}

func normalizeMenuIDs(menuIDs []uint) ([]uint, error) {
	unique := make([]uint, 0, len(menuIDs))
	seen := make(map[uint]struct{}, len(menuIDs))

	for _, menuID := range menuIDs {
		if menuID == 0 {
			return nil, apperror.BadRequest("菜单 ID 不正确")
		}
		if _, ok := seen[menuID]; ok {
			continue
		}

		seen[menuID] = struct{}{}
		unique = append(unique, menuID)
	}

	return unique, nil
}

func normalizeRolePage(page int, pageSize int) (int, int) {
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

func validRoleStatus(status model.RoleStatus) bool {
	return status == model.RoleStatusEnabled || status == model.RoleStatusDisabled
}

func roleIDParam(c *gin.Context, log *zap.Logger) (uint, bool) {
	rawID := c.Param("id")
	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil || id == 0 {
		response.Error(c, apperror.BadRequest("角色 ID 不正确"), log)
		return 0, false
	}

	return uint(id), true
}

func ensureRoleCodeAvailable(db *gorm.DB, code string) error {
	var role model.Role
	err := db.Unscoped().Where("code = ?", code).First(&role).Error
	if err == nil {
		return apperror.BadRequest("角色编码已存在")
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return err
}

func ensureMenusUsable(db *gorm.DB, menuIDs []uint) error {
	if len(menuIDs) == 0 {
		return nil
	}

	var count int64
	err := db.Model(&model.Menu{}).
		Where("id IN ?", menuIDs).
		Where("status = ?", model.MenuStatusEnabled).
		Count(&count).Error
	if err != nil {
		return err
	}

	if count != int64(len(menuIDs)) {
		return apperror.BadRequest("菜单不存在或已禁用")
	}

	return nil
}

func replaceRolePermissions(db *gorm.DB, roleCode string, permissions []rolePermissionItem) error {
	if err := db.Where("ptype = ? AND v0 = ?", "p", roleCode).Delete(&model.CasbinRule{}).Error; err != nil {
		return err
	}

	if len(permissions) == 0 {
		return nil
	}

	rows := make([]model.CasbinRule, 0, len(permissions))
	for _, permission := range permissions {
		rows = append(rows, model.CasbinRule{
			Ptype: "p",
			V0:    roleCode,
			V1:    permission.Path,
			V2:    permission.Method,
		})
	}

	return db.Create(&rows).Error
}

func replaceRoleMenus(db *gorm.DB, roleID uint, menuIDs []uint) error {
	if err := db.Where("role_id = ?", roleID).Delete(&model.RoleMenu{}).Error; err != nil {
		return err
	}

	if len(menuIDs) == 0 {
		return nil
	}

	rows := make([]model.RoleMenu, 0, len(menuIDs))
	for _, menuID := range menuIDs {
		rows = append(rows, model.RoleMenu{
			RoleID: roleID,
			MenuID: menuID,
		})
	}

	return db.Create(&rows).Error
}

func buildRoleResponse(role model.Role, permissions []rolePermissionItem, menuIDs []uint) roleResponse {
	return roleResponse{
		ID:          role.ID,
		Code:        role.Code,
		Name:        role.Name,
		Sort:        role.Sort,
		Status:      role.Status,
		Remark:      role.Remark,
		Permissions: permissions,
		MenuIDs:     menuIDs,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}

func writeRoleError(c *gin.Context, err error, fallbackMessage string, log *zap.Logger) {
	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		response.Error(c, appErr, log)
		return
	}

	response.Error(c, apperror.Internal(fallbackMessage, err), log)
}
```

::: details 为什么不允许修改角色编码
`code` 会被写入 `casbin_rule.v0`，也是权限判断时使用的稳定标识。角色创建后可以改名称、排序、状态和备注，但不建议直接改编码。
:::

::: warning ⚠️ 修改接口权限后需要重新加载策略
本节先用“重启服务”让 Casbin 重新加载策略。后续如果要在页面里即时修改权限，需要在保存权限后调用 Enforcer 的重新加载逻辑。
:::

## 🛠️ 注册角色管理路由

修改 `server/internal/router/router.go`。这一处在系统路由中新增角色 Handler 和路由。

```go
// registerSystemRoutes 注册系统级路由。
func registerSystemRoutes(r *gin.Engine, opts Options) {
	health := systemHandler.NewHealthHandler(opts.Config, opts.DB, opts.Redis, opts.Log)
	users := systemHandler.NewUserHandler(opts.DB, opts.Log)
	roles := systemHandler.NewRoleHandler(opts.DB, opts.Log) // [!code ++]

	// /health 通常给部署探针和本地快速验证使用。
	r.GET("/health", health.Check)

	// /api/v1/system/health 放在接口版本分组下，方便统一管理后台接口。
	api := r.Group("/api/v1")
	system := api.Group("/system")
	system.Use(middleware.Auth(opts.Token, opts.Log))
	system.Use(middleware.Permission(opts.DB, opts.Permission, opts.Log))
	system.GET("/health", health.Check)
	system.GET("/users", users.List)
	system.POST("/users", users.Create)
	system.POST("/users/:id/update", users.Update)
	system.POST("/users/:id/status", users.UpdateStatus)
	system.POST("/users/:id/roles", users.UpdateRoles)
	system.GET("/roles", roles.List) // [!code ++]
	system.POST("/roles", roles.Create) // [!code ++]
	system.POST("/roles/:id/update", roles.Update) // [!code ++]
	system.POST("/roles/:id/status", roles.UpdateStatus) // [!code ++]
	system.POST("/roles/:id/permissions", roles.UpdatePermissions) // [!code ++]
	system.POST("/roles/:id/menus", roles.UpdateMenus) // [!code ++]
}
```

## 🛠️ 初始化角色管理接口权限

修改 `server/internal/bootstrap/bootstrap.go`。在上一节的 `defaultPermissionSeeds` 中继续追加角色管理接口权限：

```go
var defaultPermissionSeeds = []defaultPermissionSeed{
	{Path: "/api/v1/system/health", Method: "GET"},
	{Path: "/api/v1/system/users", Method: "GET"},
	{Path: "/api/v1/system/users", Method: "POST"},
	{Path: "/api/v1/system/users/:id/update", Method: "POST"},
	{Path: "/api/v1/system/users/:id/status", Method: "POST"},
	{Path: "/api/v1/system/users/:id/roles", Method: "POST"},
	{Path: "/api/v1/system/roles", Method: "GET"}, // [!code ++]
	{Path: "/api/v1/system/roles", Method: "POST"}, // [!code ++]
	{Path: "/api/v1/system/roles/:id/update", Method: "POST"}, // [!code ++]
	{Path: "/api/v1/system/roles/:id/status", Method: "POST"}, // [!code ++]
	{Path: "/api/v1/system/roles/:id/permissions", Method: "POST"}, // [!code ++]
	{Path: "/api/v1/system/roles/:id/menus", Method: "POST"}, // [!code ++]
}
```

## 🛠️ 初始化角色管理菜单

继续修改 `server/internal/bootstrap/bootstrap.go`。先增加角色管理菜单和按钮编码：

```go
const (
	defaultUserAssignRoleCode = "system:user:assign-role"
	defaultRoleMenuCode       = "system:role" // [!code ++]
	defaultRoleListCode       = "system:role:list" // [!code ++]
	defaultRoleCreateCode     = "system:role:create" // [!code ++]
	defaultRoleUpdateCode     = "system:role:update" // [!code ++]
	defaultRoleStatusCode     = "system:role:status" // [!code ++]
	defaultRolePermissionCode = "system:role:permission" // [!code ++]
	defaultRoleMenuAssignCode = "system:role:menu" // [!code ++]
)
```

接着修改 `seedDefaultMenus`。先找到上一节最后新增的返回语句：

```go
return menus, nil
```

把这行返回语句替换为下面整段代码。也就是说：下面代码放在用户管理按钮循环之后，原 `return menus, nil` 之前；替换完成后，函数末尾仍然只保留一个 `return menus, nil`。

```go
	roleMenu, err := seedMenu(db, model.Menu{
		ParentID:  systemMenu.ID,
		Type:      model.MenuTypeMenu,
		Code:      defaultRoleMenuCode,
		Title:     "角色管理",
		Path:      "/system/roles",
		Component: "system/RoleView",
		Icon:      "team",
		Sort:      30,
		Status:    model.MenuStatusEnabled,
		Remark:    "系统内置菜单",
	}, log)
	if err != nil {
		return nil, err
	}

	roleButtons := []model.Menu{
		{ParentID: roleMenu.ID, Type: model.MenuTypeButton, Code: defaultRoleListCode, Title: "查看角色", Sort: 10, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
		{ParentID: roleMenu.ID, Type: model.MenuTypeButton, Code: defaultRoleCreateCode, Title: "创建角色", Sort: 20, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
		{ParentID: roleMenu.ID, Type: model.MenuTypeButton, Code: defaultRoleUpdateCode, Title: "编辑角色", Sort: 30, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
		{ParentID: roleMenu.ID, Type: model.MenuTypeButton, Code: defaultRoleStatusCode, Title: "修改角色状态", Sort: 40, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
		{ParentID: roleMenu.ID, Type: model.MenuTypeButton, Code: defaultRolePermissionCode, Title: "分配接口权限", Sort: 50, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
		{ParentID: roleMenu.ID, Type: model.MenuTypeButton, Code: defaultRoleMenuAssignCode, Title: "分配菜单权限", Sort: 60, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
	}

	menus = append(menus, *roleMenu)
	for _, button := range roleButtons {
		createdButton, err := seedMenu(db, button, log)
		if err != nil {
			return nil, err
		}
		menus = append(menus, *createdButton)
	}

	return menus, nil
```

::: warning ⚠️ 这一段是替换末尾返回，不是追加到函数外面
角色管理菜单要继续复用上一节创建的 `menus` 切片，所以代码必须放在 `seedDefaultMenus` 函数内部、用户管理按钮循环之后。
:::

::: details 菜单分配为什么只保存菜单 ID
菜单的标题、路径和层级保存在 `sys_menu`。角色菜单关系表只需要保存 `role_id` 和 `menu_id`，避免重复存储菜单信息。
:::

## ✅ 整理依赖并启动

本节没有新增第三方依赖，但修改了后端文件，仍然可以整理一次：

```bash
# 在 server/ 目录执行
go mod tidy
```

确认数据库和 Redis 正在运行：

```bash
# 在项目根目录执行，确认本地依赖服务处于运行状态
docker compose -f deploy/compose.local.yml ps
```

回到 `server/` 目录启动服务：

```bash
# 在 server/ 目录启动服务
go run .
```

第一次启动后，控制台应该能看到类似日志：

```text
INFO	default permission created	{"role_code": "super_admin", "path": "/api/v1/system/roles", "method": "GET"}
INFO	default menu created	{"menu_code": "system:role"}
INFO	default role menu bound	{"role_id": 1, "menu_id": 10}
```

## ✅ 验证权限和菜单数据

先确认角色管理接口权限已经写入：

```bash
# 查看角色管理相关接口权限
docker compose -f deploy/compose.local.yml exec postgres psql -U ez_admin -d ez_admin -c "select ptype, v0, v1, v2 from casbin_rule where v1 like '/api/v1/system/roles%' order by v1, v2;"
```

应该能看到角色列表、创建、编辑、状态修改、接口权限分配、菜单权限分配对应的策略。

再确认角色管理菜单已经写入：

```bash
# 查看角色管理相关菜单和按钮
docker compose -f deploy/compose.local.yml exec postgres psql -U ez_admin -d ez_admin -c "select id, parent_id, type, code, title from sys_menu where code like 'system:role%' order by sort, id;"
```

应该能看到 `system:role` 以及几个 `system:role:*` 按钮编码。

## ✅ 验证角色管理接口

先登录拿到 Token：

::: code-group

```powershell [Windows PowerShell]
$body = @{
  username = "admin"
  password = "EzAdmin@123456"
} | ConvertTo-Json

$login = Invoke-RestMethod `
  -Method Post `
  -Uri http://localhost:8080/api/v1/auth/login `
  -ContentType "application/json" `
  -Body $body

$token = $login.data.access_token
```

```bash [macOS / Linux]
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"EzAdmin@123456"}' | jq -r '.data.access_token')
```

:::

查看角色列表：

::: code-group

```powershell [Windows PowerShell]
Invoke-RestMethod `
  -Method Get `
  -Uri "http://localhost:8080/api/v1/system/roles?page=1&page_size=10" `
  -Headers @{ Authorization = "Bearer $token" }
```

```bash [macOS / Linux]
curl "http://localhost:8080/api/v1/system/roles?page=1&page_size=10" \
  -H "Authorization: Bearer ${TOKEN}"
```

:::

应该能看到包含 `super_admin` 的分页结果。

创建一个测试角色：

::: warning ⚠️ Windows PowerShell 发送中文 JSON 时要显式使用 UTF-8
如果请求体中包含中文，建议先把 JSON 转成 UTF-8 字节再发送，避免中文在发送阶段变成 `????`。
:::

::: code-group

```powershell [Windows PowerShell]
$body = @{
  code = "demo_role"
  name = "演示角色"
  sort = 100
  status = 1
  remark = "用于验证角色管理接口"
} | ConvertTo-Json

$utf8Body = [System.Text.Encoding]::UTF8.GetBytes($body)

Invoke-RestMethod `
  -Method Post `
  -Uri http://localhost:8080/api/v1/system/roles `
  -ContentType "application/json; charset=utf-8" `
  -Headers @{ Authorization = "Bearer $token" } `
  -Body $utf8Body
```

```bash [macOS / Linux]
curl -X POST http://localhost:8080/api/v1/system/roles \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"code":"demo_role","name":"演示角色","sort":100,"status":1,"remark":"用于验证角色管理接口"}'
```

:::

给测试角色分配接口权限。下面示例假设新角色 ID 是 `2`：

::: code-group

```powershell [Windows PowerShell]
$roleId = 2
$body = @{
  permissions = @(
    @{ path = "/api/v1/system/users"; method = "GET" },
    @{ path = "/api/v1/system/roles"; method = "GET" }
  )
} | ConvertTo-Json -Depth 4

Invoke-RestMethod `
  -Method Post `
  -Uri "http://localhost:8080/api/v1/system/roles/$roleId/permissions" `
  -ContentType "application/json" `
  -Headers @{ Authorization = "Bearer $token" } `
  -Body $body
```

```bash [macOS / Linux]
ROLE_ID=2

curl -X POST "http://localhost:8080/api/v1/system/roles/${ROLE_ID}/permissions" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"permissions":[{"path":"/api/v1/system/users","method":"GET"},{"path":"/api/v1/system/roles","method":"GET"}]}'
```

:::

给测试角色分配菜单权限。下面示例里的菜单 ID 需要按本地 `sys_menu` 查询结果替换：

::: code-group

```powershell [Windows PowerShell]
$roleId = 2
$body = @{
  menu_ids = @(1, 4, 5)
} | ConvertTo-Json

Invoke-RestMethod `
  -Method Post `
  -Uri "http://localhost:8080/api/v1/system/roles/$roleId/menus" `
  -ContentType "application/json" `
  -Headers @{ Authorization = "Bearer $token" } `
  -Body $body
```

```bash [macOS / Linux]
ROLE_ID=2

curl -X POST "http://localhost:8080/api/v1/system/roles/${ROLE_ID}/menus" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"menu_ids":[1,4,5]}'
```

:::

最后用 SQL 确认数据已经写入：

```bash
# 查看 demo_role 的接口权限和菜单权限
docker compose -f deploy/compose.local.yml exec postgres psql -U ez_admin -d ez_admin -c "select ptype, v0, v1, v2 from casbin_rule where v0 = 'demo_role' order by v1, v2;"
docker compose -f deploy/compose.local.yml exec postgres psql -U ez_admin -d ez_admin -c "select rm.role_id, r.code, rm.menu_id, m.code as menu_code from sys_role_menu rm join sys_role r on r.id = rm.role_id join sys_menu m on m.id = rm.menu_id where r.code = 'demo_role' order by rm.menu_id;"
```

::: warning ⚠️ 修改接口权限后记得重启服务再验证访问效果
直接写入 `casbin_rule` 后，当前进程里的 Enforcer 不会自动刷新。现在先重启服务，让新权限生效。
:::

## 常见问题

::: details 创建角色时提示“角色编码已存在”
换一个新的角色编码即可。角色编码唯一规则见：[数据库建表语句 - `sys_role`](../../reference/database-ddl#sys-role)。
:::

::: details 分配菜单时提示“菜单不存在或已禁用”
请求里的 `menu_ids` 必须对应已经存在且启用的菜单或按钮。可以先执行下面的 SQL 查看菜单：

```sql
select id, parent_id, type, code, title, status from sys_menu order by sort, id;
```
:::

::: details 为什么不能在这个接口里修改 `super_admin`
`super_admin` 是默认兜底角色。教程阶段先保护它，避免误操作导致默认管理员无法继续访问系统。

真实项目如果要开放超级管理员权限修改，建议配合二次确认、操作日志和权限恢复方案。
:::

::: details 为什么权限分配用“整体替换”
前端通常会提交当前勾选后的完整权限集合。后端整体替换可以保证数据库最终状态和页面勾选状态一致，验证也更直接。
:::

下一节会继续补齐菜单自身的管理能力：[菜单管理](./menu-management)。
