---
name: module-generator
description: 当需要在本仓库中新增业务模块（CRUD）时使用。描述后端分层约定、前端分层约定、命名规范、权限常量、Swagger 注释模式和接入步骤。AI 按此规范生成代码，无需模板引擎。
---

# 业务模块生成规范

本 skill 用于指导新增 CRUD 业务模块。参照 `system/dict`（字典模块）作为完整范例，保持项目架构一致。

## 目录结构总览

一个标准业务模块（以 `{group}/{module}` 为例，如 `system/dict`）涉及以下文件：

```
server/internal/
├── platform/model/{module}.go              # GORM 模型 + 状态类型
└── modules/{group}/{module}/
    ├── domain/types.go                     # 领域类型、校验、权限常量、BuildResponse
    ├── api/
    │   ├── dto.go                          # API 层类型别名
    │   ├── routes.go                       # HTTP 路由注册
    │   └── handler.go                      # HTTP 处理器 + Swagger 注释
    ├── application/
    │   ├── ports.go                        # Repository 接口 + Transactor 类型
    │   └── service.go                      # 业务逻辑
    ├── infra/repository.go                 # GORM 仓储实现
    ├── routes.go                           # 模块入口：路由注册（组装 service → api）
    └── services.go                         # 模块入口：依赖装配（DB → repo → service）

admin/src/modules/{group}/
├── types/
│   ├── {module}.ts                         # API 类型定义
│   └── {module}-page.ts                    # 页面表单/查询类型
├── api/{module}.ts                         # HTTP 请求函数
├── composables/
│   ├── use{Module}Page.ts                  # 核心 composable（状态 + 副作用）
│   └── {module}-page.utils.ts              # 工具函数（默认值、转换、payload 构建）
├── components/
│   ├── {Module}FilterBar.vue               # 筛选栏
│   ├── {Module}Table.vue                   # 数据表格
│   ├── {Module}FormModal.vue               # 新增/编辑表单弹窗
│   └── ...                                 # 其他组件
└── pages/{Module}View.vue                  # 页面入口
```

## 后端分层约定

### 1. GORM 模型 (`platform/model/{module}.go`)

```go
package model

import (
    "time"
    "gorm.io/gorm"
)

// {StatusName} 表示{实体}状态。
type {StatusName} int

const (
    {StatusName}Enabled  {StatusName} = 1
    {StatusName}Disabled {StatusName} = 2
)

// {EntityName} 是{中文名}表模型。
type {EntityName} struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    // ... 业务字段 ...
    Status    {StatusName}   `gorm:"type:smallint;not null;default:1" json:"status"`
    Remark    string         `gorm:"size:255;not null;default:''" json:"remark"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func ({EntityName}) TableName() string {
    return "{table_name}"
}
```

要点：
- 状态用自定义类型 + 常量，不用裸 `int`
- 必须写 `TableName()` 方法
- 遵循 database-schema-design skill 的字段命名约定

### 2. Domain 层 (`domain/types.go`)

包含：请求/响应类型、Entity 别名、权限常量、Normalize 函数、BuildResponse 函数。

```go
package domain

type ListQuery struct {
    Page     int    `form:"page"`
    PageSize int    `form:"page_size"`
    Keyword  string `form:"keyword"`
    Status   int    `form:"status"`
}

type CreateRequest struct {
    Name   string             `json:"name"`
    Status model.XxxStatus    `json:"status"`
    // ...
}

type UpdateRequest struct { /* ... */ }
type UpdateStatusRequest struct { Status model.XxxStatus `json:"status"` }
type Response struct { /* ... */ }
type ListResponse struct {
    Items    []Response `json:"items"`
    Total    int64      `json:"total"`
    Page     int        `json:"page"`
    PageSize int        `json:"page_size"`
}

type Entity = model.XxxEntity

// 权限常量：{group}:{module}:{action}
const (
    PermissionList         = "{group}:{module}:list"
    PermissionCreate       = "{group}:{module}:create"
    PermissionUpdate       = "{group}:{module}:update"
    PermissionDelete       = "{group}:{module}:delete"
    PermissionUpdateStatus = "{group}:{module}:update_status"
)
```

Normalize 函数模式：
- 每个请求类型对应一个 `NormalizeXxxRequest(req) (T, error)` 函数
- 内部调用 `strings.TrimSpace`、长度校验、状态校验
- 返回清洗后的请求体，或 `errorsx.BadRequest(...)` 错误

BuildResponse 函数：将 model 实体映射为 API 响应体。

### 3. API 层 (`api/`)

**dto.go** — 纯类型别名，避免 handler 直接依赖 domain：

```go
package api

import {module}domain "ez-admin-gin/server/internal/modules/{group}/{module}/domain"

type ListQuery = {module}domain.ListQuery
type CreateRequest = {module}domain.CreateRequest
type Response = {module}domain.Response
type ListResponse = {module}domain.ListResponse
```

**routes.go** — 路由注册，接收 Service 而非 DB：

```go
type RouteOptions struct {
    Service *{module}app.Service
    Log     *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
    handler := NewHandler(opts.Service, opts.Log)
    group.GET("/{module}s", handler.List)
    group.POST("/{module}s", handler.Create)
    group.POST("/{module}s/:id/update", handler.Update)
    group.POST("/{module}s/:id/status", handler.UpdateStatus)
}
```

**handler.go** — HTTP 处理器 + Swagger 注释：

```go
// List godoc
// @Summary      查询{中文名}列表
// @Tags         {Group} / {中文名}管理
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "页码"
// @Param        page_size  query     int     false  "每页条数"
// @Success      200  {object}  httpx.Body{data={module}domain.ListResponse}
// @Failure      400  {object}  httpx.Body
// @Security     BearerAuth
// @Router       /{group}/{module}s [get]
func (h *Handler) List(c *gin.Context) {
    var query ListQuery
    if err := c.ShouldBindQuery(&query); err != nil {
        httpx.Error(c, errorsx.BadRequest("查询参数不正确"), h.log)
        return
    }
    result, err := h.service.List(query)
    if err != nil {
        httpx.WriteError(c, err, "查询{中文名}列表失败", h.log)
        return
    }
    httpx.Success(c, result)
}
```

Handler 统一模式：
- 参数绑定失败 → `errorsx.BadRequest` + `httpx.Error`
- 业务错误 → `httpx.WriteError`（自动记录日志）
- 成功 → `httpx.Success`
- 路径参数用 `httpx.UintIDParam`

### 4. Application 层 (`application/`)

**ports.go** — 接口定义：

```go
type {Module}Transactor = database.Transactor

type {Module}Repository interface {
    List(query {module}domain.ListQuery, page int, pageSize int, status *model.XxxStatus) ([]{module}domain.Entity, int64, error)
    FindByID(db *gorm.DB, id uint) ({module}domain.Entity, error)
    Create(db *gorm.DB, item *{module}domain.Entity) error
    Update(db *gorm.DB, item *{module}domain.Entity, req {module}domain.UpdateRequest) error
    UpdateStatus(db *gorm.DB, item *{module}domain.Entity, status model.XxxStatus) error
}
```

**service.go** — 业务逻辑：

```go
type Service struct {
    tx   {Module}Transactor
    repo {Module}Repository
}

func NewService(tx {Module}Transactor, repo {Module}Repository) *Service {
    return &Service{tx: tx, repo: repo}
}

// List — 查询列表
func (s *Service) List(query {module}domain.ListQuery) ({module}domain.ListResponse, error) {
    page, pageSize := paging.NormalizePage(query.Page, query.PageSize)
    // ... 状态过滤处理 ...
    items, total, err := s.repo.List(query, page, pageSize, status)
    if err != nil {
        return {module}domain.ListResponse{}, err
    }
    result := make([]{module}domain.Response, 0, len(items))
    for _, item := range items {
        result = append(result, {module}domain.BuildResponse(item))
    }
    return {module}domain.ListResponse{Items: result, Total: total, Page: page, PageSize: pageSize}, nil
}

// Create — 在事务中校验唯一性后创建
func (s *Service) Create(req {module}domain.CreateRequest) ({module}domain.Response, error) {
    req, err := {module}domain.NormalizeCreateRequest(req)
    if err != nil {
        return {module}domain.Response{}, err
    }
    // 构建实体 → 事务内校验唯一性 → 创建
}

// Update — 在事务中查找 → 更新
// UpdateStatus — 在事务中查找 → 校验状态 → 更新
```

事务模式统一用 `s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error { ... })`。

### 5. Infra 层 (`infra/repository.go`)

```go
type Repository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
    return &Repository{db: db}
}
```

关键模式：
- `List*` 方法用 `r.db`（不带事务）
- `FindByID`、`Create`、`Update*` 方法接收 `db *gorm.DB` 参数（事务内由 service 传入）
- `FindByID` 找不到记录返回 `errorsx.NotFound(...)`
- 唯一性检查用 `db.Unscoped().Where(...).First(...)` 包含已软删除记录
- 列表查询排序用 `Order("sort ASC, id ASC")`
- Update 方法执行后同步更新传入的 entity 指针字段（保持 service 层能拿到最新值）

### 6. 模块入口

**routes.go** — 被外部调用的路由注册入口：

```go
type RouteOptions struct {
    DB  *gorm.DB
    Log *zap.Logger
}

func RegisterRoutes(group *gin.RouterGroup, opts RouteOptions) {
    service := NewService(ServiceOptions{DB: opts.DB})
    {module}api.RegisterRoutes(group, {module}api.RouteOptions{
        Service: service,
        Log:     opts.Log,
    })
}
```

**services.go** — 依赖装配：

```go
type ServiceOptions struct {
    DB *gorm.DB
}

func NewService(opts ServiceOptions) *{module}app.Service {
    repo := {module}infra.NewRepository(opts.DB)
    transactor := platformDatabase.NewTransactor(opts.DB)
    return {module}app.NewService(transactor, repo)
}
```

## 前端分层约定

### 1. Types (`types/{module}.ts`)

```typescript
export enum XxxStatus {
  Enabled = 1,
  Disabled = 2,
}

export interface XxxItem {
  id: number
  name: string
  status: XxxStatus
  // ...
  created_at: string
  updated_at: string
}

export interface XxxListQuery {
  page: number
  page_size: number
  keyword?: string
  status?: number
}

export interface XxxListResponse {
  items: XxxItem[]
  total: number
  page: number
  page_size: number
}

export interface CreateXxxPayload { /* ... */ }
export interface UpdateXxxPayload { /* ... */ }
```

### 2. API (`api/{module}.ts`)

```typescript
import { http } from '@/api/http'
import type { XxxListQuery, XxxListResponse, XxxItem, CreateXxxPayload, UpdateXxxPayload } from '../types/{module}'

export function getXxxList(params: XxxListQuery) {
  return http.get<XxxListResponse>('/{group}/{module}s', { params })
}

export function getXxxById(id: number) {
  return http.get<XxxItem>(`/{group}/{module}s/${id}`)
}

export function createXxx(data: CreateXxxPayload) {
  return http.post<XxxItem>('/{group}/{module}s', data)
}

export function updateXxx(id: number, data: UpdateXxxPayload) {
  return http.post<XxxItem>(`/{group}/{module}s/${id}/update`, data)
}

export function updateXxxStatus(id: number, status: number) {
  return http.post(`/{group}/{module}s/${id}/status`, { status })
}
```

### 3. Utils (`composables/{module}-page.utils.ts`)

纯函数，负责：
- 默认查询参数构建
- 表单初始值
- API 响应 → 表单模型转换
- 表单模型 → 请求 payload 构建

### 4. Composable (`composables/use{Module}Page.ts`)

核心状态管理，返回页面所需的全部响应式数据和方法：

```typescript
export function useXxxPage() {
  const query = ref<XxxListQuery>(defaultQuery())
  const tableData = ref<XxxItem[]>([])
  const total = ref(0)
  const loading = ref(false)
  const formVisible = ref(false)
  const editingId = ref<number | null>(null)
  const formModel = ref<XxxFormModel>(defaultFormModel())

  const permissions = {
    list: 'group:module:list',
    create: 'group:module:create',
    update: 'group:module:update',
    delete: 'group:module:delete',
  }

  async function fetchList() { /* ... */ }
  function handleCreate() { /* ... */ }
  function handleEdit(row: XxxItem) { /* ... */ }
  async function handleSubmit() { /* ... */ }
  async function handleStatusChange(row: XxxItem, status: number) { /* ... */ }
  function handleSearch() { fetchList() }
  function handleReset() { /* defaultQuery + fetchList */ }

  return {
    query, tableData, total, loading,
    formVisible, editingId, formModel,
    permissions,
    fetchList, handleCreate, handleEdit, handleSubmit,
    handleStatusChange, handleSearch, handleReset,
  }
}
```

### 5. Components

- **FilterBar.vue** — 关键词搜索 + 状态筛选 + 新增按钮
- **Table.vue** — NDataTable 展示列表 + 操作按钮（编辑、状态切换）
- **FormModal.vue** — NModal + NForm，新增/编辑共用

组件通过 props 接收数据、通过 emits 上报交互，不直接调用 API。

### 6. Page (`pages/{Module}View.vue`)

编排层，只做拼装：

```vue
<script setup lang="ts">
import { useXxxPage } from '../composables/useXxxPage'
import XxxFilterBar from '../components/XxxFilterBar.vue'
import XxxTable from '../components/XxxTable.vue'
import XxxFormModal from '../components/XxxFormModal.vue'

const {
  query, tableData, total, loading,
  formVisible, editingId, formModel,
  permissions,
  fetchList, handleCreate, handleEdit, handleSubmit,
  handleStatusChange, handleSearch, handleReset,
} = useXxxPage()
</script>
```

## 接入步骤

新模块代码生成后，需要手动完成以下接入：

1. **注册路由** — 在对应 group 的父路由文件中 import 并调用 `RegisterRoutes`
2. **创建数据库表** — 编写 SQL 迁移脚本（遵循 database-schema-design skill）
3. **添加菜单图标** — 在 `admin/src/router/dynamic-menu.ts` 中添加图标映射
4. **添加菜单/权限种子数据** — 菜单 code: `{group}:{module}`，权限: `{group}:{module}:list`、`create`、`update`、`delete`、`update_status`
5. **前端路由自动发现** — 页面放在 `admin/src/modules/{group}/pages/{Module}View.vue` 即可被自动发现

## 命名速查

| 概念 | Go | TypeScript | 数据库 |
|------|-----|-----------|--------|
| 模块名 | `dict`（小写） | `dict`（小写） | `sys_dict_type` |
| 实体名 | `SystemDictType`（Pascal） | `DictTypeItem`（Pascal） | — |
| 状态类型 | `SystemDictStatus`（Pascal） | `DictStatus`（Pascal 枚举） | `smallint` |
| 权限常量 | `system:dict:type:list` | 同左（字符串） | — |
| 路由路径 | `/system/dict-types` | `/system/dict-types` | — |
| API 函数 | — | `getDictTypeList` | — |
| Composable | — | `useDictPage` | — |

## 生成前检查

输出代码前，确认：

- [ ] 后端分层完整：model → domain → api → application → infra → routes + services
- [ ] 前端分层完整：types → api → utils → composable → components → page
- [ ] 权限常量格式正确：`{group}:{module}:{action}`
- [ ] Swagger 注释包含 Summary、Tags、Success/Failure、Security、Router
- [ ] Repository 方法参数：查询用 `r.db`，写操作接收 `*gorm.DB`
- [ ] 事务模式：`s.tx.WithinTransaction(context.Background(), ...)`
- [ ] Normalize 函数覆盖所有必填字段校验
- [ ] BuildResponse 映射所有展示字段
- [ ] 前端组件通过 props/emits 通信，不直接调 API
