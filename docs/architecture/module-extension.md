---
title: 模块扩展机制
description: 后端模块固定结构、接入流程、前端页面接入、模块工具包
---

# 模块扩展机制

EZ Admin 的模块化设计让添加新业务模块变得标准化。无论是后端 API 还是前端页面，都遵循固定的结构和接入流程。

## 设计理念

- **模块自治**：每个模块有完整的四层结构，不依赖其他模块的内部实现
- **标准接入**：新模块只需按照约定创建文件，注册路由即可
- **平台层共享**：认证、授权、数据权限等横切关注点由 platform 层统一提供

## 后端模块结构

每个模块遵循固定的目录结构：

```
internal/modules/{module_name}/
├── api/
│   ├── handlers.go    HTTP 处理器（参数绑定、响应序列化）
│   ├── dto.go         请求/响应 DTO 定义
│   └── routes.go      路由注册
├── application/
│   └── *.service.go   业务逻辑（编排、校验、跨模块协调）
├── infra/
│   └── repository.go  数据访问（GORM 查询 + 数据权限）
└── domain/
    └── types.go       领域类型、常量、枚举
```

## 后端模块接入流程

### 1. 创建模块目录

```
internal/modules/mymodule/
├── api/
├── application/
├── infra/
└── domain/
```

### 2. 定义领域类型

在 `domain/types.go` 中定义模块的常量、枚举和 DTO：

```go
package domain

const (
    PermissionList   = "mymodule:item:list"
    PermissionCreate = "mymodule:item:create"
    PermissionUpdate = "mymodule:item:update"
    PermissionDelete = "mymodule:item:delete"
)
```

### 3. 实现 Repository

在 `infra/repository.go` 中实现数据访问：

```go
type Repository struct{}

func NewRepository() *Repository {
    return &Repository{}
}

func (r *Repository) List(db *gorm.DB, actor *actorx.Actor, query ListQuery) ([]Item, int64, error) {
    db = datascope.ApplyScopes(db, actor)  // 注入数据权限
    // ... 查询逻辑
}
```

### 4. 实现 Service

在 `application/` 中实现业务逻辑：

```go
type Service struct {
    repo *infra.Repository
}

func NewService(repo *infra.Repository) *Service {
    return &Service{repo: repo}
}
```

### 5. 定义 DTO 和 Handler

在 `api/dto.go` 定义请求/响应结构，在 `api/handlers.go` 实现 HTTP 处理器。

### 6. 注册路由

在 `api/routes.go` 中注册路由，并在 `bootstrap/router.go` 中引入：

```go
func RegisterRoutes(r *gin.RouterGroup, opts RouteOptions) {
    h := NewHandler(opts.Service)
    g := r.Group("/mymodule")
    g.GET("/items", h.List)
    g.POST("/items", h.Create)
    // ...
}
```

### 7. 添加数据库迁移

在 `server/migrations/` 中添加迁移 SQL（MySQL + PostgreSQL）。

### 8. 配置 Casbin 策略

为新接口添加 Casbin 权限策略，通常在迁移的种子数据中写入。

## 前端模块接入流程

### 1. 创建模块目录

```
admin/src/modules/mymodule/
├── api/
├── types/
├── composables/
├── components/
└── pages/
```

### 2. 定义类型

```typescript
// types/item.ts
export interface Item {
  id: number
  name: string
  status: number
}
```

### 3. 封装 API

```typescript
// api/item.ts
import http from '@/api/http'
import type { Item } from '../types/item'

export function getItemList(params: ListQuery) {
  return http.get<{ items: Item[]; total: number }>('/mymodule/items', { params })
}
```

### 4. 实现 Composable

```typescript
// composables/useItemPage.ts
export function useItemPage() {
  const items = ref<Item[]>([])
  const loading = ref(false)

  async function fetchList() {
    loading.value = true
    // ... API 调用
  }

  return { items, loading, fetchList }
}
```

### 5. 实现组件和页面

- `components/` 中放展示组件（FilterBar、Table、FormModal）
- `pages/` 中放页面组件（编排 composable + component）

### 6. 注册动态菜单

在后端菜单管理中添加新菜单项，`component` 字段填写 `mymodule/ItemView`。

### 7. 添加组件白名单

在 `admin/src/router/dynamic-menu.ts` 的组件映射中添加新页面：

```typescript
const componentMap = {
  'mymodule/ItemView': () => import('@/modules/mymodule/pages/ItemView.vue'),
}
```

## modulekit 工具包

`internal/modules/modulekit/` 提供模块注册的公共接口和工具函数，减少重复代码。

## 完整接入检查清单

- [ ] 后端模块目录结构完整（api/application/infra/domain）
- [ ] 数据库迁移文件已添加（MySQL + PostgreSQL）
- [ ] Casbin 权限策略已配置
- [ ] 菜单种子数据已更新
- [ ] 前端模块目录结构完整（api/types/composables/components/pages）
- [ ] 前端组件白名单已更新
- [ ] API 代理路径正确
- [ ] 列表接口支持数据权限过滤
