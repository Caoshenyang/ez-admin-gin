---
title: 示例业务模块
description: "用公告管理模块走完一整条接入链路，证明前面定义的规范可以落地。"
---

# 示例业务模块

前四页已经把模块接入的每一步拆开了：目录放哪、后端怎么接、权限菜单怎么挂、前端页面怎么写。但拆开看和串起来跑是两件事。这一页用一个完整的公告管理模块，从 Model 到页面，把前面所有约定串成一条能跑通的链路。

::: tip 这一页做完你能得到什么
一个完整的公告管理模块，包含后端接口、数据库迁移、权限菜单种子和前端 CRUD 页面。更重要的是，你会看到前面定义的每一条规范在真实代码里是怎么落地的，以后照着这个模式接新模块即可。
:::

## 为什么选公告管理

公告管理是后台系统里常见的轻量模块，数据结构简单、操作清晰，但同时又覆盖了分页查询、关键字搜索、状态切换、新建编辑这些后台页面最常见的交互。用它做示例，既能讲清楚接入流程，又不会因为业务本身太复杂而分散注意力。

## 后端：Model

公告表需要记录标题、正文、排序、状态和备注，同时支持软删除。下面是完整的 Model 定义：

```go
package model

import (
	"time"

	"gorm.io/gorm"
)

// NoticeStatus 表示公告状态。
type NoticeStatus int

const (
	// NoticeStatusEnabled 表示公告可见。
	NoticeStatusEnabled NoticeStatus = 1
	// NoticeStatusDisabled 表示公告已隐藏。
	NoticeStatusDisabled NoticeStatus = 2
)

// Notice 是公告表模型。
type Notice struct {
	ID        uint          `gorm:"primaryKey" json:"id"`
	Title     string        `gorm:"size:128;not null" json:"title"`
	Content   string        `gorm:"type:text;not null" json:"content"`
	Sort      int           `gorm:"not null;default:0" json:"sort"`
	Status    NoticeStatus  `gorm:"type:smallint;not null;default:1" json:"status"`
	Remark    string        `gorm:"size:255;not null;default:''" json:"remark"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 固定公告表名。
func (Notice) TableName() string {
	return "sys_notice"
}
```

几个要点：

- `TableName()` 固定表名为 `sys_notice`，与系统表保持 `sys_` 前缀一致。
- `DeletedAt` 使用 `gorm.DeletedAt`，GORM 会自动处理软删除，`json:"-"` 表示不返回给前端。
- `Status` 使用自定义类型 `NoticeStatus`，配合常量 `Enabled = 1` / `Disabled = 2`，让代码语义更清晰。
- 排序字段 `Sort` 默认为 `0`，列表查询时按 `sort ASC, id DESC` 排序。

## 后端：Handler

公告 Handler 包含四个方法：`List`、`Create`、`Update`、`UpdateStatus`，对应分页查询、新建、编辑和状态变更。文件较长，折叠查看：

::: details `server/internal/handler/system/notices.go` — 公告 Handler 完整实现
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

// NoticeHandler 负责公告管理接口。
type NoticeHandler struct {
	db  *gorm.DB
	log *zap.Logger
}

// NewNoticeHandler 创建公告 Handler。
func NewNoticeHandler(db *gorm.DB, log *zap.Logger) *NoticeHandler {
	return &NoticeHandler{db: db, log: log}
}

type noticeListQuery struct {
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
	Keyword   string `form:"keyword"`
	Status    int    `form:"status"`
}

type createNoticeRequest struct {
	Title   string          `json:"title"`
	Content string          `json:"content"`
	Sort    int             `json:"sort"`
	Status  model.NoticeStatus `json:"status"`
	Remark  string          `json:"remark"`
}

type updateNoticeRequest struct {
	Title   string          `json:"title"`
	Content string          `json:"content"`
	Sort    int             `json:"sort"`
	Status  model.NoticeStatus `json:"status"`
	Remark  string          `json:"remark"`
}

type updateNoticeStatusRequest struct {
	Status model.NoticeStatus `json:"status"`
}

type noticeResponse struct {
	ID        uint             `json:"id"`
	Title     string           `json:"title"`
	Content   string           `json:"content"`
	Sort      int              `json:"sort"`
	Status    model.NoticeStatus `json:"status"`
	Remark    string           `json:"remark"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

type noticeListResponse struct {
	Items    []noticeResponse `json:"items"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
}

// List 返回公告分页列表。
func (h *NoticeHandler) List(c *gin.Context) {
	var query noticeListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, apperror.BadRequest("查询参数不正确"), h.log)
		return
	}

	page, pageSize := normalizeNoticePage(query.Page, query.PageSize)
	queryDB := h.db.Model(&model.Notice{})

	keyword := strings.TrimSpace(query.Keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		queryDB = queryDB.Where("title LIKE ?", like)
	}

	if query.Status != 0 {
		status := model.NoticeStatus(query.Status)
		if !validNoticeStatus(status) {
			response.Error(c, apperror.BadRequest("公告状态不正确"), h.log)
			return
		}
		queryDB = queryDB.Where("status = ?", status)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		response.Error(c, apperror.Internal("查询公告总数失败", err), h.log)
		return
	}

	var notices []model.Notice
	if err := queryDB.
		Order("sort ASC, id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&notices).Error; err != nil {
		response.Error(c, apperror.Internal("查询公告列表失败", err), h.log)
		return
	}

	items := make([]noticeResponse, 0, len(notices))
	for _, n := range notices {
		items = append(items, buildNoticeResponse(n))
	}

	response.Success(c, noticeListResponse{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// Create 创建公告。
func (h *NoticeHandler) Create(c *gin.Context) {
	var req createNoticeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		response.Error(c, apperror.BadRequest("公告标题不能为空"), h.log)
		return
	}

	if len(title) > 128 {
		response.Error(c, apperror.BadRequest("公告标题不能超过 128 个字符"), h.log)
		return
	}

	status := req.Status
	if status == 0 {
		status = model.NoticeStatusEnabled
	}
	if !validNoticeStatus(status) {
		response.Error(c, apperror.BadRequest("公告状态不正确"), h.log)
		return
	}

	notice := model.Notice{
		Title:   title,
		Content: req.Content,
		Sort:    req.Sort,
		Status:  status,
		Remark:  strings.TrimSpace(req.Remark),
	}

	if err := h.db.Create(&notice).Error; err != nil {
		response.Error(c, apperror.Internal("创建公告失败", err), h.log)
		return
	}

	response.Success(c, buildNoticeResponse(notice))
}

// Update 编辑公告。
func (h *NoticeHandler) Update(c *gin.Context) {
	noticeID, ok := noticeIDParam(c, h.log)
	if !ok {
		return
	}

	var req updateNoticeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		response.Error(c, apperror.BadRequest("公告标题不能为空"), h.log)
		return
	}

	if !validNoticeStatus(req.Status) {
		response.Error(c, apperror.BadRequest("公告状态不正确"), h.log)
		return
	}

	var notice model.Notice
	if err := h.db.First(&notice, noticeID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, apperror.NotFound("公告不存在"), h.log)
			return
		}
		response.Error(c, apperror.Internal("查询公告失败", err), h.log)
		return
	}

	if err := h.db.Model(&notice).Updates(map[string]any{
		"title":   title,
		"content": req.Content,
		"sort":    req.Sort,
		"status":  req.Status,
		"remark":  strings.TrimSpace(req.Remark),
	}).Error; err != nil {
		response.Error(c, apperror.Internal("更新公告失败", err), h.log)
		return
	}

	notice.Title = title
	notice.Content = req.Content
	notice.Sort = req.Sort
	notice.Status = req.Status
	notice.Remark = strings.TrimSpace(req.Remark)

	response.Success(c, buildNoticeResponse(notice))
}

// UpdateStatus 修改公告状态。
func (h *NoticeHandler) UpdateStatus(c *gin.Context) {
	noticeID, ok := noticeIDParam(c, h.log)
	if !ok {
		return
	}

	var req updateNoticeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	if !validNoticeStatus(req.Status) {
		response.Error(c, apperror.BadRequest("公告状态不正确"), h.log)
		return
	}

	var notice model.Notice
	if err := h.db.First(&notice, noticeID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, apperror.NotFound("公告不存在"), h.log)
			return
		}
		response.Error(c, apperror.Internal("查询公告失败", err), h.log)
		return
	}

	if err := h.db.Model(&notice).Update("status", req.Status).Error; err != nil {
		response.Error(c, apperror.Internal("更新公告状态失败", err), h.log)
		return
	}

	response.Success(c, gin.H{
		"id":     noticeID,
		"status": req.Status,
	})
}

func normalizeNoticePage(page int, pageSize int) (int, int) {
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

func validNoticeStatus(status model.NoticeStatus) bool {
	return status == model.NoticeStatusEnabled || status == model.NoticeStatusDisabled
}

func noticeIDParam(c *gin.Context, log *zap.Logger) (uint, bool) {
	rawID := c.Param("id")
	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil || id == 0 {
		response.Error(c, apperror.BadRequest("公告 ID 不正确"), log)
		return 0, false
	}
	return uint(id), true
}

func buildNoticeResponse(n model.Notice) noticeResponse {
	return noticeResponse{
		ID:        n.ID,
		Title:     n.Title,
		Content:   n.Content,
		Sort:      n.Sort,
		Status:    n.Status,
		Remark:    n.Remark,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}
}
```
:::

Handler 的写法和前面系统模块完全一致，值得关注的几个设计：

- **请求 / 响应结构体定义在文件内部**。`noticeListQuery`、`createNoticeRequest` 等结构体只在 Handler 里使用，不需要对外暴露，所以不放到 Model 包。
- **分页参数归一化**。`normalizeNoticePage` 把非法的 `page` 和 `page_size` 修正为合理值，上限 100，避免一次查太多数据。
- **关键字搜索用 `LIKE`**。公告量通常不大，`LIKE` 足够；如果后续数据量变大，可以换全文检索。
- **`buildNoticeResponse` 统一响应格式**。从 Model 转成响应结构体时集中在一个函数里处理，后续加字段只需要改一处。
- **`Update` 用 `map[string]any` 做批量更新**。GORM 的 `Updates` 方法传入 struct 时会忽略零值字段，用 map 可以避免这个问题。

::: warning 为什么 UpdateStatus 单独拆一个方法
状态变更是高频操作，而且只需要传一个字段。如果复用 Update 方法，前端每次切换状态都要把公告的全部字段回传，既浪费带宽又容易出错。拆出来后，状态切换只需传 `status` 一个值，接口更轻。
:::

## 后端：Router

路由注册只需要在 `registerSystemRoutes` 里新增两行：创建 Handler 实例，注册路由。下面用 diff 标记标出新增的部分：

```go
func registerSystemRoutes(r *gin.Engine, opts Options) {
    health := systemHandler.NewHealthHandler(opts.Config, opts.DB, opts.Redis, opts.Log)
    users := systemHandler.NewUserHandler(opts.DB, opts.Log)
    roles := systemHandler.NewRoleHandler(opts.DB, opts.Log)
    menus := systemHandler.NewMenuAdminHandler(opts.DB, opts.Log)
    configs := systemHandler.NewSystemConfigHandler(opts.DB, opts.Redis, opts.Log)
    files := systemHandler.NewFileHandler(opts.DB, opts.Config.Upload, opts.Log)
    operationLogs := systemHandler.NewOperationLogHandler(opts.DB, opts.Log)
    loginLogs := systemHandler.NewLoginLogHandler(opts.DB, opts.Log)
    notices := systemHandler.NewNoticeHandler(opts.DB, opts.Log) // [!code ++]

    // ... 省略中间代码 ...

    system.GET("/login-logs", loginLogs.List)
    system.GET("/notices", notices.List)              // [!code ++]
    system.POST("/notices", notices.Create)            // [!code ++]
    system.POST("/notices/:id/update", notices.Update) // [!code ++]
    system.POST("/notices/:id/status", notices.UpdateStatus) // [!code ++]
}
```

公告路由注册在 `system` 分组下，自动继承了三条中间件：

1. **Auth** — 验证登录状态，未登录返回 401。
2. **OperationLog** — 记录操作日志，方便审计。
3. **Permission** — 校验角色是否有对应接口的访问权限。

::: details 路由路径为什么要统一用复数
`/notices` 而不是 `/notice`，与已有的 `/users`、`/roles`、`/menus` 保持一致。REST 风格里资源名用复数是常见约定，团队统一一种写法比争论哪一种更正确更有价值。
:::

## 后端：数据库迁移

公告模块需要通过数据库迁移文件来初始化权限种子和菜单种子。创建新的迁移文件来添加公告管理的权限和菜单：

### 权限种子

在 `server/migrations/{postgres,mysql}/` 目录下创建新的迁移文件，添加公告的接口权限：

::: code-group

```sql [PostgreSQL — 000003_notice_seed_data.up.sql]
-- 公告管理接口权限
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/notices', 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/notices', 'POST');
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/notices/:id/update', 'POST');
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', '/api/v1/system/notices/:id/status', 'POST');
```

```sql [MySQL — 000003_notice_seed_data.up.sql]
-- 公告管理接口权限
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/notices', 'GET');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/notices', 'POST');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/notices/:id/update', 'POST');
INSERT INTO `casbin_rule` (`ptype`, `v0`, `v1`, `v2`) VALUES ('p', 'super_admin', '/api/v1/system/notices/:id/status', 'POST');
```

:::

### 菜单种子

在同一个迁移文件中，添加公告菜单和按钮：

::: code-group

```sql [PostgreSQL — 000003_notice_seed_data.up.sql]
-- 公告管理菜单
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (20, 1, 2, 'system:notice', '公告管理', '/system/notices', 'system/NoticeView', 'notification', 90, 1, '系统内置菜单', NOW(), NOW());

-- 公告管理按钮
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1000, 20, 3, 'system:notice:list', '查看公告', '', '', '', 10, 1, '系统内置按钮', NOW(), NOW());
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1001, 20, 3, 'system:notice:create', '创建公告', '', '', '', 20, 1, '系统内置按钮', NOW(), NOW());
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1002, 20, 3, 'system:notice:update', '编辑公告', '', '', '', 30, 1, '系统内置按钮', NOW(), NOW());
INSERT INTO sys_menu (id, parent_id, type, code, title, path, component, icon, sort, status, remark, created_at, updated_at)
VALUES (1003, 20, 3, 'system:notice:status', '修改公告状态', '', '', '', 40, 1, '系统内置按钮', NOW(), NOW());

-- 绑定到 super_admin 角色
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 20, NOW(), NOW());
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1000, NOW(), NOW());
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1001, NOW(), NOW());
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1002, NOW(), NOW());
INSERT INTO sys_role_menu (role_id, menu_id, created_at, updated_at) VALUES (1, 1003, NOW(), NOW());
```

```sql [MySQL — 000003_notice_seed_data.up.sql]
-- 公告管理菜单
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (20, 1, 2, 'system:notice', '公告管理', '/system/notices', 'system/NoticeView', 'notification', 90, 1, '系统内置菜单', NOW(), NOW());

-- 公告管理按钮
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1000, 20, 3, 'system:notice:list', '查看公告', '', '', '', 10, 1, '系统内置按钮', NOW(), NOW());
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1001, 20, 3, 'system:notice:create', '创建公告', '', '', '', 20, 1, '系统内置按钮', NOW(), NOW());
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1002, 20, 3, 'system:notice:update', '编辑公告', '', '', '', 30, 1, '系统内置按钮', NOW(), NOW());
INSERT INTO `sys_menu` (`id`, `parent_id`, `type`, `code`, `title`, `path`, `component`, `icon`, `sort`, `status`, `remark`, `created_at`, `updated_at`)
VALUES (1003, 20, 3, 'system:notice:status', '修改公告状态', '', '', '', 40, 1, '系统内置按钮', NOW(), NOW());

-- 绑定到 super_admin 角色
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 20, NOW(), NOW());
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1000, NOW(), NOW());
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1001, NOW(), NOW());
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1002, NOW(), NOW());
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`, `created_at`, `updated_at`) VALUES (1, 1003, NOW(), NOW());
```

:::

::: warning 菜单种子的 `Component` 字段必须与前端路由映射一致
迁移文件里写的 `component: "system/NoticeView"` 必须和前端 `dynamic-menu.ts` 中 `routeComponentMap` 的 key 完全匹配。如果这里写 `Notice` 而前端写 `system/NoticeView`，菜单能查到但页面会加载占位组件，不会报错但也不会显示真实页面。这类问题排查起来很费时间，建议在接入新模块时把 `Component` 值直接从前端 `routeComponentMap` 里复制过来。
:::

## 前端：Types

类型定义是前端接入的起点。公告模块的类型文件包含状态枚举、列表项、查询参数、响应结构和请求载荷：

```ts
export const NoticeStatus = {
  Enabled: 1,
  Disabled: 2,
} as const

export type NoticeStatus = (typeof NoticeStatus)[keyof typeof NoticeStatus]

export interface NoticeItem {
  id: number
  title: string
  content: string
  sort: number
  status: NoticeStatus
  remark: string
  created_at: string
  updated_at: string
}

export interface NoticeListQuery {
  page: number
  page_size: number
  keyword?: string
  status?: NoticeStatus | 0
}

export interface NoticeListResponse {
  items: NoticeItem[]
  total: number
  page: number
  page_size: number
}

export interface CreateNoticePayload {
  title: string
  content: string
  sort: number
  status: NoticeStatus
  remark: string
}

export interface UpdateNoticePayload {
  title: string
  content: string
  sort: number
  status: NoticeStatus
  remark: string
}

export interface UpdateNoticeStatusPayload {
  status: NoticeStatus
}
```

类型定义和后端 Model 一一对应，几个设计考虑：

- `NoticeStatus` 用 `as const` 定义常量对象，同时导出类型和值，在模板和逻辑中都能直接使用。
- `NoticeListQuery` 的 `status` 类型写成 `NoticeStatus | 0`，`0` 表示"查询全部"，不传给后端。
- `CreateNoticePayload` 和 `UpdateNoticePayload` 结构相同，但分开定义。如果后续创建和编辑的字段出现差异（比如编辑时多一个版本号），改动不会互相影响。

## 前端：API

API 层负责类型安全的请求封装，每个函数对应一个后端接口：

```ts
import http from './http'

import type { ApiResponse } from '../types/http'
import type {
  NoticeItem,
  NoticeListQuery,
  NoticeListResponse,
  CreateNoticePayload,
  UpdateNoticePayload,
  UpdateNoticeStatusPayload,
} from '../types/notice'

export async function getNotices(params: NoticeListQuery) {
  const response = await http.get<ApiResponse<NoticeListResponse>>('/system/notices', { params })
  return response.data.data
}

export async function createNotice(payload: CreateNoticePayload) {
  const response = await http.post<ApiResponse<NoticeItem>>('/system/notices', payload)
  return response.data.data
}

export async function updateNotice(id: number, payload: UpdateNoticePayload) {
  const response = await http.post<ApiResponse<NoticeItem>>(`/system/notices/${id}/update`, payload)
  return response.data.data
}

export async function updateNoticeStatus(id: number, payload: UpdateNoticeStatusPayload) {
  const response = await http.post<ApiResponse<{ id: number; status: number }>>(
    `/system/notices/${id}/status`,
    payload,
  )
  return response.data.data
}
```

注意接口路径和后端路由的对应关系：

| 前端函数 | HTTP 方法 | 路径 |
| --- | --- | --- |
| `getNotices` | GET | `/system/notices` |
| `createNotice` | POST | `/system/notices` |
| `updateNotice` | POST | `/system/notices/:id/update` |
| `updateNoticeStatus` | POST | `/system/notices/:id/status` |

所有函数都通过 `http` 实例发送请求，自动带上 Token 和错误处理。返回值直接解包为业务数据，页面调用时不需要再处理 `response.data.data`。

## 前端：页面

公告管理页面包含搜索栏、数据表格、分页和弹窗表单，是一个典型的后台 CRUD 页面。文件较长，折叠查看：

::: details `admin/src/pages/system/NoticeView.vue` — 公告管理页面完整代码
```vue
<script setup lang="ts">
import { CloseOutline } from '@vicons/ionicons5'
import type { DataTableColumns, FormInst, FormRules } from 'naive-ui'
import {
  NButton,
  NCard,
  NDataTable,
  NForm,
  NFormItem,
  NIcon,
  NInput,
  NInputNumber,
  NModal,
  NPagination,
  NPopconfirm,
  NSelect,
  NSpace,
  NTag,
  useMessage,
} from 'naive-ui'
import { h, onMounted, reactive, ref } from 'vue'

import { createNotice, getNotices, updateNotice, updateNoticeStatus } from '../../api/notice'
import { buttonPermissionCodes } from '../../router/dynamic-menu'
import {
  NoticeStatus,
  type NoticeItem,
  type NoticeListQuery,
} from '../../types/notice'

interface NoticeFormModel {
  id: number
  title: string
  content: string
  sort: number
  status: NoticeStatus
  remark: string
}

const message = useMessage()
const loading = ref(false)
const saving = ref(false)
const notices = ref<NoticeItem[]>([])
const total = ref(0)

const query = reactive<NoticeListQuery>({
  page: 1,
  page_size: 10,
  keyword: '',
  status: 0,
})

const formRef = ref<FormInst | null>(null)
const formVisible = ref(false)
const formMode = ref<'create' | 'edit'>('create')
const formModel = reactive<NoticeFormModel>({
  id: 0,
  title: '',
  content: '',
  sort: 0,
  status: NoticeStatus.Enabled,
  remark: '',
})

const statusFilterOptions = [
  { label: '状态：全部', value: 0 },
  { label: '启用', value: NoticeStatus.Enabled },
  { label: '禁用', value: NoticeStatus.Disabled },
]

const statusFormOptions = [
  { label: '启用', value: NoticeStatus.Enabled },
  { label: '禁用', value: NoticeStatus.Disabled },
]

const rules: FormRules = {
  title: [{ required: true, message: '请输入公告标题', trigger: 'blur' }],
}

const columns: DataTableColumns<NoticeItem> = [
  {
    title: '标题',
    key: 'title',
    width: 220,
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', { class: 'font-semibold text-[#111827]' }, row.title)
    },
  },
  {
    title: '内容',
    key: 'content',
    minWidth: 240,
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', { class: 'text-[#374151]' }, row.content || '-')
    },
  },
  {
    title: '排序',
    key: 'sort',
    width: 80,
    align: 'center',
  },
  {
    title: '状态',
    key: 'status',
    width: 90,
    align: 'center',
    render(row) {
      return h(
        NTag,
        { type: row.status === NoticeStatus.Enabled ? 'success' : 'error', bordered: false },
        { default: () => (row.status === NoticeStatus.Enabled ? '启用' : '禁用') },
      )
    },
  },
  {
    title: '更新时间',
    key: 'updated_at',
    width: 160,
    render(row) {
      return h('span', { class: 'tabular-nums text-[#6B7280]' }, formatTime(row.updated_at))
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 180,
    fixed: 'right',
    render(row) {
      const nextStatus = row.status === NoticeStatus.Enabled ? NoticeStatus.Disabled : NoticeStatus.Enabled

      return h(
        NSpace,
        { size: 8, align: 'center' },
        {
          default: () =>
            [
              canUse('system:notice:update')
                ? h(
                    NButton,
                    { size: 'small', ghost: true, type: 'info', onClick: () => openEdit(row) },
                    { default: () => '编辑' },
                  )
                : null,
              canUse('system:notice:status')
                ? h(
                    NPopconfirm,
                    { onPositiveClick: () => handleToggleStatus(row, nextStatus) },
                    {
                      trigger: () =>
                        h(
                          NButton,
                          {
                            size: 'small',
                            ghost: true,
                            type: nextStatus === NoticeStatus.Disabled ? 'error' : 'success',
                          },
                          { default: () => (nextStatus === NoticeStatus.Disabled ? '禁用' : '启用') },
                        ),
                      default: () => `确认${nextStatus === NoticeStatus.Disabled ? '禁用' : '启用'}该公告？`,
                    },
                  )
                : null,
            ].filter(Boolean),
        },
      )
    },
  },
]

function canUse(code: string) {
  return buttonPermissionCodes.value.includes(code)
}

function formatTime(value: string) {
  if (!value) return '-'
  const d = new Date(value)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

function resetForm() {
  Object.assign(formModel, {
    id: 0,
    title: '',
    content: '',
    sort: 0,
    status: NoticeStatus.Enabled,
    remark: '',
  })
}

function handleSearch() {
  query.page = 1
  void loadNotices()
}

function handleReset() {
  query.page = 1
  query.page_size = 10
  query.keyword = ''
  query.status = 0
  void loadNotices()
}

function handlePageChange(page: number) {
  query.page = page
  void loadNotices()
}

function handlePageSizeChange(pageSize: number) {
  query.page = 1
  query.page_size = pageSize
  void loadNotices()
}

function openCreate() {
  formMode.value = 'create'
  resetForm()
  formVisible.value = true
}

function openEdit(row: NoticeItem) {
  formMode.value = 'edit'
  Object.assign(formModel, {
    id: row.id,
    title: row.title,
    content: row.content,
    sort: row.sort,
    status: row.status,
    remark: row.remark,
  })
  formVisible.value = true
}

async function loadNotices() {
  loading.value = true
  try {
    const data = await getNotices({
      ...query,
      keyword: query.keyword?.trim() || undefined,
      status: query.status === 0 ? undefined : query.status,
    })
    notices.value = data.items
    total.value = data.total
  } finally {
    loading.value = false
  }
}

async function handleSubmit() {
  await formRef.value?.validate()
  saving.value = true
  try {
    if (formMode.value === 'create') {
      await createNotice({
        title: formModel.title,
        content: formModel.content,
        sort: formModel.sort,
        status: formModel.status,
        remark: formModel.remark,
      })
      message.success('公告创建成功')
    } else {
      await updateNotice(formModel.id, {
        title: formModel.title,
        content: formModel.content,
        sort: formModel.sort,
        status: formModel.status,
        remark: formModel.remark,
      })
      message.success('公告更新成功')
    }

    formVisible.value = false
    await loadNotices()
  } finally {
    saving.value = false
  }
}

async function handleToggleStatus(row: NoticeItem, status: NoticeStatus) {
  await updateNoticeStatus(row.id, { status })
  message.success('公告状态已更新')
  await loadNotices()
}

onMounted(() => {
  void loadNotices()
})
</script>

<template>
  <main class="h-full overflow-hidden">
    <section class="flex h-full flex-col gap-4 overflow-hidden">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-[26px] font-bold text-[#111827]">公告管理</h1>
          <p class="mt-1 text-sm text-[#6B7280]">管理系统公告，支持按标题搜索和状态筛选。</p>
        </div>

        <NButton v-if="canUse('system:notice:create')" type="primary" @click="openCreate">
          + 新增公告
        </NButton>
      </div>

      <NCard :bordered="false" class="rounded-lg">
        <NSpace align="center" :wrap="true">
          <NInput
            v-model:value="query.keyword"
            clearable
            placeholder="公告标题"
            class="w-56"
            @keyup.enter="handleSearch"
          />
          <NSelect v-model:value="query.status" :options="statusFilterOptions" class="w-36" />
          <NButton type="primary" @click="handleSearch">查询</NButton>
          <NButton @click="handleReset">重置</NButton>
        </NSpace>
      </NCard>

      <NCard
        class="min-h-0 flex-1 rounded-lg"
        :bordered="false"
        content-style="height: 100%; padding: 0;"
      >
        <div class="flex items-center justify-between border-b border-[#E5E7EB] px-4 py-3">
          <span class="text-sm text-[#6B7280]">共 {{ total }} 条</span>
          <NButton text type="primary" @click="loadNotices">刷新</NButton>
        </div>

        <NDataTable
          remote
          class="notice-table h-full"
          style="height: calc(100% - 105px)"
          :columns="columns"
          :data="notices"
          :loading="loading"
          :pagination="false"
          :row-key="(row: NoticeItem) => row.id"
          :bordered="false"
          flex-height
        />

        <div
          class="flex items-center justify-between border-t border-[#E5E7EB] px-4 py-3 text-sm text-[#6B7280]"
        >
          <span>共 {{ total }} 条</span>
          <NPagination
            :page="query.page"
            :page-size="query.page_size"
            :item-count="total"
            :page-sizes="[10, 20, 50]"
            show-size-picker
            @update:page="handlePageChange"
            @update:page-size="handlePageSizeChange"
          />
        </div>
      </NCard>
    </section>

    <NModal
      v-model:show="formVisible"
      preset="card"
      :closable="false"
      class="compact-notice-modal"
      style="width: 600px; max-width: calc(100vw - 32px)"
    >
      <template #header>
        <div class="modal-header modal-header--hero">
          <h2 class="modal-header__title">
            {{ formMode === 'create' ? '新增公告' : '编辑公告' }}
          </h2>
          <p class="modal-header__hero-title">
            {{
              formMode === 'create'
                ? '填写公告标题和内容，保存后可立即展示'
                : '修改公告标题和内容，状态变更即时生效'
            }}
          </p>
          <button type="button" class="modal-close" @click="formVisible = false">
            <NIcon :size="18">
              <CloseOutline />
            </NIcon>
          </button>
        </div>
      </template>

      <div class="notice-modal-shell">
        <NForm
          ref="formRef"
          class="compact-notice-form"
          :model="formModel"
          :rules="rules"
          label-placement="left"
          label-width="76"
        >
          <section class="form-section form-section--primary">
            <div class="form-section__head">
              <h3>公告信息</h3>
              <p>标题不超过 128 个字符，内容支持任意文本。</p>
            </div>

            <div class="form-section-grid">
              <NFormItem label="标题" path="title">
                <NInput v-model:value="formModel.title" placeholder="公告标题" />
              </NFormItem>

              <NFormItem label="排序">
                <NInputNumber v-model:value="formModel.sort" :min="0" class="w-full" />
              </NFormItem>
            </div>
          </section>

          <section class="form-section form-section--muted">
            <div class="form-section__head">
              <h3>公告内容</h3>
            </div>

            <NFormItem label="内容" class="mb-0">
              <NInput
                v-model:value="formModel.content"
                type="textarea"
                :rows="4"
                placeholder="请输入公告内容"
              />
            </NFormItem>
          </section>

          <section class="form-section form-section--muted">
            <div class="form-section-grid">
              <NFormItem label="状态">
                <NSelect v-model:value="formModel.status" :options="statusFormOptions" />
              </NFormItem>

              <NFormItem label="备注">
                <NInput v-model:value="formModel.remark" placeholder="可选" />
              </NFormItem>
            </div>
          </section>
        </NForm>
      </div>

      <template #footer>
        <div class="modal-footer-actions">
          <NButton quaternary class="modal-footer-button" @click="formVisible = false">
            取消
          </NButton>
          <NButton
            type="primary"
            class="modal-footer-button modal-footer-button--primary"
            :loading="saving"
            @click="handleSubmit"
          >
            保存
          </NButton>
        </div>
      </template>
    </NModal>
  </main>
</template>

<style scoped>
.notice-table :deep(.n-data-table-th) {
  font-weight: 700;
  color: #4B5563;
  background: #F9FAFB;
  font-size: 13px;
}

.notice-table :deep(.n-data-table-td) {
  color: #374151;
  font-size: 14px;
  padding: 10px 16px;
}

.notice-table :deep(.n-data-table-tr:hover .n-data-table-td) {
  background: unset !important;
}

.notice-table :deep(.n-data-table-tr) {
  transition: none;
}

.notice-table :deep(.n-data-table-tr:hover) {
  filter: brightness(0.97);
}

.compact-notice-modal :deep(.n-card) {
  overflow: hidden;
  border-radius: 32px;
  border: 1px solid #dfe9f5;
  background: #ffffff;
  box-shadow: 0 24px 72px rgba(15, 23, 42, 0.16);
}

.compact-notice-modal :deep(.n-card-header) {
  padding: 0;
  border-bottom: 1px solid #dfe9f5;
  background: linear-gradient(135deg, #eff6ff 0%, #e8f2ff 58%, #f4f9ff 100%);
}

.compact-notice-modal :deep(.n-card-header__main) {
  font-size: 19px;
  font-weight: 600;
  letter-spacing: 0.01em;
  color: #111827;
}

.compact-notice-modal :deep(.n-card__content) {
  padding: 20px 28px 10px;
}

.compact-notice-modal :deep(.n-card__footer) {
  padding: 16px 28px 24px;
  border-top: 1px solid #edf2f7;
  background: rgba(248, 250, 252, 0.85);
}

.compact-notice-form :deep(.n-form-item) {
  margin-bottom: 16px;
}

.compact-notice-form :deep(.n-form-item-label) {
  white-space: nowrap;
  align-items: center;
  padding-right: 14px;
  font-weight: 600;
  color: #374151;
}

.compact-notice-form :deep(.n-form-item-blank) {
  min-height: 40px;
}

.compact-notice-form :deep(.n-input-wrapper) {
  border-radius: 10px;
  background: #fbfcfe;
}

.compact-notice-form :deep(.n-base-selection) {
  border-radius: 10px;
  background: #fbfcfe;
}

.compact-notice-form :deep(.n-input),
.compact-notice-form :deep(.n-base-selection) {
  box-shadow: none;
}

.compact-notice-form {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.notice-modal-shell {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.modal-header {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.modal-header--hero {
  position: relative;
  overflow: hidden;
  min-height: 120px;
  padding: 26px 28px 22px;
  background:
    radial-gradient(circle at top right, rgba(34, 197, 94, 0.12), transparent 24%),
    linear-gradient(135deg, #eff6ff 0%, #e8f2ff 58%, #f4f9ff 100%);
}

.modal-header__title {
  position: relative;
  z-index: 1;
  font-size: 19px;
  font-weight: 600;
  line-height: 1.3;
  color: #111827;
}

.modal-header__hero-title {
  position: relative;
  z-index: 1;
  font-size: 14px;
  font-weight: 600;
  line-height: 1.6;
  color: #0f172a;
}

.modal-close {
  position: absolute;
  top: 20px;
  right: 22px;
  z-index: 2;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 38px;
  height: 38px;
  border: none;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.76);
  color: #64748b;
  box-shadow: 0 10px 24px rgba(148, 163, 184, 0.12);
  backdrop-filter: blur(8px);
  cursor: pointer;
  transition:
    background-color 0.2s ease,
    color 0.2s ease,
    box-shadow 0.2s ease,
    transform 0.2s ease;
}

.modal-close:hover {
  background: #ffffff;
  color: #111827;
  box-shadow: 0 14px 28px rgba(148, 163, 184, 0.18);
  transform: translateY(-1px);
}

.form-section {
  border: 1px solid #e9eff6;
  border-radius: 14px;
  background: #ffffff;
  padding: 18px 18px 4px;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.9);
}

.form-section--primary {
  border-color: #d9e7f8;
  background: linear-gradient(180deg, #ffffff 0%, #fcfdff 100%);
}

.form-section--muted {
  background: linear-gradient(180deg, #fcfdff 0%, #f9fbff 100%);
}

.form-section__head {
  margin-bottom: 12px;
}

.form-section__head h3 {
  font-size: 15px;
  font-weight: 600;
  color: #111827;
}

.form-section__head p {
  margin-top: 4px;
  font-size: 12px;
  line-height: 1.6;
  color: #6b7280;
}

.form-section-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
  column-gap: 20px;
}

.modal-footer-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.modal-footer-button {
  min-width: 92px;
  height: 40px;
  border-radius: 10px;
}

.modal-footer-button--primary {
  box-shadow: 0 10px 24px rgba(34, 197, 94, 0.18);
}

.mb-0 {
  margin-bottom: 0;
}

@media (max-width: 720px) {
  .form-section-grid {
    grid-template-columns: minmax(0, 1fr);
  }

  .compact-notice-modal :deep(.n-card-header),
  .compact-notice-modal :deep(.n-card__content),
  .compact-notice-modal :deep(.n-card__footer) {
    padding-left: 20px;
    padding-right: 20px;
  }

  .compact-notice-modal :deep(.n-card-header) {
    padding-bottom: 0;
  }

  .modal-header--hero {
    padding: 22px 20px 18px;
    min-height: 110px;
  }

  .modal-close {
    top: 18px;
    right: 18px;
  }
}
</style>
```
:::

页面的核心结构可以拆成四个部分来理解：

1. **搜索区** — 关键字输入框 + 状态下拉 + 查询/重置按钮。查询时重置到第一页，重置时清空所有条件。
2. **表格区** — `NDataTable` 使用 `remote` 模式，分页、排序都由后端处理。列定义中用 `render` 函数自定义了标题加粗、状态标签、时间格式化和操作按钮。
3. **分页区** — `NPagination` 放在表格底部，支持切换页码和每页条数。
4. **弹窗表单** — `NModal` + `NForm`，支持新建和编辑两种模式。表单校验规则只要求标题必填。

::: details 按钮权限是怎么生效的
页面上每个操作按钮都用 `canUse('system:notice:create')` 这样的方式控制可见性。`canUse` 函数读取 `dynamic-menu.ts` 中导出的 `buttonPermissionCodes`，这个值是从后端 `/auth/menus` 接口返回的按钮权限列表中收集的。只有当前用户所属角色被授权了对应的按钮权限编码，按钮才会渲染出来。

如果没有看到某个按钮，排查顺序是：角色管理里是否勾选了该按钮权限 → 菜单管理里按钮是否启用 → Bootstrap 里菜单种子是否正确创建。
:::

## 前端：路由映射

最后一步，在 `dynamic-menu.ts` 的 `routeComponentMap` 中加一行，把后端菜单的 `Component` 值映射到实际的 Vue 组件：

```ts
const routeComponentMap: Record<string, RouteComponent> = {
  'system/HealthView': () => import('../pages/system/HealthView.vue'),
  'system/UserView': () => import('../pages/system/UserView.vue'),
  'system/RoleView': () => import('../pages/system/RoleView.vue'),
  'system/MenuView': () => import('../pages/system/MenuView.vue'),
  'system/ConfigView': () => import('../pages/system/ConfigView.vue'),
  'system/FileView': () => import('../pages/system/FileView.vue'),
  'system/OperationLogView': () => import('../pages/system/OperationLogView.vue'),
  'system/LoginLogView': () => import('../pages/system/LoginLogView.vue'),
  'system/NoticeView': () => import('../pages/system/NoticeView.vue'), // [!code ++]
}
```

这一行是菜单能加载到真实页面的关键。`dynamic-menu.ts` 中的 `resolveRouteComponent` 函数会拿后端返回的 `Component` 字段（这里是 `"system/NoticeView"`）去 `routeComponentMap` 里查找对应的懒加载函数。找到就加载真实组件，找不到就降级到占位页面。

## 验证

模块接入完成后，按下面的步骤逐一验证。

### 1. 数据库迁移

`sys_notice` 表的建表语句放在 `server/migrations/` 下的迁移文件中。重启服务后 golang-migrate 会自动执行，不需要手动建表。

启动后端服务：

```bash
cd server
go run main.go
```

启动日志中应该能看到类似输出：

```text
INFO	database migrations applied
INFO	server started	{"addr": ":8080", "env": "dev"}
```

### 2. 创建管理员账号

服务启动后，先通过初始化接口创建管理员账号：

```bash
# 创建管理员账号
curl -X POST http://localhost:8080/api/v1/setup/init \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"YourPassword123","nickname":"管理员"}'
```

### 3. 接口验证

使用 `curl` 验证接口是否正常工作。先登录获取 Token：

```bash
TOKEN=$(curl -s http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"YourPassword123"}' \
  | jq -r '.data.access_token')
```

查询公告列表（应该返回空列表）：

```bash
curl -s http://localhost:8080/api/v1/system/notices \
  -H "Authorization: Bearer $TOKEN" | jq
```

期望输出：

```json
{
  "code": 0,
  "data": {
    "items": [],
    "total": 0,
    "page": 1,
    "page_size": 10
  }
}
```

创建一条公告：

```bash
curl -s -X POST http://localhost:8080/api/v1/system/notices \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{
    "title": "系统上线公告",
    "content": "后台管理系统已正式上线，欢迎使用。",
    "sort": 0,
    "status": 1,
    "remark": "首条公告"
  }' | jq
```

期望输出中 `data.id` 大于 0，`data.title` 为 `"系统上线公告"`。

再次查询列表，`total` 应为 `1`，`items` 中包含刚创建的记录。

### 3. 前端页面验证

1. 打开浏览器，登录后台管理系统。
2. 侧边栏"系统管理"下应该出现"公告管理"菜单项（图标为 `notification`）。
3. 点击进入，页面顶部显示"公告管理"标题和"+ 新增公告"按钮。
4. 点击"新增公告"，弹窗中填写标题和内容，点击"保存"。
5. 表格中出现新建的公告，状态显示绿色"启用"标签。
6. 点击"禁用"按钮，确认后状态切换为红色"禁用"标签。

::: warning 菜单看不到的排查顺序
如果侧边栏没有出现"公告管理"，按这个顺序检查：

1. 后端是否正常启动，日志里有没有 `default menu created menu_code=system:notice`。
2. 角色管理中 `super_admin` 角色的菜单权限是否包含公告相关条目（Bootstrap 会自动绑定，但如果数据库里已有旧数据，可能需要手动勾选）。
3. 浏览器控制台 Network 面板，查看 `/auth/menus` 接口返回的菜单列表是否包含 `system:notice`。
4. 清除浏览器缓存后重新登录。
:::

## 小结

公告管理模块走完了一条完整的接入链路，涉及的所有文件和改动点可以汇总成一张表：

| 层 | 文件 | 改动类型 |
| --- | --- | --- |
| Model | `server/internal/model/notice.go` | 新增 |
| Handler | `server/internal/handler/system/notices.go` | 新增 |
| Router | `server/internal/router/router.go` | 追加 5 行 |
| Migration | `server/migrations/{postgres,mysql}/000003_notice_seed_data.up.sql` | 新增 |
| Types | `admin/src/types/notice.ts` | 新增 |
| API | `admin/src/api/notice.ts` | 新增 |
| Page | `admin/src/pages/system/NoticeView.vue` | 新增 |
| Route | `admin/src/router/dynamic-menu.ts` | 追加 1 行 |

这就是[模块固定结构](./module-structure)里定义的约定在真实代码里的落地方式。以后接入新模块，按同样的顺序和结构走一遍就行：先写 Model，再写 Handler，然后注册路由，创建数据库迁移文件来初始化权限和菜单，最后接前端。

回到本章目录：[第 6 章：业务模块接入规范](./index)。
