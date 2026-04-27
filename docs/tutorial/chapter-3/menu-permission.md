---
title: 菜单权限设计
description: "设计菜单、按钮和角色菜单关系，并返回当前用户可见的菜单树。"
---

# 菜单权限设计

前面已经能判断接口访问权限。这一节继续补齐前端管理台需要的菜单权限：用户登录后，根据角色拿到自己能看到的目录、菜单和按钮。

::: tip 🎯 本节目标
完成后，数据库中会新增 `sys_menu` 和 `sys_role_menu` 两张表；启动服务时会初始化系统管理菜单，并把它授权给 `super_admin`；访问 `/api/v1/auth/menus` 可以返回当前用户菜单树。
:::

## 本节会改什么

本节会新增或修改下面这些文件：

```text
server/
├─ internal/
│  ├─ handler/
│  │  └─ auth/
│  │     └─ menus.go
│  ├─ model/
│  │  ├─ menu.go
│  │  └─ role_menu.go
│  └─ router/
│     └─ router.go
└─ migrations/
   ├─ pgsql/
   │  └─ 000002_seed_data.up.sql
   └─ mysql/
      └─ 000002_seed_data.up.sql
```

| 位置 | 用途 |
| --- | --- |
| `internal/model/menu.go` | 定义目录、菜单、按钮模型 |
| `internal/model/role_menu.go` | 定义角色和菜单的绑定关系 |
| `migrations/{pgsql,mysql}/000002_seed_data.up.sql` | 初始化默认菜单，并授权给超级管理员 |
| `internal/handler/auth/menus.go` | 返回当前用户可见菜单树 |
| `internal/router/router.go` | 注册 `/api/v1/auth/menus` |

## 菜单权限关系

本节落地下面这条关系：

```text
用户 sys_user
  ↓
用户角色关系 sys_user_role
  ↓
角色 sys_role
  ↓
角色菜单关系 sys_role_menu
  ↓
菜单 sys_menu
```

`sys_menu` 同时承载三类数据：

| 类型 | 含义 | 示例 |
| --- | --- | --- |
| `1` | 目录 | 系统管理 |
| `2` | 菜单 | 系统状态 |
| `3` | 按钮 | 查看系统状态 |

::: warning ⚠️ 菜单权限不替代接口权限
菜单权限控制“前端展示什么”；Casbin 控制“接口能不能访问”。即使前端隐藏了某个按钮，后端接口仍然必须做权限判断。
:::

## 先创建数据表

本节新增 `sys_menu` 和 `sys_role_menu`，分别用于保存目录、菜单、按钮权限点，以及角色和菜单的绑定关系。

::: tip 建表 SQL
字段说明、菜单类型、唯一编码、关系表约定和 PostgreSQL / MySQL 建表语句统一放在参考手册：

- [数据库建表语句 - `sys_menu`](../../reference/database-ddl#sys-menu)
- [数据库建表语句 - `sys_role_menu`](../../reference/database-ddl#sys-role-menu)
:::

## 🛠️ 创建菜单模型

创建 `server/internal/model/menu.go`。这是新增文件，直接完整写入即可。

```go
package model

import (
	"time"

	"gorm.io/gorm"
)

// MenuType 表示菜单节点类型。
type MenuType int

const (
	// MenuTypeDirectory 表示目录节点。
	MenuTypeDirectory MenuType = 1
	// MenuTypeMenu 表示可访问页面。
	MenuTypeMenu MenuType = 2
	// MenuTypeButton 表示页面内按钮或操作点。
	MenuTypeButton MenuType = 3
)

// MenuStatus 表示菜单状态。
type MenuStatus int

const (
	// MenuStatusEnabled 表示菜单正常启用。
	MenuStatusEnabled MenuStatus = 1
	// MenuStatusDisabled 表示菜单已禁用。
	MenuStatusDisabled MenuStatus = 2
)

// Menu 是后台菜单和按钮权限模型。
type Menu struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ParentID  uint           `gorm:"not null;default:0;index" json:"parent_id"`
	Type      MenuType       `gorm:"type:smallint;not null" json:"type"`
	Code      string         `gorm:"size:128;not null;uniqueIndex" json:"code"`
	Title     string         `gorm:"size:64;not null" json:"title"`
	Path      string         `gorm:"size:255;not null;default:''" json:"path"`
	Component string         `gorm:"size:255;not null;default:''" json:"component"`
	Icon      string         `gorm:"size:64;not null;default:''" json:"icon"`
	Sort      int            `gorm:"not null;default:0" json:"sort"`
	Status    MenuStatus     `gorm:"type:smallint;not null;default:1" json:"status"`
	Remark    string         `gorm:"size:255;not null;default:''" json:"remark"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 固定菜单表名。
func (Menu) TableName() string {
	return "sys_menu"
}
```

## 🛠️ 创建角色菜单关系模型

创建 `server/internal/model/role_menu.go`。这是新增文件，直接完整写入即可。

```go
package model

import "time"

// RoleMenu 是角色和菜单的绑定关系。
type RoleMenu struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	RoleID    uint      `gorm:"not null;uniqueIndex:uk_sys_role_menu_role_menu;index:idx_sys_role_menu_role_id" json:"role_id"`
	MenuID    uint      `gorm:"not null;uniqueIndex:uk_sys_role_menu_role_menu;index:idx_sys_role_menu_menu_id" json:"menu_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 固定角色菜单关系表名。
func (RoleMenu) TableName() string {
	return "sys_role_menu"
}
```

## 🛠️ 初始化默认菜单

默认菜单已经在数据库迁移文件中初始化。迁移文件会在服务启动时自动执行，创建系统管理相关的菜单数据和角色菜单绑定关系。

::: tip 💡 菜单初始化
- 菜单数据：在 `migrations/{pgsql,mysql}/000002_seed_data.up.sql` 中插入系统管理目录、菜单和按钮
- 角色菜单绑定：在同一迁移文件中绑定 `super_admin` 角色到系统管理菜单
- 初始菜单包括：系统管理目录、系统状态菜单和查看系统状态按钮
:::

::: warning ⚠️ 菜单初始化只提供最小起步数据
迁移文件中初始化了一组菜单，方便验证菜单权限链路。后续真正的菜单新增、编辑、排序和授权，要放在系统管理接口中完成。
:::

## 🛠️ 创建当前用户菜单接口

创建 `server/internal/handler/auth/menus.go`。这是新增文件，直接完整写入即可。

```go
package auth

import (
	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/middleware"
	"ez-admin-gin/server/internal/model"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// MenuHandler 负责当前用户菜单相关接口。
type MenuHandler struct {
	db  *gorm.DB
	log *zap.Logger
}

// NewMenuHandler 创建菜单 Handler。
func NewMenuHandler(db *gorm.DB, log *zap.Logger) *MenuHandler {
	return &MenuHandler{
		db:  db,
		log: log,
	}
}

type menuResponse struct {
	ID        uint           `json:"id"`
	ParentID  uint           `json:"parent_id"`
	Type      model.MenuType `json:"type"`
	Code      string         `json:"code"`
	Title     string         `json:"title"`
	Path      string         `json:"path"`
	Component string         `json:"component"`
	Icon      string         `json:"icon"`
	Sort      int            `json:"sort"`
	Children  []menuResponse `json:"children,omitempty"`
}

type menuNode struct {
	menuResponse
	children []*menuNode
}

// Menus 返回当前登录用户可见的菜单树。
func (h *MenuHandler) Menus(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, apperror.Unauthorized("请先登录"), h.log)
		return
	}

	var menus []model.Menu
	err := h.db.
		Table("sys_menu AS m").
		Select("DISTINCT m.*").
		Joins("JOIN sys_role_menu AS rm ON rm.menu_id = m.id").
		Joins("JOIN sys_user_role AS ur ON ur.role_id = rm.role_id").
		Joins("JOIN sys_role AS r ON r.id = ur.role_id").
		Where("ur.user_id = ?", userID).
		Where("m.status = ?", model.MenuStatusEnabled).
		Where("r.status = ?", model.RoleStatusEnabled).
		Where("m.deleted_at IS NULL").
		Where("r.deleted_at IS NULL").
		Order("m.sort ASC, m.id ASC").
		Find(&menus).Error
	if err != nil {
		response.Error(c, apperror.Internal("查询菜单失败", err), h.log)
		return
	}

	response.Success(c, buildMenuTree(menus))
}

func buildMenuTree(menus []model.Menu) []menuResponse {
	nodes := make(map[uint]*menuNode, len(menus))

	for _, menu := range menus {
		nodes[menu.ID] = &menuNode{
			menuResponse: menuResponse{
				ID:        menu.ID,
				ParentID:  menu.ParentID,
				Type:      menu.Type,
				Code:      menu.Code,
				Title:     menu.Title,
				Path:      menu.Path,
				Component: menu.Component,
				Icon:      menu.Icon,
				Sort:      menu.Sort,
			},
		}
	}

	roots := make([]*menuNode, 0)
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

	return menuNodesToResponses(roots)
}

func menuNodesToResponses(nodes []*menuNode) []menuResponse {
	result := make([]menuResponse, 0, len(nodes))
	for _, node := range nodes {
		item := node.menuResponse
		item.Children = menuNodesToResponses(node.children)
		result = append(result, item)
	}

	return result
}
```

::: details 为什么 `/auth/menus` 不再挂 Casbin 权限
这个接口本身就是“根据当前登录用户返回自己的菜单”。只要用户已经登录，就可以请求；真正能看到哪些菜单，由 `sys_role_menu` 决定。

如果再给它加 Casbin 权限，容易出现“没有菜单权限就连菜单列表也拿不到”的绕口问题。
:::

## 🛠️ 注册菜单接口

修改 `server/internal/router/router.go`。这一处只需要在认证路由里增加菜单 Handler 和路由。

```go
// registerAuthRoutes 注册认证相关路由。
func registerAuthRoutes(r *gin.Engine, opts Options) {
	login := authHandler.NewLoginHandler(opts.DB, opts.Log, opts.Token)
	me := authHandler.NewMeHandler(opts.Log)
	menus := authHandler.NewMenuHandler(opts.DB, opts.Log) // [!code ++]

	api := r.Group("/api/v1")
	auth := api.Group("/auth")
	auth.POST("/login", login.Login)

	protectedAuth := auth.Group("")
	protectedAuth.Use(middleware.Auth(opts.Token, opts.Log))
	protectedAuth.GET("/me", me.Me)
	protectedAuth.GET("/menus", menus.Menus) // [!code ++]
}
```

## ✅ 整理依赖并启动

本节没有新增第三方依赖，但修改了模型、初始化逻辑和路由，仍然可以整理一次：

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
INFO	default menu created	{"menu_code": "system"}
INFO	default menu created	{"menu_code": "system:health"}
INFO	default menu created	{"menu_code": "system:health:view"}
INFO	default role menu bound	{"role_id": 1, "menu_id": 1}
```

## ✅ 验证菜单和授权数据

打开另一个终端，在项目根目录执行：

```bash
# 查看默认菜单
docker compose -f deploy/compose.local.yml exec postgres psql -U ez_admin -d ez_admin -c "select id, parent_id, type, code, title, path, sort from sys_menu order by sort, id;"
```

应该看到类似结果：

```text
 id | parent_id | type |        code         |    title     |      path       | sort
----+-----------+------+---------------------+--------------+-----------------+------
  1 |         0 |    1 | system              | 系统管理     | /system         |   10
  2 |         1 |    2 | system:health       | 系统状态     | /system/health  |   10
  3 |         2 |    3 | system:health:view  | 查看系统状态 |                 |   10
```

继续查看角色菜单绑定：

```bash
# 查看超级管理员绑定了哪些菜单
docker compose -f deploy/compose.local.yml exec postgres psql -U ez_admin -d ez_admin -c "select rm.id, r.code as role_code, m.code as menu_code from sys_role_menu rm join sys_role r on r.id = rm.role_id join sys_menu m on m.id = rm.menu_id order by rm.id;"
```

应该看到 `super_admin` 已经绑定上面三个菜单节点。

## ✅ 验证当前用户菜单接口

先登录拿到 Token，再请求 `/api/v1/auth/menus`：

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

Invoke-RestMethod `
  -Method Get `
  -Uri http://localhost:8080/api/v1/auth/menus `
  -Headers @{ Authorization = "Bearer $token" }
```

```bash [macOS / Linux]
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"EzAdmin@123456"}' | jq -r '.data.access_token')

curl -X GET http://localhost:8080/api/v1/auth/menus \
  -H "Authorization: Bearer ${TOKEN}"
```

:::

应该看到类似结果：

```json
{
  "code": 0,
  "message": "ok",
  "data": [
    {
      "id": 1,
      "parent_id": 0,
      "type": 1,
      "code": "system",
      "title": "系统管理",
      "path": "/system",
      "component": "",
      "icon": "setting",
      "sort": 10,
      "children": [
        {
          "id": 2,
          "parent_id": 1,
          "type": 2,
          "code": "system:health",
          "title": "系统状态",
          "path": "/system/health",
          "component": "system/HealthView",
          "icon": "monitor",
          "sort": 10,
          "children": [
            {
              "id": 3,
              "parent_id": 2,
              "type": 3,
              "code": "system:health:view",
              "title": "查看系统状态",
              "path": "",
              "component": "",
              "icon": "",
              "sort": 10
            }
          ]
        }
      ]
    }
  ]
}
```

::: details 为什么示例里按钮可能显示在 children 里
本节接口返回的是完整权限树，按钮节点也会作为子节点返回。后续前端可以按 `type` 区分：`1`、`2` 用来生成菜单和路由，`3` 用来控制按钮或操作点。
:::

## 常见问题

::: details 请求 `/api/v1/auth/menus` 提示 `请先登录`
这个接口需要登录。先调用 `/api/v1/auth/login` 拿到 `access_token`，再按下面格式传请求头：

```http
Authorization: Bearer <access_token>
```
:::

::: details 菜单接口返回空数组
优先检查三件事：

- `sys_menu` 里是否有启用菜单。
- `sys_role_menu` 是否已经把菜单绑定给 `super_admin`。
- 当前用户是否通过 `sys_user_role` 绑定了 `super_admin`。
:::

::: details 菜单和 Casbin 策略有什么区别
菜单权限控制前端展示；Casbin 策略控制后端接口访问。

两者可以有关联，但不要互相替代。隐藏菜单不等于接口安全，接口安全必须由后端权限校验保证。
:::

到这里，第 3 章的认证与权限主链路就完整了。
