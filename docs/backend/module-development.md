---
title: 模块开发
description: 后端模块的固定结构、开发流程和约定
---

# 模块开发

## 模块目录结构

每个后端模块遵循固定的四层目录结构：

```
internal/modules/{module}/
├── api/
│   ├── handlers.go    HTTP 处理器
│   ├── dto.go         请求/响应 DTO
│   └── routes.go      路由注册
├── application/
│   └── *.service.go   业务逻辑
├── infra/
│   └── repository.go  数据访问
└── domain/
    └── types.go       领域类型和常量
```

## 各层职责

### Domain 层

定义模块的核心类型、常量和枚举：

```go
package domain

type ItemStatus int

const (
    ItemStatusEnabled  ItemStatus = 1
    ItemStatusDisabled ItemStatus = 2
)

const (
    PermissionList   = "mymodule:item:list"
    PermissionCreate = "mymodule:item:create"
    PermissionUpdate = "mymodule:item:update"
    PermissionDelete = "mymodule:item:delete"
)
```

### Repository 层

数据访问，GORM 查询，数据权限注入：

```go
type Repository struct{}

func (r *Repository) List(db *gorm.DB, actor *actorx.Actor, query ListQuery) ([]Item, int64, error) {
    db = datascope.ApplyScopes(db, actor)
    var items []Item
    var total int64
    // ... 查询逻辑
    return items, total, nil
}
```

### Service 层

业务逻辑编排，调用 Repository，不直接操作数据库：

```go
type Service struct {
    repo *infra.Repository
}

func (s *Service) List(ctx *gin.Context, query ListQuery) ([]Item, int64, error) {
    actor := actorx.FromContext(ctx)
    return s.repo.List(db, actor, query)
}
```

### Handler 层

HTTP 请求处理，参数绑定和响应：

```go
type Handler struct {
    svc *application.Service
}

func (h *Handler) List(ctx *gin.Context) {
    var query ListQuery
    if err := ctx.ShouldBindQuery(&query); err != nil {
        httpx.BadRequest(ctx, err.Error())
        return
    }
    items, total, err := h.svc.List(ctx, query)
    if err != nil {
        httpx.InternalError(ctx, err.Error())
        return
    }
    httpx.Success(ctx, gin.H{"items": items, "total": total})
}
```

## DTO 约定

### 请求 DTO

使用 `binding` 标签进行参数验证：

```go
type CreateRequest struct {
    Name   string `json:"name" binding:"required"`
    Status int    `json:"status" binding:"oneof=1 2"`
    Remark string `json:"remark"`
}
```

### 列表查询 DTO

支持分页和关键词搜索：

```go
type ListQuery struct {
    Page     int    `form:"page"`
    PageSize int    `form:"page_size"`
    Keyword  string `form:"keyword"`
    Status   int    `form:"status"`
}
```

### 响应 DTO

从领域模型中提取需要返回的字段：

```go
type Response struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    Status    int       `json:"status"`
    CreatedAt time.Time `json:"created_at"`
}
```

## 路由注册

在 `api/routes.go` 中定义路由注册函数：

```go
func RegisterRoutes(rg *gin.RouterGroup, opts RouteOptions) {
    h := NewHandler(opts)

    g := rg.Group("/mymodule")
    g.GET("/items", h.List)
    g.POST("/items", h.Create)
    g.GET("/items/:id", h.Get)
    g.PUT("/items/:id", h.Update)
    g.DELETE("/items/:id", h.Delete)
}
```

然后在 `bootstrap/router.go` 中引入。

## 数据权限接入

在 Repository 的列表查询中使用 `datascope.ApplyScopes`：

```go
func (r *Repository) List(db *gorm.DB, actor *actorx.Actor) ([]Item, int64, error) {
    db = datascope.ApplyScopes(db, actor)
    // 后续查询自动带上数据权限过滤
}
```

业务表需要有 `department_id` 或 `creator_id` 字段才能参与数据权限过滤。

## 错误处理

使用 `errorsx` 包统一错误处理：

```go
errorsx.BadRequest("参数错误")
errorsx.Unauthorized("未认证")
errorsx.Forbidden("无权限")
errorsx.NotFound("资源不存在")
errorsx.InternalError("内部错误")
```

## 新增模块检查清单

1. 模块目录结构完整
2. GORM 模型已添加到 `platform/model/`
3. 数据库迁移文件已创建（MySQL + PostgreSQL）
4. 路由已注册到 `bootstrap/router.go`
5. Casbin 权限策略已配置
6. 菜单种子数据已更新
7. 列表接口已接入数据权限
8. DTO 已定义请求验证
