---
title: 菜单管理
description: "实现后台菜单维护，为动态菜单和按钮权限提供配置入口。"
---

# 菜单管理

前面已经能根据角色返回当前用户可见菜单。现在补齐菜单本身的管理能力：查询菜单树、创建菜单、编辑菜单、禁用菜单和删除菜单。

::: tip 🎯 本节目标
完成后，`super_admin` 可以访问菜单管理接口；系统会初始化菜单管理菜单和按钮；通过接口可以维护目录、菜单和按钮节点。
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
│  │     └─ menus.go
│  └─ router/
│     └─ router.go
```

| 位置 | 用途 |
| --- | --- |
| `internal/handler/system/menus.go` | 菜单管理接口 |
| `internal/router/router.go` | 注册菜单管理路由 |
| `internal/bootstrap/bootstrap.go` | 初始化菜单管理权限和菜单 |

::: info 本节不新增数据库表
菜单管理复用 `sys_menu` 和 `sys_role_menu`。`sys_menu` 保存目录、菜单、按钮；`sys_role_menu` 保存角色拥有哪些菜单和按钮。
:::

## 先区分两个菜单接口

项目里会有两类菜单接口：

| 接口 | 用途 |
| --- | --- |
| `/api/v1/auth/menus` | 当前登录用户可见菜单，用于前端渲染侧边栏和按钮 |
| `/api/v1/system/menus` | 菜单配置管理，用于管理员维护菜单树 |

::: warning ⚠️ 不要把两个接口混在一起
`/api/v1/auth/menus` 要按当前用户角色过滤；`/api/v1/system/menus` 是管理接口，返回系统内菜单配置，只有有权限的管理员才能访问。
:::

## 接口规划

本节先实现 5 个接口：

| 方法 | 路径 | 用途 |
| --- | --- | --- |
| `GET` | `/api/v1/system/menus` | 菜单树 |
| `POST` | `/api/v1/system/menus` | 创建菜单 |
| `POST` | `/api/v1/system/menus/:id/update` | 编辑菜单 |
| `POST` | `/api/v1/system/menus/:id/status` | 修改菜单状态 |
| `POST` | `/api/v1/system/menus/:id/delete` | 删除菜单 |

`sys_menu.type` 继续沿用前面的约定：

| 类型 | 含义 | 示例 |
| --- | --- | --- |
| `1` | 目录 | 系统管理 |
| `2` | 菜单 | 用户管理 |
| `3` | 按钮 | 创建用户 |

## 🛠️ 创建菜单管理 Handler

创建 `server/internal/handler/system/menus.go`。这是新增文件，直接完整写入即可。

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

// MenuAdminHandler 负责后台菜单管理接口。
type MenuAdminHandler struct {
	db  *gorm.DB
	log *zap.Logger
}

// NewMenuAdminHandler 创建菜单管理 Handler。
func NewMenuAdminHandler(db *gorm.DB, log *zap.Logger) *MenuAdminHandler {
	return &MenuAdminHandler{
		db:  db,
		log: log,
	}
}

type createMenuRequest struct {
	ParentID  uint             `json:"parent_id"`
	Type      model.MenuType   `json:"type"`
	Code      string           `json:"code"`
	Title     string           `json:"title"`
	Path      string           `json:"path"`
	Component string           `json:"component"`
	Icon      string           `json:"icon"`
	Sort      int              `json:"sort"`
	Status    model.MenuStatus `json:"status"`
	Remark    string           `json:"remark"`
}

type updateMenuRequest struct {
	ParentID  uint             `json:"parent_id"`
	Type      model.MenuType   `json:"type"`
	Title     string           `json:"title"`
	Path      string           `json:"path"`
	Component string           `json:"component"`
	Icon      string           `json:"icon"`
	Sort      int              `json:"sort"`
	Status    model.MenuStatus `json:"status"`
	Remark    string           `json:"remark"`
}

type updateMenuStatusRequest struct {
	Status model.MenuStatus `json:"status"`
}

type menuAdminResponse struct {
	ID        uint                `json:"id"`
	ParentID  uint                `json:"parent_id"`
	Type      model.MenuType      `json:"type"`
	Code      string              `json:"code"`
	Title     string              `json:"title"`
	Path      string              `json:"path"`
	Component string              `json:"component"`
	Icon      string              `json:"icon"`
	Sort      int                 `json:"sort"`
	Status    model.MenuStatus    `json:"status"`
	Remark    string              `json:"remark"`
	Children  []menuAdminResponse `json:"children,omitempty"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

type menuAdminNode struct {
	menuAdminResponse
	children []*menuAdminNode
}

// Tree 返回完整菜单树。
func (h *MenuAdminHandler) Tree(c *gin.Context) {
	var menus []model.Menu
	if err := h.db.Order("sort ASC, id ASC").Find(&menus).Error; err != nil {
		response.Error(c, apperror.Internal("查询菜单树失败", err), h.log)
		return
	}

	response.Success(c, buildMenuAdminTree(menus))
}

// Create 创建菜单、目录或按钮。
func (h *MenuAdminHandler) Create(c *gin.Context) {
	var req createMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	menu, err := normalizeCreateMenuRequest(req)
	if err != nil {
		response.Error(c, err, h.log)
		return
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := ensureMenuCodeAvailable(tx, menu.Code); err != nil {
			return err
		}

		if err := ensureParentMenuUsable(tx, menu.ParentID, menu.Type, 0); err != nil {
			return err
		}

		return tx.Create(&menu).Error
	})
	if err != nil {
		writeMenuError(c, err, "创建菜单失败", h.log)
		return
	}

	response.Success(c, buildMenuAdminResponse(menu))
}

// Update 编辑菜单基础信息。
func (h *MenuAdminHandler) Update(c *gin.Context) {
	menuID, ok := menuIDParam(c, h.log)
	if !ok {
		return
	}

	var req updateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	update, err := normalizeUpdateMenuRequest(req)
	if err != nil {
		response.Error(c, err, h.log)
		return
	}

	var menu model.Menu
	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&menu, menuID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("菜单不存在")
			}
			return err
		}

		if err := ensureParentMenuUsable(tx, update.ParentID, update.Type, menuID); err != nil {
			return err
		}

		if err := tx.Model(&menu).Updates(map[string]any{
			"parent_id":  update.ParentID,
			"type":       update.Type,
			"title":      update.Title,
			"path":       update.Path,
			"component":  update.Component,
			"icon":       update.Icon,
			"sort":       update.Sort,
			"status":     update.Status,
			"remark":     update.Remark,
		}).Error; err != nil {
			return err
		}

		menu.ParentID = update.ParentID
		menu.Type = update.Type
		menu.Title = update.Title
		menu.Path = update.Path
		menu.Component = update.Component
		menu.Icon = update.Icon
		menu.Sort = update.Sort
		menu.Status = update.Status
		menu.Remark = update.Remark
		return nil
	})
	if err != nil {
		writeMenuError(c, err, "更新菜单失败", h.log)
		return
	}

	response.Success(c, buildMenuAdminResponse(menu))
}

// UpdateStatus 修改菜单状态。
func (h *MenuAdminHandler) UpdateStatus(c *gin.Context) {
	menuID, ok := menuIDParam(c, h.log)
	if !ok {
		return
	}

	var req updateMenuStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	if !validMenuStatus(req.Status) {
		response.Error(c, apperror.BadRequest("菜单状态不正确"), h.log)
		return
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		var menu model.Menu
		if err := tx.First(&menu, menuID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("菜单不存在")
			}
			return err
		}

		return tx.Model(&menu).Update("status", req.Status).Error
	})
	if err != nil {
		writeMenuError(c, err, "更新菜单状态失败", h.log)
		return
	}

	response.Success(c, gin.H{
		"id":     menuID,
		"status": req.Status,
	})
}

// Delete 删除菜单。
func (h *MenuAdminHandler) Delete(c *gin.Context) {
	menuID, ok := menuIDParam(c, h.log)
	if !ok {
		return
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		var menu model.Menu
		if err := tx.First(&menu, menuID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("菜单不存在")
			}
			return err
		}

		if err := ensureMenuCanDelete(tx, menuID); err != nil {
			return err
		}

		return tx.Delete(&menu).Error
	})
	if err != nil {
		writeMenuError(c, err, "删除菜单失败", h.log)
		return
	}

	response.Success(c, gin.H{
		"id": menuID,
	})
}

func normalizeCreateMenuRequest(req createMenuRequest) (model.Menu, error) {
	code := strings.TrimSpace(req.Code)
	if code == "" {
		return model.Menu{}, apperror.BadRequest("菜单编码不能为空")
	}
	if len(code) > 128 {
		return model.Menu{}, apperror.BadRequest("菜单编码不能超过 128 个字符")
	}

	title, path, component, icon, status, remark, err := normalizeMenuFields(
		req.Type,
		req.Title,
		req.Path,
		req.Component,
		req.Icon,
		req.Status,
		req.Remark,
	)
	if err != nil {
		return model.Menu{}, err
	}

	return model.Menu{
		ParentID:  req.ParentID,
		Type:      req.Type,
		Code:      code,
		Title:     title,
		Path:      path,
		Component: component,
		Icon:      icon,
		Sort:      req.Sort,
		Status:    status,
		Remark:    remark,
	}, nil
}

func normalizeUpdateMenuRequest(req updateMenuRequest) (model.Menu, error) {
	title, path, component, icon, status, remark, err := normalizeMenuFields(
		req.Type,
		req.Title,
		req.Path,
		req.Component,
		req.Icon,
		req.Status,
		req.Remark,
	)
	if err != nil {
		return model.Menu{}, err
	}

	return model.Menu{
		ParentID:  req.ParentID,
		Type:      req.Type,
		Title:     title,
		Path:      path,
		Component: component,
		Icon:      icon,
		Sort:      req.Sort,
		Status:    status,
		Remark:    remark,
	}, nil
}

func normalizeMenuFields(menuType model.MenuType, title string, path string, component string, icon string, status model.MenuStatus, remark string) (string, string, string, string, model.MenuStatus, string, error) {
	if !validMenuType(menuType) {
		return "", "", "", "", 0, "", apperror.BadRequest("菜单类型不正确")
	}

	title = strings.TrimSpace(title)
	if title == "" {
		return "", "", "", "", 0, "", apperror.BadRequest("菜单名称不能为空")
	}
	if len(title) > 64 {
		return "", "", "", "", 0, "", apperror.BadRequest("菜单名称不能超过 64 个字符")
	}

	path = strings.TrimSpace(path)
	component = strings.TrimSpace(component)
	icon = strings.TrimSpace(icon)
	remark = strings.TrimSpace(remark)

	if len(path) > 255 {
		return "", "", "", "", 0, "", apperror.BadRequest("路由路径不能超过 255 个字符")
	}
	if len(component) > 255 {
		return "", "", "", "", 0, "", apperror.BadRequest("组件路径不能超过 255 个字符")
	}
	if len(icon) > 64 {
		return "", "", "", "", 0, "", apperror.BadRequest("图标标识不能超过 64 个字符")
	}
	if len(remark) > 255 {
		return "", "", "", "", 0, "", apperror.BadRequest("备注不能超过 255 个字符")
	}

	if status == 0 {
		status = model.MenuStatusEnabled
	}
	if !validMenuStatus(status) {
		return "", "", "", "", 0, "", apperror.BadRequest("菜单状态不正确")
	}

	if menuType == model.MenuTypeMenu && path == "" {
		return "", "", "", "", 0, "", apperror.BadRequest("菜单节点需要填写路由路径")
	}

	return title, path, component, icon, status, remark, nil
}

func ensureMenuCodeAvailable(db *gorm.DB, code string) error {
	var menu model.Menu
	err := db.Unscoped().Where("code = ?", code).First(&menu).Error
	if err == nil {
		return apperror.BadRequest("菜单编码已存在")
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return err
}

func ensureParentMenuUsable(db *gorm.DB, parentID uint, menuType model.MenuType, currentID uint) error {
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

func ensureMenuCanDelete(db *gorm.DB, menuID uint) error {
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

func validMenuType(menuType model.MenuType) bool {
	return menuType == model.MenuTypeDirectory ||
		menuType == model.MenuTypeMenu ||
		menuType == model.MenuTypeButton
}

func validMenuStatus(status model.MenuStatus) bool {
	return status == model.MenuStatusEnabled || status == model.MenuStatusDisabled
}

func menuIDParam(c *gin.Context, log *zap.Logger) (uint, bool) {
	rawID := c.Param("id")
	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil || id == 0 {
		response.Error(c, apperror.BadRequest("菜单 ID 不正确"), log)
		return 0, false
	}

	return uint(id), true
}

func buildMenuAdminTree(menus []model.Menu) []menuAdminResponse {
	nodes := make(map[uint]*menuAdminNode, len(menus))

	for _, menu := range menus {
		nodes[menu.ID] = &menuAdminNode{
			menuAdminResponse: buildMenuAdminResponse(menu),
		}
	}

	roots := make([]*menuAdminNode, 0)
	for _, menu := range menus {
		node := nodes[menu.ID]
		if menu.ParentID == 0 {
			roots = append(roots, node)
			continue
		}

		parent, ok := nodes[menu.ParentID]
		if !ok {
			roots = append(roots, node)
			continue
		}

		parent.children = append(parent.children, node)
	}

	return menuAdminNodesToResponses(roots)
}

func menuAdminNodesToResponses(nodes []*menuAdminNode) []menuAdminResponse {
	result := make([]menuAdminResponse, 0, len(nodes))
	for _, node := range nodes {
		item := node.menuAdminResponse
		item.Children = menuAdminNodesToResponses(node.children)
		result = append(result, item)
	}

	return result
}

func buildMenuAdminResponse(menu model.Menu) menuAdminResponse {
	return menuAdminResponse{
		ID:        menu.ID,
		ParentID:  menu.ParentID,
		Type:      menu.Type,
		Code:      menu.Code,
		Title:     menu.Title,
		Path:      menu.Path,
		Component: menu.Component,
		Icon:      menu.Icon,
		Sort:      menu.Sort,
		Status:    menu.Status,
		Remark:    menu.Remark,
		CreatedAt: menu.CreatedAt,
		UpdatedAt: menu.UpdatedAt,
	}
}

func writeMenuError(c *gin.Context, err error, fallbackMessage string, log *zap.Logger) {
	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		response.Error(c, appErr, log)
		return
	}

	response.Error(c, apperror.Internal(fallbackMessage, err), log)
}
```

::: details 为什么创建后不允许修改 `code`
`code` 是前端判断按钮权限时最稳定的标识，也会出现在角色授权配置里。创建后如果随意改编码，前端按钮判断和历史授权会变得难排查。
:::

::: details 为什么删除菜单前要检查子菜单和角色绑定
本项目不使用数据库外键，关联约束要放在业务逻辑里维护。删除菜单前先检查子节点和 `sys_role_menu`，可以避免留下不可见但仍被授权的数据。
:::

## 🛠️ 注册菜单管理路由

修改 `server/internal/router/router.go`。这一处在系统路由中新增菜单管理 Handler 和路由。

```go
// registerSystemRoutes 注册系统级路由。
func registerSystemRoutes(r *gin.Engine, opts Options) {
	health := systemHandler.NewHealthHandler(opts.Config, opts.DB, opts.Redis, opts.Log)
	users := systemHandler.NewUserHandler(opts.DB, opts.Log)
	roles := systemHandler.NewRoleHandler(opts.DB, opts.Log)
	menus := systemHandler.NewMenuAdminHandler(opts.DB, opts.Log) // [!code ++]

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
	system.GET("/roles", roles.List)
	system.POST("/roles", roles.Create)
	system.POST("/roles/:id/update", roles.Update)
	system.POST("/roles/:id/status", roles.UpdateStatus)
	system.POST("/roles/:id/permissions", roles.UpdatePermissions)
	system.POST("/roles/:id/menus", roles.UpdateMenus)
	system.GET("/menus", menus.Tree) // [!code ++]
	system.POST("/menus", menus.Create) // [!code ++]
	system.POST("/menus/:id/update", menus.Update) // [!code ++]
	system.POST("/menus/:id/status", menus.UpdateStatus) // [!code ++]
	system.POST("/menus/:id/delete", menus.Delete) // [!code ++]
}
```

## 🛠️ 初始化菜单管理接口权限

修改 `server/internal/bootstrap/bootstrap.go`。在 `defaultPermissionSeeds` 中继续追加菜单管理接口权限：

```go
var defaultPermissionSeeds = []defaultPermissionSeed{
	{Path: "/api/v1/system/health", Method: "GET"},
	{Path: "/api/v1/system/users", Method: "GET"},
	{Path: "/api/v1/system/users", Method: "POST"},
	{Path: "/api/v1/system/users/:id/update", Method: "POST"},
	{Path: "/api/v1/system/users/:id/status", Method: "POST"},
	{Path: "/api/v1/system/users/:id/roles", Method: "POST"},
	{Path: "/api/v1/system/roles", Method: "GET"},
	{Path: "/api/v1/system/roles", Method: "POST"},
	{Path: "/api/v1/system/roles/:id/update", Method: "POST"},
	{Path: "/api/v1/system/roles/:id/status", Method: "POST"},
	{Path: "/api/v1/system/roles/:id/permissions", Method: "POST"},
	{Path: "/api/v1/system/roles/:id/menus", Method: "POST"},
	{Path: "/api/v1/system/menus", Method: "GET"}, // [!code ++]
	{Path: "/api/v1/system/menus", Method: "POST"}, // [!code ++]
	{Path: "/api/v1/system/menus/:id/update", Method: "POST"}, // [!code ++]
	{Path: "/api/v1/system/menus/:id/status", Method: "POST"}, // [!code ++]
	{Path: "/api/v1/system/menus/:id/delete", Method: "POST"}, // [!code ++]
}
```

## 🛠️ 初始化菜单管理菜单

继续修改 `server/internal/bootstrap/bootstrap.go`。先增加菜单管理菜单和按钮编码：

```go
const (
	defaultRoleMenuAssignCode = "system:role:menu"
	defaultMenuManageCode     = "system:menu" // [!code ++]
	defaultMenuListCode       = "system:menu:list" // [!code ++]
	defaultMenuCreateCode     = "system:menu:create" // [!code ++]
	defaultMenuUpdateCode     = "system:menu:update" // [!code ++]
	defaultMenuStatusCode     = "system:menu:status" // [!code ++]
	defaultMenuDeleteCode     = "system:menu:delete" // [!code ++]
)
```

接着修改 `seedDefaultMenus`。先找到上一节最后新增的返回语句：

```go
return menus, nil
```

把这行返回语句替换为下面整段代码。也就是说：下面代码放在角色管理按钮循环之后，原 `return menus, nil` 之前；替换完成后，函数末尾仍然只保留一个 `return menus, nil`。

```go
	menuManage, err := seedMenu(db, model.Menu{
		ParentID:  systemMenu.ID,
		Type:      model.MenuTypeMenu,
		Code:      defaultMenuManageCode,
		Title:     "菜单管理",
		Path:      "/system/menus",
		Component: "system/MenuView",
		Icon:      "menu",
		Sort:      40,
		Status:    model.MenuStatusEnabled,
		Remark:    "系统内置菜单",
	}, log)
	if err != nil {
		return nil, err
	}

	menuButtons := []model.Menu{
		{ParentID: menuManage.ID, Type: model.MenuTypeButton, Code: defaultMenuListCode, Title: "查看菜单", Sort: 10, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
		{ParentID: menuManage.ID, Type: model.MenuTypeButton, Code: defaultMenuCreateCode, Title: "创建菜单", Sort: 20, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
		{ParentID: menuManage.ID, Type: model.MenuTypeButton, Code: defaultMenuUpdateCode, Title: "编辑菜单", Sort: 30, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
		{ParentID: menuManage.ID, Type: model.MenuTypeButton, Code: defaultMenuStatusCode, Title: "修改菜单状态", Sort: 40, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
		{ParentID: menuManage.ID, Type: model.MenuTypeButton, Code: defaultMenuDeleteCode, Title: "删除菜单", Sort: 50, Status: model.MenuStatusEnabled, Remark: "系统内置按钮"},
	}

	menus = append(menus, *menuManage)
	for _, button := range menuButtons {
		createdButton, err := seedMenu(db, button, log)
		if err != nil {
			return nil, err
		}
		menus = append(menus, *createdButton)
	}

	return menus, nil
```

::: warning ⚠️ 这一段仍然是在 `seedDefaultMenus` 函数内部
菜单管理菜单要继续追加到已有的 `menus` 切片里，不要新建第二个 `menus`，也不要把代码放到 `seedDefaultMenus` 函数外。
:::

::: warning ⚠️ 这里的变量名容易撞
前面已经有 `/api/v1/auth/menus` 的 `menus` Handler。为了避免在 `system` 包里读起来混乱，本节把管理端 Handler 命名为 `MenuAdminHandler`。
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
INFO	default permission created	{"role_code": "super_admin", "path": "/api/v1/system/menus", "method": "GET"}
INFO	default menu created	{"menu_code": "system:menu"}
INFO	default role menu bound	{"role_id": 1, "menu_id": 18}
```

## ✅ 验证权限和菜单数据

先确认菜单管理接口权限已经写入：

```bash
# 查看菜单管理相关接口权限
docker compose -f deploy/compose.local.yml exec postgres psql -U ez_admin -d ez_admin -c "select ptype, v0, v1, v2 from casbin_rule where v1 like '/api/v1/system/menus%' order by v1, v2;"
```

应该能看到菜单树、创建、编辑、状态修改、删除对应的策略。

再确认菜单管理菜单已经写入：

```bash
# 查看菜单管理相关菜单和按钮
docker compose -f deploy/compose.local.yml exec postgres psql -U ez_admin -d ez_admin -c "select id, parent_id, type, code, title from sys_menu where code like 'system:menu%' order by sort, id;"
```

应该能看到 `system:menu` 以及几个 `system:menu:*` 按钮编码。

## ✅ 验证菜单管理接口

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

查看菜单树：

::: code-group

```powershell [Windows PowerShell]
Invoke-RestMethod `
  -Method Get `
  -Uri http://localhost:8080/api/v1/system/menus `
  -Headers @{ Authorization = "Bearer $token" }
```

```bash [macOS / Linux]
curl http://localhost:8080/api/v1/system/menus \
  -H "Authorization: Bearer ${TOKEN}"
```

:::

应该能看到完整菜单树，包括目录、菜单和按钮。

创建一个测试菜单。下面示例假设 `系统管理` 的菜单 ID 是 `1`，如果不确定，可以先查询 `sys_menu`：

::: warning ⚠️ Windows PowerShell 发送中文 JSON 时要显式使用 UTF-8
请求体中包含中文时，先把 JSON 转成 UTF-8 字节再发送，避免中文在发送阶段变成 `????`。
:::

::: code-group

```powershell [Windows PowerShell]
$body = @{
  parent_id = 1
  type = 2
  code = "system:demo"
  title = "演示菜单"
  path = "/system/demo"
  component = "system/DemoView"
  icon = "experiment"
  sort = 90
  status = 1
  remark = "用于验证菜单管理接口"
} | ConvertTo-Json

$utf8Body = [System.Text.Encoding]::UTF8.GetBytes($body)

Invoke-RestMethod `
  -Method Post `
  -Uri http://localhost:8080/api/v1/system/menus `
  -ContentType "application/json; charset=utf-8" `
  -Headers @{ Authorization = "Bearer $token" } `
  -Body $utf8Body
```

```bash [macOS / Linux]
curl -X POST http://localhost:8080/api/v1/system/menus \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"parent_id":1,"type":2,"code":"system:demo","title":"演示菜单","path":"/system/demo","component":"system/DemoView","icon":"experiment","sort":90,"status":1,"remark":"用于验证菜单管理接口"}'
```

:::

创建成功后，用 SQL 确认菜单已经写入：

```bash
# 查看演示菜单
docker compose -f deploy/compose.local.yml exec postgres psql -U ez_admin -d ez_admin -c "select id, parent_id, type, code, title, path, status from sys_menu where code = 'system:demo';"
```

修改菜单状态。把上一步返回的菜单 ID 替换到路径里：

::: code-group

```powershell [Windows PowerShell]
$menuId = 20
$body = @{ status = 2 } | ConvertTo-Json

Invoke-RestMethod `
  -Method Post `
  -Uri "http://localhost:8080/api/v1/system/menus/$menuId/status" `
  -ContentType "application/json" `
  -Headers @{ Authorization = "Bearer $token" } `
  -Body $body
```

```bash [macOS / Linux]
MENU_ID=20

curl -X POST "http://localhost:8080/api/v1/system/menus/${MENU_ID}/status" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"status":2}'
```

:::

`status = 2` 表示禁用。禁用后，这个菜单不会出现在当前用户菜单接口中。

删除测试菜单：

::: code-group

```powershell [Windows PowerShell]
$menuId = 20

Invoke-RestMethod `
  -Method Post `
  -Uri "http://localhost:8080/api/v1/system/menus/$menuId/delete" `
  -Headers @{ Authorization = "Bearer $token" }
```

```bash [macOS / Linux]
MENU_ID=20

curl -X POST "http://localhost:8080/api/v1/system/menus/${MENU_ID}/delete" \
  -H "Authorization: Bearer ${TOKEN}"
```

:::

::: warning ⚠️ 不要直接删除已经分配给角色的菜单
如果菜单已经存在子菜单，或者已经写入 `sys_role_menu`，删除接口会拒绝操作。验证删除逻辑时，建议使用刚创建且未分配给角色的测试菜单。
:::

## 常见问题

::: details 创建菜单时提示“根节点只能是目录”
`parent_id = 0` 表示根节点。根节点只能创建目录，菜单和按钮都应该挂在已有父级下面。
:::

::: details 创建按钮时提示“按钮只能挂在菜单下面”
按钮代表页面内操作点，应该挂在具体菜单下面。目录下面可以挂目录或菜单，菜单下面可以挂按钮。
:::

::: details 删除菜单时提示“请先删除子菜单”
先删除或迁移它的子节点，再删除当前菜单。这样可以避免菜单树出现断层。
:::

::: details 删除菜单时提示“菜单已分配给角色，不能删除”
先在角色管理中取消这个菜单的角色绑定，再删除菜单。
:::

下一节会继续补齐系统配置能力：[系统配置](./system-config)。
